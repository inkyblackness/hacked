package object

const (
	// CommonPropertiesSize specifies, in bytes, the length a common properties structure has.
	CommonPropertiesSize = 27
)

// CommonProperties are generic ones available for all objects.
type CommonProperties struct {
	Mass          int32
	Hitpoints     int16
	Armor         byte
	RenderType    RenderType
	PhysicsModel  PhysicsModel
	Hardness      byte
	Pep           byte
	PhysicsXR     byte
	PhysicsY      byte
	PhysicsZ      byte
	Resistances   uint32
	DefenseValue  byte
	Toughness     byte
	Flags         uint16
	MfdId         uint16
	Bitmap3D      uint16
	DestroyEffect DestroyEffect
}
