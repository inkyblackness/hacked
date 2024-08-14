package level

import "fmt"

// TileFlag describes simple properties of a map tile.
type TileFlag uint32

// RealWorldFlag describes simple properties of a map tile in the real world.
type RealWorldFlag TileFlag

// CyberspaceFlag describes simple properties of a map tile in cyberspace.
type CyberspaceFlag TileFlag

// ForRealWorld interprets the flag value for the real world.
func (flag TileFlag) ForRealWorld() RealWorldFlag {
	return RealWorldFlag(flag)
}

// ForCyberspace interprets the flag value for cyberspace.
func (flag TileFlag) ForCyberspace() CyberspaceFlag {
	return CyberspaceFlag(flag)
}

// MusicZone returns the music zone.
func (flag TileFlag) MusicZone() TileMusicZone {
	return TileMusicZone((flag & 0x0000E000) >> 13)
}

// WithMusicZone returns a new flag value with the given music index set.
func (flag TileFlag) WithMusicZone(value TileMusicZone) TileFlag {
	return TileFlag(uint32(flag&^0x0000E000) | (uint32(value) << 13))
}

// Peril returns whether the tile is marked as peril (may play peril music).
func (flag TileFlag) Peril() bool {
	return (flag & 0x00001000) != 0
}

// WithPeril returns a flag with the given peril set.
func (flag TileFlag) WithPeril(value bool) TileFlag {
	var valueFlag uint32
	if value {
		valueFlag = 0x00001000
	}
	return TileFlag(uint32(flag&^0x00001000) | valueFlag)
}

// SlopeControl returns the slope control as per flags.
func (flag TileFlag) SlopeControl() TileSlopeControl {
	return TileSlopeControl((flag & 0x00000C00) >> 10)
}

// WithSlopeControl returns a new flag value with the given slope control set.
func (flag TileFlag) WithSlopeControl(ctrl TileSlopeControl) TileFlag {
	return TileFlag(uint32(flag&^0x00000C00) | (uint32(ctrl&0x3) << 10))
}

// AsTileFlag returns the flags as regular tile flag value.
func (flag RealWorldFlag) AsTileFlag() TileFlag {
	return TileFlag(flag)
}

// WallTextureOffset returns the vertical offset (in tile height units) to apply for wall textures.
func (flag RealWorldFlag) WallTextureOffset() TileHeightUnit {
	return TileHeightUnit(flag & 0x0000001F)
}

// WithWallTextureOffset returns a flag value with the given wall texture offset.
func (flag RealWorldFlag) WithWallTextureOffset(value TileHeightUnit) RealWorldFlag {
	return RealWorldFlag(uint32(flag&^0x0000001F) | uint32(value&0x1F))
}

// WallTexturePattern returns the pattern to apply for walls.
func (flag RealWorldFlag) WallTexturePattern() WallTexturePattern {
	return WallTexturePattern(byte(flag&0x00000060) >> 5)
}

// WithWallTexturePattern returns a flag with the given pattern set.
func (flag RealWorldFlag) WithWallTexturePattern(value WallTexturePattern) RealWorldFlag {
	return RealWorldFlag(uint32(flag&^0x00000060) | (uint32(value&0x3) << 5))
}

// UseAdjacentWallTexture returns whether the wall texture from the adjacent tile should be used for each side.
func (flag RealWorldFlag) UseAdjacentWallTexture() bool {
	return (flag & 0x00000100) != 0
}

// WithUseAdjacentWallTexture returns a flag with the given usage set.
func (flag RealWorldFlag) WithUseAdjacentWallTexture(value bool) RealWorldFlag {
	var valueFlag uint32
	if value {
		valueFlag = 0x00000100
	}
	return RealWorldFlag(uint32(flag&^0x00000100) | valueFlag)
}

// Deconstructed returns whether the tile is marked as heavily deconstructed (should play spooky music).
func (flag RealWorldFlag) Deconstructed() bool {
	return (flag & 0x00000200) != 0
}

// WithDeconstructed returns a flag with the given deconstruction set.
func (flag RealWorldFlag) WithDeconstructed(value bool) RealWorldFlag {
	var valueFlag uint32
	if value {
		valueFlag = 0x00000200
	}
	return RealWorldFlag(uint32(flag&^0x00000200) | valueFlag)
}

// FloorShadow returns the floor shadow value. Range: [0..GradesOfShadow].
func (flag RealWorldFlag) FloorShadow() int {
	return int((flag & 0x000F0000) >> 16)
}

