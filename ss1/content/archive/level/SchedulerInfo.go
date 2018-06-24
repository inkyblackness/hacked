package level

const (
	schedulerDefaultSize        = 64
	schedulerDefaultElementSize = 8
)

// SchedulerInfo describes the basic parameters of a level scheduler.
// Its entries are stored in the Schedules table.
type SchedulerInfo struct {
	// Size describes the maximum amount of schedule entries in the table.
	Size int32
	// ScheduleCount is the current amount of active schedules.
	ScheduleCount int32
	// ElementSize is the size of an schedule entry, in bytes.
	ElementSize int32
	// GrowFlag indicates whether the schedule table can be resized.
	GrowFlag byte
	_        [4]byte
	_        [4]byte
}

// DefaultSchedulerInfo returns a new scheduler info structure with default values.
func DefaultSchedulerInfo() SchedulerInfo {
	return SchedulerInfo{
		Size:        schedulerDefaultSize,
		ElementSize: schedulerDefaultElementSize,
	}
}
