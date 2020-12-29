package level

// WallHeights describes the height delta, in HeightUnits, one would take to
// cross to the other tile. A value of 0 means equal floor, a negative value a fall down.
// A value of 0x20 (= maximum of height units) means a solid wall.
//
// For each cardinal direction, facing to the specified direction,
// the first entry is the left side, the second the center, and the third the right side.
type WallHeights struct {
	North [3]float32
	East  [3]float32
	South [3]float32
	West  [3]float32
}

// Reset sets all the heights to zero.
func (heights *WallHeights) Reset() {
	*heights = WallHeights{}
}

// WallHeightsMap is rectangular table of wall heights, typically mirroring a tile map.
type WallHeightsMap [][]WallHeights

// NewWallHeightsMap returns a new, initialized map.
func NewWallHeightsMap(width, height int) WallHeightsMap {
	m := make([][]WallHeights, height)
	for y := 0; y < height; y++ {
		m[y] = make([]WallHeights, width)
		for x := 0; x < width; x++ {
			m[y][x].Reset()
		}
	}
	return m
}

// Tile returns a pointer to the tile within the map for given position.
func (m WallHeightsMap) Tile(x, y int) *WallHeights {
	return &m[y][x]
}

// CalculateFrom updates all the wall heights according to the specified map.
func (m *WallHeightsMap) CalculateFrom(tileMap interface{ Tile(x, y int) *TileMapEntry }) {
	for y, row := range *m {
		for x := 0; x < len(row); x++ {
			tile := tileMap.Tile(x, y)
			if tile == nil {
				continue
			}
			heights := &row[x]
			heights.North = calculateWallHeight(tile, DirNorth, tileMap.Tile(x, y+1), DirSouth)
			heights.East = calculateWallHeight(tile, DirEast, tileMap.Tile(x+1, y), DirWest)
			heights.South = calculateWallHeight(tile, DirSouth, tileMap.Tile(x, y-1), DirNorth)
			heights.West = calculateWallHeight(tile, DirWest, tileMap.Tile(x-1, y), DirEast)
		}
	}
}

func calculateWallHeight(entry *TileMapEntry, entrySide Direction, other *TileMapEntry, otherSide Direction) [3]float32 {
	var result [3]float32
	baseHeights := func(tile *TileMapEntry, side Direction) (floor TileHeightUnit, slope TileHeightUnit, ceiling TileHeightUnit) {
		floor = TileHeightUnitMax
		slope = TileHeightUnitMin
		ceiling = TileHeightUnitMax
		if (tile != nil) && ((tile.Type.Info().SolidSides & side.AsMask()) == 0) {
			floor = tile.Floor.AbsoluteHeight()
			slope = tile.SlopeHeight
			ceiling = tile.Ceiling.AbsoluteHeight()
		}
		return
	}

	otherFloorHeight, otherSlopeHeight, otherCeilingHeight := baseHeights(other, otherSide)
	var otherFloorFactors SlopeFactors
	var otherCeilingFactors SlopeFactors
	if other != nil {
		otherSlopeControl := other.Flags.SlopeControl()
		otherFloorFactors = otherSlopeControl.FloorSlopeFactors(other.Type)
		otherCeilingFactors = otherSlopeControl.CeilingSlopeFactors(other.Type)
	}

	entryFloorHeight, entrySlopeHeight, entryCeilingHeight := baseHeights(entry, entrySide)
	entrySlopeControl := entry.Flags.SlopeControl()
	entryFloorFactors := entrySlopeControl.FloorSlopeFactors(entry.Type)
	entryCeilingFactors := entrySlopeControl.CeilingSlopeFactors(entry.Type)

	sideHeights := func(side Direction, sideOffset int, factors SlopeFactors, slopeHeight TileHeightUnit, absoluteHeight TileHeightUnit) [3]float32 {
		return [3]float32{
			factors[side.Offset(-sideOffset)]*float32(slopeHeight) + float32(absoluteHeight),
			factors[side.Offset(0)]*float32(slopeHeight) + float32(absoluteHeight),
			factors[side.Offset(+sideOffset)]*float32(slopeHeight) + float32(absoluteHeight),
		}
	}
	entryFloorSideHeights := sideHeights(entrySide, 1, entryFloorFactors, entrySlopeHeight, entryFloorHeight)
	entryCeilingSideHeights := sideHeights(entrySide, 1, entryCeilingFactors.Negated(), entrySlopeHeight, entryCeilingHeight)
	otherFloorSideHeights := sideHeights(otherSide, -1, otherFloorFactors, otherSlopeHeight, otherFloorHeight)
	otherCeilingSideHeights := sideHeights(otherSide, -1, otherCeilingFactors.Negated(), otherSlopeHeight, otherCeilingHeight)

	for i := 0; i < 3; i++ {
		switch {
		case (otherFloorSideHeights[i] == entryFloorSideHeights[i]) && (entryFloorSideHeights[i] == float32(TileHeightUnitMax)):
			// shortcut for solid-to-solid - no "wall" to see here
			result[i] = 0
		case (otherFloorSideHeights[i] >= otherCeilingSideHeights[i]) || // other ceiling side seals off at its floor
			(otherCeilingSideHeights[i] <= entryFloorSideHeights[i]) || // other ceiling side seals off at this floor
			(otherFloorSideHeights[i] >= entryCeilingSideHeights[i]) || // other floor side seals off at this ceiling
			(entryFloorSideHeights[i] >= entryCeilingSideHeights[i]): // own ceiling side seals off at floor
			result[i] = float32(TileHeightUnitMax)
		default:
			result[i] = otherFloorSideHeights[i] - entryFloorSideHeights[i]
		}
	}

	return result
}