// WithFloorShadow returns a new flag value with the given floor shadow set. Values beyond allowed range are ignored.
func (flag RealWorldFlag) WithFloorShadow(value int) RealWorldFlag {
	if (value < 0) || (value >= GradesOfShadow) {
		return flag
	}
	return RealWorldFlag(uint32(flag&^0x000F0000) | (uint32(value) << 16))
}

// CeilingShadow returns the ceiling shadow value. Range: [0..GradesOfShadow].
func (flag RealWorldFlag) CeilingShadow() int {
	return int((flag & 0x0F000000) >> 24)
}

// WithCeilingShadow returns a new flag value with the given ceiling shadow set. Values beyond allowed range are ignored.
func (flag RealWorldFlag) WithCeilingShadow(value int) RealWorldFlag {
	if (value < 0) || (value >= GradesOfShadow) {
		return flag
	}
	return RealWorldFlag(uint32(flag&^0x0F000000) | (uint32(value) << 24))
}

// TileVisited returns whether the tile is marked as being visited (seen).
func (flag RealWorldFlag) TileVisited() bool {
	return (flag & 0x80000000) != 0
}

// WithTileVisited returns a flag with the given visited set.
func (flag RealWorldFlag) WithTileVisited(value bool) RealWorldFlag {
	var valueFlag uint32
	if value {
		valueFlag = 0x80000000
	}
	return RealWorldFlag(uint32(flag&^0x80000000) | valueFlag)
}

// AsTileFlag returns the flags as regular tile flag value.
func (flag CyberspaceFlag) AsTileFlag() TileFlag {
	return TileFlag(flag)
}

// GameOfLifeState returns the current game-of-life state.
func (flag CyberspaceFlag) GameOfLifeState() int {
	return int(byte(flag&0x00000060) >> 5)
}

// WithGameOfLifeState returns a flag with the given game-of-life state set.
func (flag CyberspaceFlag) WithGameOfLifeState(value int) CyberspaceFlag {
	return CyberspaceFlag(uint32(flag&^0x00000060) | (uint32(value&0x3) << 5))
}

// FlightPull returns the pull applying in a tile.
func (flag CyberspaceFlag) FlightPull() CyberspaceFlightPull {
	return CyberspaceFlightPull((uint32(flag&0x01000000) >> 20) | (uint32(flag&0x000F0000) >> 16))
}

// WithFlightPull returns a flag instance with the given pull applied.
func (flag CyberspaceFlag) WithFlightPull(value CyberspaceFlightPull) CyberspaceFlag {
	newFlag := uint32(flag &^ 0x010F0000)
	newFlag |= uint32(value&0x0F) << 16
	newFlag |= uint32(value&0x10) << 20
	return CyberspaceFlag(newFlag)
}

type TileMusicZone byte

const (
	TileMusicZoneNoMusic    TileMusicZone = 0
	TileMusicZoneHospital   TileMusicZone = 1
	TileMusicZoneExecutive  TileMusicZone = 2
	TileMusicZoneIndustrial TileMusicZone = 3
	TileMusicZoneMetal      TileMusicZone = 4
	TileMusicZonePark       TileMusicZone = 5
	TileMusicZoneBridge     TileMusicZone = 6
	TileMusicZoneElevator   TileMusicZone = 7
)

// TileSlopeControls returns all control values.
func TileMusicZones() []TileMusicZone {
	return []TileMusicZone{
		TileMusicZoneNoMusic,
		TileMusicZoneHospital,
		TileMusicZoneExecutive,
		TileMusicZoneIndustrial,
		TileMusicZoneMetal,
		TileMusicZonePark,
		TileMusicZoneBridge,
		TileMusicZoneElevator,
	}
}

func (ctrl TileMusicZone) String() string {
	switch ctrl {
	case TileMusicZoneNoMusic:
		return "No Music"
	case TileMusicZoneHospital:
		return "Hospital"
	case TileMusicZoneExecutive:
		return "Executive"
	case TileMusicZoneIndustrial:
		return "Industrial"
	case TileMusicZoneMetal:
		return "Metal"
	case TileMusicZonePark:
		return "Park"
	case TileMusicZoneBridge:
		return "Bridge"
	case TileMusicZoneElevator:
		return "Elevator"
	default:
		return fmt.Sprintf("Unknown%02X", int(ctrl))
	}
}
