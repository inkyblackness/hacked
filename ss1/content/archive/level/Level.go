package level

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlobj"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Level is the complete structure defining all necessary data for a level.
type Level struct {
	id        int
	localizer resource.Localizer

	resStart resource.ID
	resEnd   resource.ID

	baseInfo       BaseInfo
	tileMap        TileMap
	wallHeightsMap WallHeightsMap
	textureAtlas   TextureAtlas

	objectMainTable     ObjectMainTable
	objectCrossRefTable ObjectCrossReferenceTable
	objectClassTables   [object.ClassCount]ObjectClassTable

	textureAnimations      []TextureAnimationEntry
	surveillanceSources    [SurveillanceObjectCount]ObjectID
	surveillanceSurrogates [SurveillanceObjectCount]ObjectID
	parameters             Parameters
}

// NewLevel returns a new instance.
func NewLevel(resourceBase resource.ID, id int, localizer resource.Localizer) *Level {
	lvl := &Level{
		id: id,

		localizer: localizer,

		resStart: resourceBase.Plus(lvlids.PerLevel * id),
	}
	lvl.resEnd = lvl.resStart.Plus(lvlids.PerLevel)

	lvl.reloadBaseInfo()
	lvl.reloadTileMap()
	lvl.reloadTextureAtlas()
	lvl.reloadObjectMainTable()
	lvl.reloadObjectCrossRefTable()
	lvl.reloadObjectClassTables()
	lvl.reloadTextureAnimations()
	lvl.reloadSurveillanceSources()
	lvl.reloadSurveillanceSurrogates()
	lvl.reloadParameters()

	return lvl
}

// ID returns the identifier of the level.
func (lvl *Level) ID() int {
	return lvl.id
}

// InvalidateResources resets all internally cached data.
func (lvl *Level) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		if (id >= lvl.resStart) && (id < lvl.resEnd) {
			lvl.onLevelResourceDataChanged(int(id.Value() - lvl.resStart.Value()))
		}
	}
}

// Size returns the dimensions of the level.
func (lvl *Level) Size() (x, y int, z HeightShift) {
	return int(lvl.baseInfo.XSize), int(lvl.baseInfo.YSize), lvl.baseInfo.ZShift
}

// SetHeightShift changes the vertical scale of the level.
func (lvl *Level) SetHeightShift(value HeightShift) {
	lvl.baseInfo.ZShift = value
}

// TextureAnimations returns the list of animations.
func (lvl *Level) TextureAnimations() []TextureAnimationEntry {
	return lvl.textureAnimations
}

// SurveillanceSources returns the current sources of surveillance.
func (lvl *Level) SurveillanceSources() [SurveillanceObjectCount]ObjectID {
	return lvl.surveillanceSources
}

// SetSurveillanceSource sets a surveillance source.
func (lvl *Level) SetSurveillanceSource(index int, id ObjectID) {
	if (index >= 0) && (index < SurveillanceObjectCount) {
		lvl.surveillanceSources[index] = id
	}
}

// SurveillanceSurrogates returns the current surrogates of surveillance.
func (lvl *Level) SurveillanceSurrogates() [SurveillanceObjectCount]ObjectID {
	return lvl.surveillanceSurrogates
}

// SetSurveillanceSurrogate sets a surveillance source.
func (lvl *Level) SetSurveillanceSurrogate(index int, id ObjectID) {
	if (index >= 0) && (index < SurveillanceObjectCount) {
		lvl.surveillanceSurrogates[index] = id
	}
}

// IsCyberspace returns true if the level describes a cyberspace.
func (lvl *Level) IsCyberspace() bool {
	return lvl.baseInfo.Cyberspace != 0
}

// Parameters returns the extra properties the level has.
func (lvl *Level) Parameters() *Parameters {
	return &lvl.parameters
}

// TextureAtlas returns the atlas for textures.
func (lvl *Level) TextureAtlas() TextureAtlas {
	return lvl.textureAtlas
}

// SetTextureAtlasEntry changes the entry of the atlas.
func (lvl *Level) SetTextureAtlasEntry(index int, textureIndex TextureIndex) {
	if (index >= 0) && (index < len(lvl.textureAtlas)) {
		lvl.textureAtlas[index] = textureIndex
	}
}

