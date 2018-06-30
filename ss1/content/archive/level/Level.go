package level

import (
	"encoding/binary"
	"errors"
	"io"

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

	return lvl
}

// InvalidateResources resets all internally cached data.
func (lvl *Level) InvalidateResources(ids []resource.ID) {
	for _, id := range ids {
		if (id >= lvl.resStart) && (id < lvl.resEnd) {
			lvl.onLevelResourceDataChanged(int(id.Value() - lvl.resStart.Value()))
		}
	}
}

// MapGridInfo returns the information necessary to draw a 2D map.
func (lvl *Level) MapGridInfo(x, y int) (TileType, WallHeights) {
	tile := lvl.tileMap.Tile(x, y)
	if tile == nil {
		return TileTypeSolid, WallHeights{}
	}
	return tile.Type, *lvl.wallHeightsMap.Tile(x, y)
}

func (lvl *Level) onLevelResourceDataChanged(id int) {
	switch id {
	case lvlids.Information:
		lvl.reloadBaseInfo()
	case lvlids.TileMap:
		lvl.reloadTileMap()
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
