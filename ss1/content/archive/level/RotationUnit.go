package level

// RotationUnit describes a rotation in the range 0..255.
type RotationUnit byte

// ToDegrees returns the value converted to degrees [0..360).
func (unit RotationUnit) ToDegrees() float64 {
	return (float64(unit) * 360.0) / 256.0
}