// Tile returns the tile entry at given position.
func (lvl *Level) Tile(pos TilePosition) *TileMapEntry {
	return lvl.tileMap.Tile(pos, int(lvl.baseInfo.XShift))
}

// MapGridInfo returns the information necessary to draw a 2D map.
func (lvl *Level) MapGridInfo(pos TilePosition) (TileType, TileSlopeControl, WallHeights) {
	tile := lvl.Tile(pos)
	if tile == nil {
		return TileTypeSolid, TileSlopeControlCeilingInverted, WallHeights{}
	}
	return tile.Type, tile.Flags.SlopeControl(), *lvl.wallHeightsMap.Tile(pos)
}

// ObjectCapacity returns the number of objects that can be stored in this level.
func (lvl *Level) ObjectCapacity() int {
	return lvl.objectMainTable.Capacity()
}

// ObjectClassStats returns the number of used and total possible entries of given class.
func (lvl *Level) ObjectClassStats(class object.Class) (active, capacity int) {
	if int(class) >= len(lvl.objectClassTables) {
		return 0, 0
	}
	objectClassTable := lvl.objectClassTables[class]
	return objectClassTable.AllocatedCount(), objectClassTable.Capacity()
}

// ForEachObject iterates over all active objects and calls the given handler.
func (lvl *Level) ForEachObject(handler func(ObjectID, ObjectMainEntry)) {
	tableSize := len(lvl.objectMainTable)
	if tableSize > 0 {
		id := ObjectID(lvl.objectMainTable[0].CrossReferenceTableIndex)
		for (id > 0) && (int(id) < tableSize) {
			entry := lvl.objectMainTable[id]
			handler(id, entry)
			id = entry.Next
		}
	}
}

// HasRoomForObjectOf returns true if the level can hold one object of given class.
func (lvl *Level) HasRoomForObjectOf(class object.Class) bool {
	classActive, classCapacity := lvl.ObjectClassStats(class)
	levelActive := lvl.objectMainTable.AllocatedCount()
	levelCapacity := lvl.objectMainTable.Capacity()
	return (levelActive < levelCapacity) && (classActive < classCapacity)
}

// NewObject attempts to allocate a new object for given class.
// The created object has no place in the world and must be placed.
// Returns 0 and an error if not possible.
func (lvl *Level) NewObject(class object.Class) (ObjectID, error) {
	if int(class) >= len(lvl.objectClassTables) {
		return 0, errInvalidClass
	}
	classTable := lvl.objectClassTables[class]
	classIndex := classTable.Allocate()
	if classIndex == 0 {
		return 0, errNoMoreRoomForClass
	}
	id := lvl.objectMainTable.Allocate()
	if id == 0 {
		classTable.Release(classIndex)
		return 0, errNoMoreRoomForObjects
	}
	obj := &lvl.objectMainTable[id]
	classEntry := &classTable[classIndex]
	classEntry.ObjectID = id
	obj.ClassTableIndex = int16(classIndex)
	obj.Class = class

	lvl.addCrossReferenceTo(id, obj, offMapReferencePosition())

	return id, nil
}

// UpdateObjectLocation updates the reference table between object and tiles based on its current location.
func (lvl *Level) UpdateObjectLocation(id ObjectID) {
	obj := lvl.Object(id)
	if (obj == nil) || (obj.InUse == 0) {
		return
	}
	lvl.removeCrossReferences(int(obj.CrossReferenceTableIndex), func(entry ObjectCrossReferenceEntry) int { return int(entry.NextTileForObj) })
	lvl.addCrossReferenceTo(id, obj, obj.TilePosition())
}

// DelObject removes the identified object from the level.
func (lvl *Level) DelObject(id ObjectID) {
	obj := lvl.Object(id)
	if (obj == nil) || (obj.InUse == 0) {
		return
	}
	lvl.removeCrossReferences(int(obj.CrossReferenceTableIndex), func(entry ObjectCrossReferenceEntry) int { return int(entry.NextTileForObj) })
	if int(obj.Class) < len(lvl.objectClassTables) {
		lvl.objectClassTables[obj.Class].Release(int(obj.ClassTableIndex))
	}
	lvl.objectMainTable.Release(id)
}

