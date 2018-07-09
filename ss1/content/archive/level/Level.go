package level

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
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

	objectMasterTable ObjectMasterTable

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

		resStart:       resourceBase.Plus(lvlids.PerLevel * id),
		tileMap:        NewTileMap(64, 64),
		wallHeightsMap: NewWallHeightsMap(64, 64),
	}
	lvl.resEnd = lvl.resStart.Plus(lvlids.PerLevel)

	lvl.reloadBaseInfo()
	lvl.reloadTileMap()
	lvl.reloadTextureAtlas()
	lvl.reloadObjectMasterTable()
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
func (lvl *Level) Tile(x, y int) *TileMapEntry {
	return lvl.tileMap.Tile(x, y)
}

// MapGridInfo returns the information necessary to draw a 2D map.
func (lvl *Level) MapGridInfo(x, y int) (TileType, TileSlopeControl, WallHeights) {
	tile := lvl.tileMap.Tile(x, y)
	if tile == nil {
		return TileTypeSolid, TileSlopeControlCeilingInverted, WallHeights{}
	}
	return tile.Type, tile.Flags.SlopeControl(), *lvl.wallHeightsMap.Tile(x, y)
}

// ObjectLimit returns the highest object ID that can be stored in this level.
func (lvl *Level) ObjectLimit() ObjectID {
	size := len(lvl.objectMasterTable)
	if size == 0 {
		return 0
	}
	return ObjectID(size - 1)
}

// ForEachObject iterates over all active objects and calls the given handler.
func (lvl *Level) ForEachObject(handler func(ObjectID, ObjectMasterEntry)) {
	tableSize := len(lvl.objectMasterTable)
	if tableSize > 0 {
		id := ObjectID(lvl.objectMasterTable[0].CrossReferenceTableIndex)
		for (id > 0) && (int(id) < tableSize) {
			entry := lvl.objectMasterTable[id]
			handler(id, entry)
			id = entry.Next
		}
	}
}

// EncodeState returns a subset of encoded level data, which only includes
// data that is loaded (modified) by the level structure.
// For any data block that is not relevant, a zero length slice is returned.
func (lvl *Level) EncodeState() [lvlids.PerLevel][]byte {
	var levelData [lvlids.PerLevel][]byte

	levelData[lvlids.Information] = encode(&lvl.baseInfo)

	levelData[lvlids.TextureAtlas] = encode(lvl.textureAtlas)
	levelData[lvlids.TileMap] = encode(lvl.tileMap)
	levelData[lvlids.ObjectMasterTable] = encode(lvl.objectMasterTable)

	levelData[lvlids.TextureAnimations] = encode(lvl.textureAnimations)
	levelData[lvlids.SurveillanceSources] = encode(&lvl.surveillanceSources)
	levelData[lvlids.SurveillanceSurrogates] = encode(&lvl.surveillanceSurrogates)
	levelData[lvlids.Parameters] = encode(&lvl.parameters)

	return levelData
}

func (lvl Level) encode(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(data)
	return buf.Bytes()
}

func (lvl *Level) onLevelResourceDataChanged(id int) {
	switch id {
	case lvlids.Information:
		lvl.reloadBaseInfo()
	case lvlids.TextureAtlas:
		lvl.reloadTextureAtlas()
	case lvlids.TileMap:
		lvl.reloadTileMap()
	case lvlids.ObjectMasterTable:
		lvl.reloadObjectMasterTable()
	case lvlids.TextureAnimations:
		lvl.reloadTextureAnimations()
	case lvlids.SurveillanceSources:
		lvl.reloadSurveillanceSources()
	case lvlids.SurveillanceSurrogates:
		lvl.reloadSurveillanceSurrogates()
	case lvlids.Parameters:
		lvl.reloadParameters()
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
	reader, err := lvl.reader(lvlids.TileMap)
	if err == nil {
		coder := serial.NewDecoder(reader)
		lvl.tileMap.Code(coder)
		err = coder.FirstError()
	}
	if err != nil {
		lvl.clearTileMap()
	}
	lvl.wallHeightsMap.CalculateFrom(lvl.tileMap)
}

func (lvl *Level) reloadObjectMasterTable() {
	reader, err := lvl.reader(lvlids.ObjectMasterTable)
	if err != nil {
		lvl.objectMasterTable = nil
		return
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		lvl.objectMasterTable = nil
		return
	}
	lvl.objectMasterTable = make([]ObjectMasterEntry, len(data)/ObjectMasterEntrySize)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, lvl.objectMasterTable)
	if err != nil {
		lvl.objectMasterTable = nil
	}
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
	for _, row := range lvl.tileMap {
		for i := 0; i < len(row); i++ {
			row[i].Reset()
		}
	}
}

func (lvl *Level) reader(block int) (io.Reader, error) {
	res, err := lvl.localizer.LocalizedResources(resource.LangAny).Select(lvl.resStart.Plus(block))
	if err != nil {
		return nil, err
	}
	if res.ContentType() != resource.Archive {
		return nil, errors.New("resource is not for archive")
	}
	if res.BlockCount() != 1 {
		return nil, errors.New("resource has invalid block count")
	}
	return res.Block(0)
}
