package object

import "fmt"

// PhysicsModel defines how an object can interact with the world.
type PhysicsModel byte

// String returns the textual representation of the value.
func (model PhysicsModel) String() string {
	if int(model) >= len(physicsModelNames) {
		return fmt.Sprintf("Unknown 0x%02X", int(model))
	}
	return physicsModelNames[model]
}

// PhysicsModel constants.
const (
	PhysicsModelInsubstantial PhysicsModel = 0
	PhysicsModelRegular       PhysicsModel = 1
	PhysicsModelStrange       PhysicsModel = 2
)

var physicsModelNames = []string{
	"Insubstantial",
	"Regular",
	"Strange",
}

// PhysicsModels returns all known constants.
func PhysicsModels() []PhysicsModel {
	return []PhysicsModel{
		PhysicsModelInsubstantial, PhysicsModelRegular, PhysicsModelStrange,
	}
}