func (lvl *Level) addCrossReferenceTo(id ObjectID, obj *ObjectMainEntry, pos TilePosition) {
	newIndex := lvl.objectCrossRefTable.Allocate()
	if newIndex == 0 {
		return
	}
	entry := &lvl.objectCrossRefTable[newIndex]
	entry.ObjectID = id
	entry.SetTilePosition(pos)

	tile := lvl.Tile(pos)
	if tile != nil {
		entry.NextInTile = tile.FirstObjectIndex
		tile.FirstObjectIndex = int16(newIndex)
	}
	if obj.CrossReferenceTableIndex != 0 {
		entry.NextTileForObj = obj.CrossReferenceTableIndex
		endIndex := obj.CrossReferenceTableIndex
		for lvl.objectCrossRefTable[endIndex].NextTileForObj != obj.CrossReferenceTableIndex {
			endIndex = lvl.objectCrossRefTable[endIndex].NextTileForObj
		}
		lvl.objectCrossRefTable[endIndex].NextTileForObj = int16(newIndex)
	} else {
		entry.NextTileForObj = int16(newIndex)
	}
	obj.CrossReferenceTableIndex = int16(newIndex)
}

func (lvl *Level) removeCrossReferences(start int, next func(ObjectCrossReferenceEntry) int) {
	index := start
	for (index != 0) && (index < len(lvl.objectCrossRefTable)) {
		entry := lvl.objectCrossRefTable[index]

		tile := lvl.Tile(entry.TilePosition())
		if tile != nil {
			if tile.FirstObjectIndex == int16(index) {
				tile.FirstObjectIndex = entry.NextInTile
			} else {
				otherIndex := tile.FirstObjectIndex
				for otherIndex != 0 {
					otherEntry := &lvl.objectCrossRefTable[otherIndex]
					otherIndex = otherEntry.NextInTile
					if otherEntry.NextInTile == int16(index) {
						otherEntry.NextInTile = entry.NextInTile
					}
				}
			}
		}

		obj := lvl.Object(entry.ObjectID)
		if obj != nil {
			switch {
			case entry.NextTileForObj == int16(start):
				obj.CrossReferenceTableIndex = 0
			case obj.CrossReferenceTableIndex == int16(index):
				obj.CrossReferenceTableIndex = entry.NextTileForObj
			default:
				otherIndex := obj.CrossReferenceTableIndex
				for otherIndex != 0 {
					otherEntry := &lvl.objectCrossRefTable[otherIndex]
					otherIndex = otherEntry.NextTileForObj
					if otherEntry.NextTileForObj == int16(index) {
						otherEntry.NextTileForObj = entry.NextTileForObj
					}
				}
			}
		}

		lvl.objectCrossRefTable.Release(index)

		index = next(entry)
		if index == start {
			index = 0
		}
	}
}

// Object returns the main entry for the identified object.
func (lvl *Level) Object(id ObjectID) *ObjectMainEntry {
	var entry *ObjectMainEntry
	if (id > 0) && (int(id) < len(lvl.objectMainTable)) {
		entry = &lvl.objectMainTable[id]
	}
	return entry
}

// ObjectClassData returns the raw class data for the given object.
func (lvl *Level) ObjectClassData(obj *ObjectMainEntry) *interpreters.Instance {
	if (obj == nil) || (obj.InUse == 0) {
		return interpreters.New().For(nil)
	}
	if int(obj.Class) >= len(lvl.objectClassTables) {
		return interpreters.New().For(nil)
	}
	classTable := lvl.objectClassTables[obj.Class]
	if (obj.ClassTableIndex < 1) || (int(obj.ClassTableIndex) >= len(classTable)) {
		return interpreters.New().For(nil)
	}

	interpreterFactory := lvlobj.ForRealWorld
	if lvl.IsCyberspace() {
		interpreterFactory = lvlobj.ForCyberspace
	}
	return interpreterFactory(obj.Triple(), classTable[obj.ClassTableIndex].Data)
}

