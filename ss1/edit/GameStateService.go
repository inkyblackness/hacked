package edit

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/citadel"
)

// VariableBaseContextIdentifier identifies which context for variable naming should be used.
type VariableBaseContextIdentifier int

// Context identifier are listed below.
const (
	VariableContextCitadel = VariableBaseContextIdentifier(0)
	VariableContextEngine  = VariableBaseContextIdentifier(1)
)

// VariableContextIdentifiers returns a list of all available identifiers.
func VariableContextIdentifiers() []VariableBaseContextIdentifier {
	return []VariableBaseContextIdentifier{VariableContextCitadel, VariableContextEngine}
}

// GameStateSettings describe the game state details.
type GameStateSettings struct {
	BaseContext int

	BooleanVariables []VariableSettings
	IntegerVariables []VariableSettings
}

// VariableSettings describe game variables.
type VariableSettings struct {
	Index int

	ResetValue  *int16 `json:",omitempty"`
	Name        string
	Description string
	Limits      *struct {
		Minimum int16
		Maximum int16
	} `json:",omitempty"`
	ValueNames []struct {
		Value int16
		Name  string
	} `json:",omitempty"`
}

func variableSettingsFrom(variables archive.GameVariables) []VariableSettings {
	settings := make([]VariableSettings, 0, len(variables))
	for varIndex, varInfo := range variables {
		varSettings := VariableSettings{
			Index:       varIndex,
			ResetValue:  varInfo.InitValue,
			Name:        varInfo.Name,
			Description: varInfo.Description,
			Limits:      nil,
			ValueNames:  nil,
		}
		if varInfo.Limits != nil {
			varSettings.Limits = &struct {
				Minimum int16
				Maximum int16
			}{
				Minimum: varInfo.Limits.Minimum,
				Maximum: varInfo.Limits.Maximum,
			}
		}
		for valIndex, valName := range varInfo.ValueNames {
			varSettings.ValueNames = append(varSettings.ValueNames,
				struct {
					Value int16
					Name  string
				}{
					Value: valIndex,
					Name:  valName,
				})
		}
		settings = append(settings, varSettings)
	}
	return settings
}

func variablesFrom(list []VariableSettings) archive.GameVariables {
	variables := make(archive.GameVariables)
	for _, settings := range list {
		varInfo := archive.GameVariableInfo{
			InitValue:   settings.ResetValue,
			Name:        settings.Name,
			Description: settings.Description,
		}
		if settings.Limits != nil {
			varInfo.Limits = &archive.GameVariableLimits{
				Minimum: settings.Limits.Maximum,
				Maximum: settings.Limits.Minimum,
			}
		}
		if len(settings.ValueNames) > 0 {
			varInfo.ValueNames = make(map[int16]string)
			for _, valName := range settings.ValueNames {
				varInfo.ValueNames[valName.Value] = valName.Name
			}
		}
		variables[settings.Index] = varInfo
	}
	return variables
}

// GameStateService handles all details on the global game state.
type GameStateService struct {
	registry cmd.Registry

	currentContext VariableBaseContextIdentifier

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

// CurrentSettings returns the snapshot of the project.
func (service GameStateService) CurrentSettings() GameStateSettings {
	return GameStateSettings{
		BaseContext:      int(service.currentContext),
		BooleanVariables: variableSettingsFrom(service.booleanVariables),
		IntegerVariables: variableSettingsFrom(service.integerVariables),
	}
}

// RestoreSettings loads the given snapshot.
func (service *GameStateService) RestoreSettings(settings GameStateSettings) {
	service.currentContext = VariableBaseContextIdentifier(settings.BaseContext)
	service.booleanVariables = variablesFrom(settings.BooleanVariables)
	service.integerVariables = variablesFrom(settings.IntegerVariables)
}

// VariableBaseContext returns the current context to use as basis for variables not specified by the project.
func (service GameStateService) VariableBaseContext() VariableBaseContextIdentifier {
	return service.currentContext
}

// SetVariableBaseContext sets the current context to use as basis for variables not specified by the project.
func (service *GameStateService) SetVariableBaseContext(identifier VariableBaseContextIdentifier) {
	service.currentContext = identifier
}

func (service GameStateService) baseInfoProvider() archive.GameVariableInfoProvider {
	switch service.currentContext {
	case VariableContextCitadel:
		return citadel.MissionVariables{}
	case VariableContextEngine:
		return archive.EngineVariables{}
	default:
		panic(fmt.Sprintf("invalid identifier: %d", service.currentContext))
	}
}

// DefaultAllVariables removes all project specific overrides.
func (service *GameStateService) DefaultAllVariables() error {
	return service.registry.Register(cmd.Named("RemoveAllOverrides"),
		cmd.Nested(func() error {
			for i := 0; i < archive.BooleanVarCount; i++ {
				err := service.DefaultBooleanVariable(i)
				if err != nil {
					return err
				}
			}
			for i := 0; i < archive.IntegerVarCount; i++ {
				err := service.DefaultIntegerVariable(i)
				if err != nil {
					return err
				}
			}
			return nil
		}))
}

// IntegerVariableOverride returns true if project-specific details are stored for given variable index.
func (service GameStateService) IntegerVariableOverride(index int) bool {
	_, perProject := service.integerVariables[index]
	return perProject
}

// IntegerVariable returns the variable info as per project settings for the given index.
func (service GameStateService) IntegerVariable(index int) archive.GameVariableInfo {
	baseInfo := service.baseInfoProvider().IntegerVariable(index)
	projInfo, perProject := service.integerVariables[index]
	if perProject {
		projInfo.Hardcoded = baseInfo.Hardcoded
		return projInfo
	}
	return baseInfo
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
	baseInfo := service.baseInfoProvider().BooleanVariable(index)
	projInfo, perProject := service.booleanVariables[index]
	if perProject {
		projInfo.Hardcoded = baseInfo.Hardcoded
		return projInfo
	}
	return baseInfo
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
