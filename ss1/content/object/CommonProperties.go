package object

const (
	// CommonPropertiesSize specifies, in bytes, the length a common properties structure has.
	CommonPropertiesSize = 27

	// HardnessLimit is the maximum hardness value.
	HardnessLimit = 253
	// PhysicsXRLimit is the maximum size value.
	PhysicsXRLimit = 253

	// DefenseNoCriticals is the defense value that prohibits criticals.
	DefenseNoCriticals = 0xFF
	// ToughnessNoDamage is the toughness value that prohibits regular damage.
	ToughnessNoDamage = 3
)

// CommonProperties are generic ones available for all objects.
type CommonProperties struct {
	Mass                   int32
	Hitpoints              int16
	Armor                  byte
	RenderType             RenderType
	PhysicsModel           PhysicsModel
	Hardness               byte
	_                      byte // pep: not used
	PhysicsXR              byte
	_                      byte // physics_y: not used
	PhysicsZ               byte
	Vulnerabilities        DamageTypeMask
	SpecialVulnerabilities SpecialDamageType
	_                      uint16 // remainder of "resistances": not used
	Defense                byte
	Toughness              byte
	Flags                  CommonFlagField
	MfdOrMeshID            uint16
	Bitmap3D               Bitmap3D
	DestroyEffect          DestroyEffect
}
