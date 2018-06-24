package level

const (
	defaultMapSideLength  = 64
	defaultMapSideShift   = 6
	defaultMapHeightShift = 3
)

// BaseInfo describes the basic parameters of a level map.
type BaseInfo struct {
	// XSize is the horizontal extent of the map (West to East).
	XSize int32
	// YSize is the vertical extent of the map (South to North).
	YSize int32
	// XShift is the base value of XSize (XSize == 1 << XShift).
	XShift int32
	// YShift is the base value of YSize (YSize == 1 << YShift).
	YShift int32
	// ZShift is the base value of the height of the map.
	ZShift int32
	_      [4]byte
	// Cyberspace indicates whether the level is cyberspace.
	Cyberspace byte
	_          [12]byte
	// Scheduler contains base information about the schedule table of the level.
	Scheduler SchedulerInfo
}

// DefaultBaseInfo returns an initialized instance.
func DefaultBaseInfo(cyberspace bool) BaseInfo {
	info := BaseInfo{
		XSize:     defaultMapSideLength,
		YSize:     defaultMapSideLength,
		XShift:    defaultMapSideShift,
		YShift:    defaultMapSideShift,
		ZShift:    defaultMapHeightShift,
		Scheduler: DefaultSchedulerInfo(),
	}
	if cyberspace {
		info.Cyberspace = 1
	}
	return info
}
