package level

const (
	// ParametersSize is the serialized size, in bytes, of the Parameters structure.
	ParametersSize = 94
)

// Parameters describes level-global properties.
// It also contains information on the auto-maps, which, as per original code, should have been in the game state.
type Parameters struct {
	Size int16
	_    [2]byte

	CeilingHazardLevel   byte
	FloorHazardLevel     byte
	FloorHazardIsGravity byte
	FloorHazardOn        byte
	CeilingHazardOn      byte

	_ [4]byte
	_ [3][27]byte
}

// DefaultParameters returns a new instance of parameters.
func DefaultParameters() Parameters {
	return Parameters{
		Size:            ParametersSize,
		FloorHazardOn:   1,
		CeilingHazardOn: 1,
	}
}