// ObjectExtraData returns the raw extra data for the given object.
func (lvl *Level) ObjectExtraData(obj *ObjectMainEntry) *interpreters.Instance {
	if (obj == nil) || (obj.InUse == 0) {
		return interpreters.New().For(nil)
	}
	interpreterFactory := lvlobj.RealWorldExtra
	if lvl.IsCyberspace() {
		interpreterFactory = lvlobj.CyberspaceExtra
	}
	return interpreterFactory(obj.Triple(), obj.Extra[:])
}

// EncodeState returns a subset of encoded level data, which only includes
// data that is loaded (modified) by the level structure.
// For any data block that is not relevant, a zero length slice is returned.
func (lvl *Level) EncodeState() [lvlids.PerLevel][]byte {
	var levelData [lvlids.PerLevel][]byte

	levelData[lvlids.Information] = encode(&lvl.baseInfo)

	levelData[lvlids.TextureAtlas] = encode(lvl.textureAtlas)
	levelData[lvlids.TileMap] = encode(lvl.tileMap)
	levelData[lvlids.ObjectMainTable] = encode(lvl.objectMainTable)
	levelData[lvlids.ObjectCrossRefTable] = encode(lvl.objectCrossRefTable)
	for class := 0; class < len(lvl.objectClassTables); class++ {
		levelData[lvlids.ObjectClassTablesStart+class] = encode(lvl.objectClassTables[class])
	}

	levelData[lvlids.TextureAnimations] = encode(lvl.textureAnimations)
	levelData[lvlids.SurveillanceSources] = encode(&lvl.surveillanceSources)
	levelData[lvlids.SurveillanceSurrogates] = encode(&lvl.surveillanceSurrogates)
	levelData[lvlids.Parameters] = encode(&lvl.parameters)

	return levelData
}

func (lvl *Level) onLevelResourceDataChanged(id int) {
	switch id {
	case lvlids.Information:
		lvl.reloadBaseInfo()
	case lvlids.TextureAtlas:
		lvl.reloadTextureAtlas()
	case lvlids.TileMap:
		lvl.reloadTileMap()
	case lvlids.ObjectMainTable:
		lvl.reloadObjectMainTable()
	case lvlids.ObjectCrossRefTable:
		lvl.reloadObjectCrossRefTable()
	case lvlids.TextureAnimations:
		lvl.reloadTextureAnimations()
	case lvlids.SurveillanceSources:
		lvl.reloadSurveillanceSources()
	case lvlids.SurveillanceSurrogates:
		lvl.reloadSurveillanceSurrogates()
	case lvlids.Parameters:
		lvl.reloadParameters()
	}
	if (id >= lvlids.ObjectClassTablesStart) && (id < (lvlids.ObjectClassTablesStart + len(lvl.objectClassTables))) {
		lvl.reloadObjectClassTable(object.Class(id - lvlids.ObjectClassTablesStart))
	}
}

func (lvl *Level) reloadBaseInfo() {
	reader, err := lvl.reader(lvlids.Information)
	if err != nil {
		lvl.baseInfo = BaseInfo{}
		return
	}
	err = binary.Read(reader, binary.LittleEndian, &lvl.baseInfo)
	if err != nil {
		lvl.baseInfo = BaseInfo{}
	}
}

func (lvl *Level) reloadTextureAtlas() {
	reader, err := lvl.reader(lvlids.TextureAtlas)
	if err != nil {
		lvl.textureAtlas = nil
		return
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		lvl.textureAtlas = nil
		return
	}
	lvl.textureAtlas = make([]TextureIndex, len(data)/2)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, lvl.textureAtlas)
	if err != nil {
		lvl.textureAtlas = nil
	}
}

