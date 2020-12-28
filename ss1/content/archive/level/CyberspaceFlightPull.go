package level

import "fmt"

// CyberspaceFlightPull describes a directed pull of hacker in cyberspace.
type CyberspaceFlightPull byte

// String returns a textual representation.
func (pull CyberspaceFlightPull) String() string {
	switch pull {
	case CyberspaceFlightPullNone:
		return "None"
	case CyberspaceFlightPullWeakEast:
		return "Weak East"
	case CyberspaceFlightPullWeakWest:
		return "Weak West"
	case CyberspaceFlightPullWeakNorth:
		return "Weak North"
	case CyberspaceFlightPullWeakSouth:
		return "Weak South"
	case CyberspaceFlightPullMediumEast:
		return "Medium East"
	case CyberspaceFlightPullMediumWest:
		return "Medium West"
	case CyberspaceFlightPullMediumNorth:
		return "Medium North"
	case CyberspaceFlightPullMediumSouth:
		return "Medium South"
	case CyberspaceFlightPullStrongEast:
		return "Strong East"
	case CyberspaceFlightPullStrongWest:
		return "Strong West"
	case CyberspaceFlightPullStrongNorth:
		return "Strong North"
	case CyberspaceFlightPullStrongSouth:
		return "Strong South"
	case CyberspaceFlightPullMediumCeiling:
		return "Medium Ceiling"
	case CyberspaceFlightPullMediumFloor:
		return "Medium Floor"
	case CyberspaceFlightPullStrongCeiling:
		return "Strong Ceiling"
	case CyberspaceFlightPullStrongFloor:
		return "Strong Floor"
	default:
		return fmt.Sprintf("Unknown%02X", int(pull))
	}
}

// CyberspaceFlightPull constants are listed below.
const (
	CyberspaceFlightPullNone          CyberspaceFlightPull = 0
	CyberspaceFlightPullWeakEast      CyberspaceFlightPull = 1
	CyberspaceFlightPullWeakWest      CyberspaceFlightPull = 2
	CyberspaceFlightPullWeakNorth     CyberspaceFlightPull = 3
	CyberspaceFlightPullWeakSouth     CyberspaceFlightPull = 4
	CyberspaceFlightPullMediumEast    CyberspaceFlightPull = 5
	CyberspaceFlightPullMediumWest    CyberspaceFlightPull = 6
	CyberspaceFlightPullMediumNorth   CyberspaceFlightPull = 7
	CyberspaceFlightPullMediumSouth   CyberspaceFlightPull = 8
	CyberspaceFlightPullStrongEast    CyberspaceFlightPull = 9
	CyberspaceFlightPullStrongWest    CyberspaceFlightPull = 10
	CyberspaceFlightPullStrongNorth   CyberspaceFlightPull = 11
	CyberspaceFlightPullStrongSouth   CyberspaceFlightPull = 12
	CyberspaceFlightPullMediumCeiling CyberspaceFlightPull = 13
	CyberspaceFlightPullMediumFloor   CyberspaceFlightPull = 14
	CyberspaceFlightPullStrongCeiling CyberspaceFlightPull = 15
	CyberspaceFlightPullStrongFloor   CyberspaceFlightPull = 16
)

// CyberspaceFlightPulls returns all constants of CyberspaceFlightPull.
func CyberspaceFlightPulls() []CyberspaceFlightPull {
	return []CyberspaceFlightPull{
		CyberspaceFlightPullNone,
		CyberspaceFlightPullWeakEast,
		CyberspaceFlightPullWeakWest,
		CyberspaceFlightPullWeakNorth,
		CyberspaceFlightPullWeakSouth,
		CyberspaceFlightPullMediumEast,
		CyberspaceFlightPullMediumWest,
		CyberspaceFlightPullMediumNorth,
		CyberspaceFlightPullMediumSouth,
		CyberspaceFlightPullStrongEast,
		CyberspaceFlightPullStrongWest,
		CyberspaceFlightPullStrongNorth,
		CyberspaceFlightPullStrongSouth,
		CyberspaceFlightPullMediumCeiling,
		CyberspaceFlightPullMediumFloor,
		CyberspaceFlightPullStrongCeiling,
		CyberspaceFlightPullStrongFloor,
	}
}
