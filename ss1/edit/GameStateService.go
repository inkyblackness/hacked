package edit

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/world/citadel"
)

// VariableContextIdentifier identifies which context for variable naming should be used.
type VariableContextIdentifier int

// Context identifier
const (
	VariableContextCitadel = VariableContextIdentifier(0)
	VariableContextEngine  = VariableContextIdentifier(1)
	VariableContextProject = VariableContextIdentifier(2)
)

// VariableContextIdentifiers returns a list of all available identifiers.
func VariableContextIdentifiers() []VariableContextIdentifier {
	return []VariableContextIdentifier{VariableContextCitadel, VariableContextEngine, VariableContextProject}
}

// GameStateService handles all details on the global game state.
type GameStateService struct {
	currentContext VariableContextIdentifier
}

// NewGameStateService returns a new instance.
func NewGameStateService() *GameStateService {
	return &GameStateService{
		currentContext: VariableContextCitadel,
	}
}

// VariableContext returns the current variable context.
func (service GameStateService) VariableContext() VariableContextIdentifier {
	return service.currentContext
}

// SEtVariableContext sets the current variable context.
func (service *GameStateService) SetVariableContext(identifier VariableContextIdentifier) {
	service.currentContext = identifier
}

// GameVariableInfoProviderFor returns the provider for the current identifier.
func (service GameStateService) GameVariableInfoProvider() archive.GameVariableInfoProvider {
	switch service.currentContext {
	case VariableContextCitadel:
		return citadel.MissionVariables{}
	case VariableContextEngine:
		return archive.EngineVariables{}
	case VariableContextProject:
		return service
	default:
		panic(fmt.Sprintf("invalid identifier: %d", service.currentContext))
	}
}

// IntegerVariable returns the variable info as per project settings for the given index.
func (service GameStateService) IntegerVariable(index int) archive.GameVariableInfo {
	varInfo := archive.EngineIntegerVariable(index)
	if varInfo == nil {
		return archive.GameVariableInfoFor("(unused)").At(0)
	}
	return *varInfo
}

// BooleanVariable returns the variable info as per project settings for the given index.
func (service GameStateService) BooleanVariable(index int) archive.GameVariableInfo {
	varInfo := archive.EngineBooleanVariable(index)
	if varInfo == nil {
		return archive.GameVariableInfoFor("(unused)").At(0)
	}
	return *varInfo
}