func (lvl *Level) reloadTileMap() {
	lvl.tileMap = NewTileMap(1<<lvl.baseInfo.XShift, 1<<lvl.baseInfo.YShift)
	reader, err := lvl.reader(lvlids.TileMap)
	if err == nil {
		coder := serial.NewDecoder(reader)
		lvl.tileMap.Code(coder)
		err = coder.FirstError()
	}
	if err != nil {
		lvl.clearTileMap()
	}
	lvl.wallHeightsMap = NewWallHeightsMap(int(lvl.baseInfo.XSize), int(lvl.baseInfo.YSize))
	lvl.wallHeightsMap.CalculateFrom(lvl)
}

func (lvl *Level) reloadObjectMainTable() {
	reader, err := lvl.reader(lvlids.ObjectMainTable)
	if err != nil {
		lvl.objectMainTable = nil
		return
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		lvl.objectMainTable = nil
		return
	}
	lvl.objectMainTable = make([]ObjectMainEntry, len(data)/ObjectMainEntrySize)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, lvl.objectMainTable)
	if err != nil {
		lvl.objectMainTable = nil
	}
}

func (lvl *Level) reloadObjectCrossRefTable() {
	reader, err := lvl.reader(lvlids.ObjectCrossRefTable)
	if err != nil {
		lvl.objectCrossRefTable = nil
		return
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		lvl.objectCrossRefTable = nil
		return
	}
	lvl.objectCrossRefTable = make([]ObjectCrossReferenceEntry, len(data)/ObjectCrossReferenceEntrySize)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, lvl.objectCrossRefTable)
	if err != nil {
		lvl.objectCrossRefTable = nil
	}
}

func (lvl *Level) reloadObjectClassTables() {
	for class := 0; class < len(lvl.objectClassTables); class++ {
		lvl.reloadObjectClassTable(object.Class(class))
	}
}

func (lvl *Level) reloadObjectClassTable(class object.Class) {
	lvl.objectClassTables[class] = nil
	reader, err := lvl.reader(lvlids.ObjectClassTablesStart + int(class))
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	info := ObjectClassInfoFor(class)
	table := make(ObjectClassTable, len(data)/(ObjectClassEntryHeaderSize+info.DataSize))
	table.AllocateData(info.DataSize)

	coder := serial.NewDecoder(bytes.NewReader(data))
	table.Code(coder)
	if coder.FirstError() != nil {
		return
	}
	lvl.objectClassTables[class] = table
}

func (lvl *Level) reloadTextureAnimations() {
	reader, err := lvl.reader(lvlids.TextureAnimations)
	if err != nil {
		lvl.textureAnimations = nil
		return
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		lvl.textureAnimations = nil
		return
	}
	lvl.textureAnimations = make([]TextureAnimationEntry, len(data)/TextureAnimationEntrySize)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, lvl.textureAnimations)
	if err != nil {
		lvl.textureAnimations = nil
	}
}

func (lvl *Level) reloadSurveillanceSources() {
	reader, err := lvl.reader(lvlids.SurveillanceSources)
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &lvl.surveillanceSources)
	}
	if err != nil {
		lvl.surveillanceSources = [SurveillanceObjectCount]ObjectID{}
	}
}

func (lvl *Level) reloadSurveillanceSurrogates() {
	reader, err := lvl.reader(lvlids.SurveillanceSurrogates)
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &lvl.surveillanceSurrogates)
	}
	if err != nil {
		lvl.surveillanceSurrogates = [SurveillanceObjectCount]ObjectID{}
	}
}

func (lvl *Level) reloadParameters() {
	reader, err := lvl.reader(lvlids.Parameters)
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &lvl.parameters)
	}
	if err != nil {
		lvl.parameters = DefaultParameters()
	}
}

func (lvl *Level) clearTileMap() {
	for i := 0; i < len(lvl.tileMap.entries); i++ {
		lvl.tileMap.entries[i].Reset()
	}
}

func (lvl *Level) reader(block int) (io.Reader, error) {
	resID := lvl.resStart.Plus(block)
	res, err := lvl.localizer.LocalizedResources(resource.LangAny).Select(resID)
	if err != nil {
		return nil, err
	}
	if (res.ContentType() != resource.Archive) || (res.BlockCount() != 1) {
		return nil, resource.ErrWrongType(resource.KeyOf(resID, resource.LangAny, 0), resource.Archive)
	}
	return res.Block(0)
}
