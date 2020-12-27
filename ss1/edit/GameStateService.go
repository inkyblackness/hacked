package edit

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
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
	registry cmd.Registry

	currentContext VariableContextIdentifier

	booleanVariables archive.GameVariables
	integerVariables archive.GameVariables
}

// NewGameStateService returns a new instance.
func NewGameStateService(registry cmd.Registry) *GameStateService {
	return &GameStateService{
		registry: registry,

		currentContext: VariableContextCitadel,

		booleanVariables: make(archive.GameVariables),
		integerVariables: make(archive.GameVariables),
	}
}

// VariableContext returns the current variable context.
func (service GameStateService) VariableContext() VariableContextIdentifier {
	return service.currentContext
}

// SetVariableContext sets the current variable context.
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

// IntegerVariableOverride returns true if project-specific details are stored for given variable index.
func (service GameStateService) IntegerVariableOverride(index int) bool {
	_, perProject := service.integerVariables[index]
	return perProject
}

// IntegerVariable returns the variable info as per project settings for the given index.
func (service GameStateService) IntegerVariable(index int) archive.GameVariableInfo {
	varInfo, perProject := service.integerVariables[index]
	if perProject {
		return varInfo
	}

	varInfoPtr := archive.EngineIntegerVariable(index)
	if varInfoPtr == nil {
		return archive.GameVariableInfoFor("(unused)").At(0)
	}
	return *varInfoPtr
}

// SetIntegerVariable sets the variable info for the given index in the project settings.
func (service *GameStateService) SetIntegerVariable(index int, info archive.GameVariableInfo) error {
	return service.registry.Register(cmd.Named("SetIntegerVariable"),
		cmd.Forward(setVariableTask(service.integerVariables, index, info)),
		cmd.Reverse(restoreVariableTask(service.integerVariables, index)))
}

// DefaultIntegerVariable clears the variable info for the given index in the project settings.
func (service *GameStateService) DefaultIntegerVariable(index int) error {
	return service.registry.Register(cmd.Named("DefaultIntegerVariable"),
		cmd.Forward(deleteVariableTask(service.integerVariables, index)),
		cmd.Reverse(restoreVariableTask(service.integerVariables, index)))
}

// BooleanVariableOverride returns true if project-specific details are stored for given variable index.
func (service GameStateService) BooleanVariableOverride(index int) bool {
	_, perProject := service.booleanVariables[index]
	return perProject
}

// BooleanVariable returns the variable info as per project settings for the given index.
func (service GameStateService) BooleanVariable(index int) archive.GameVariableInfo {
	varInfo, perProject := service.booleanVariables[index]
	if perProject {
		return varInfo
	}

	varInfoPtr := archive.EngineBooleanVariable(index)
	if varInfoPtr == nil {
		return archive.GameVariableInfoFor("(unused)").At(0)
	}
	return *varInfoPtr
}

// SetBooleanVariable sets the variable info for the given index in the project settings.
func (service *GameStateService) SetBooleanVariable(index int, info archive.GameVariableInfo) error {
	return service.registry.Register(cmd.Named("SetBooleanVariable"),
		cmd.Forward(setVariableTask(service.booleanVariables, index, info)),
		cmd.Reverse(restoreVariableTask(service.booleanVariables, index)))
}

// DefaultBooleanVariable clears the variable info for the given index in the project settings.
func (service *GameStateService) DefaultBooleanVariable(index int) error {
	return service.registry.Register(cmd.Named("DefaultBooleanVariable"),
		cmd.Forward(deleteVariableTask(service.booleanVariables, index)),
		cmd.Reverse(restoreVariableTask(service.booleanVariables, index)))
}

func setVariableTask(variables archive.GameVariables, index int, info archive.GameVariableInfo) cmd.Task {
	return func(modder world.Modder) error {
		variables[index] = info
		return nil
	}
}

func deleteVariableTask(variables archive.GameVariables, index int) cmd.Task {
	return func(modder world.Modder) error {
		delete(variables, index)
		return nil
	}
}

func restoreVariableTask(variables archive.GameVariables, index int) cmd.Task {
	info, present := variables[index]
	if present {
		return setVariableTask(variables, index, info)
	}
	return deleteVariableTask(variables, index)
}
