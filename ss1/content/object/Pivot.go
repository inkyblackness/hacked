package object

// Pivot returns the assumed height offset for placing objects.
func Pivot(prop CommonProperties) float32 {
	const DefaultHeight = 1.0 / float32(0xbd00)
	const PhysicsScale = 96.0

	if (prop.RenderType == RenderTypeTextPoly) || (prop.RenderType == RenderTypeSpecial) {
		return 0.0
	} else if prop.PhysicsZ != 0 {
		return float32(prop.PhysicsZ) / PhysicsScale
	} else if prop.PhysicsXR != 0 {
		return float32(prop.PhysicsXR) / PhysicsScale
	}
	return DefaultHeight
}
