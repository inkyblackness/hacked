package level

const (
	// MapVersion is the version of the map structures.
	MapVersion int32 = 11
	// ObjectVersion is the version of object structures.
	ObjectVersion int32 = 27
	// SavefileVersion is the version of the save.
	SavefileVersion int32 = 13

	// InventorySize specifies how many items the inventory could hold.
	InventorySize int = 20

	// GradesOfShadow is the number of how much shadow (light) is possible.
	GradesOfShadow = 16
)

var mapVersionValue = MapVersion
var objectVersionValue = ObjectVersion
var savefileVersionValue = SavefileVersion
