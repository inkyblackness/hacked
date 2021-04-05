package levels

import (
	"sort"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

type hoverItem interface {
	Pos() MapPosition
	Size() float32
	IsIn(lvl *level.Level) bool
}

type tileHoverItem struct {
	pos MapPosition
}

func (item tileHoverItem) Pos() MapPosition {
	return item.pos
}

func (item tileHoverItem) Size() float32 {
	return level.FineCoordinatesPerTileSide
}

func (item tileHoverItem) IsIn(lvl *level.Level) bool {
	return lvl.Tile(item.pos.Tile()) != nil
}

type objectHoverItem struct {
	id  level.ObjectID
	pos MapPosition
}

func (item objectHoverItem) Pos() MapPosition {
	return item.pos
}

func (item objectHoverItem) Size() float32 {
	return level.FineCoordinatesPerTileSide / 4
}

func (item objectHoverItem) IsIn(lvl *level.Level) bool {
	obj := lvl.Object(item.id)
	return (obj != nil) && (obj.InUse != 0)
}

type hoverItems struct {
	available   []hoverItem
	activeIndex int
	activeItem  hoverItem
}

func (items *hoverItems) reset() {
	items.available = nil
	items.activeIndex = 0
	items.activeItem = nil
}

func (items *hoverItems) scroll(forward bool) {
	count := len(items.available)
	if count > 0 {
		diff := 1
		if !forward {
			diff = -1
		}
		items.activeIndex = (count + (items.activeIndex + diff)) % count
		items.activeItem = items.available[items.activeIndex]
	}
}

func (items *hoverItems) find(lvl *level.Level, pos MapPosition) {
	refVec := mgl.Vec2{float32(pos.X), float32(pos.Y)}
	var distances []float32

	items.available = nil
	lvl.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
		entryVec := mgl.Vec2{float32(entry.X), float32(entry.Y)}
		distance := refVec.Sub(entryVec).Len()
		if distance < level.FineCoordinatesPerTileSide/4 {
			items.available = append(items.available,
				objectHoverItem{
					id:  id,
					pos: MapPosition{X: entry.X, Y: entry.Y},
				})
			distances = append(distances, distance)
		}
	})
	items.available = append(items.available,
		tileHoverItem{
			pos: MapPosition{
				X: level.CoordinateAt(pos.X.Tile(), level.FineCoordinatesPerTileSide/2),
				Y: level.CoordinateAt(pos.Y.Tile(), level.FineCoordinatesPerTileSide/2),
			},
		})
	distances = append(distances, level.FineCoordinatesPerTileSide)
	sort.Slice(items.available, func(a, b int) bool { return distances[a] < distances[b] })

	items.activeIndex = 0
	items.activeItem = items.available[items.activeIndex]
}

func (items *hoverItems) validate(lvl *level.Level) {

}
