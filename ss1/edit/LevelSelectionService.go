package edit

import (
	"sort"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

type levelSelection struct {
	//tiles []
	objects map[level.ObjectID]struct{}
}

func (selection levelSelection) objectList() []level.ObjectID {
	list := make([]level.ObjectID, 0, len(selection.objects))
	for id := range selection.objects {
		list = append(list, id)
	}
	sort.Slice(list, func(a, b int) bool { return list[a] < list[b] })
	return list
}

// LevelInfoProvider returns details on level content.
type LevelInfoProvider interface {
	IsLevelAvailable(index int) bool
	IsObjectInUse(levelIndex int, id level.ObjectID) bool
	IsTileOnMap(levelIndex int, x, y int) bool
}

// LevelSelectionService keeps track of selections related to levels.
type LevelSelectionService struct {
	provider LevelInfoProvider

	currentLevel int

	levels map[int]*levelSelection
}

// NewLevelSelectionService returns a new instance.
func NewLevelSelectionService(provider LevelInfoProvider) *LevelSelectionService {
	return &LevelSelectionService{
		provider:     provider,
		currentLevel: 0,
		levels:       make(map[int]*levelSelection),
	}
}

// CurrentLevelID returns the identifier for the currently selected level.
func (service *LevelSelectionService) CurrentLevelID() int {
	if !service.provider.IsLevelAvailable(service.currentLevel) {
		return 0
	}
	return service.currentLevel
}

// SetCurrentLevelID sets the currently selected level.
func (service *LevelSelectionService) SetCurrentLevelID(index int) {
	service.currentLevel = index
}

// SetCurrentLevelIDTask returns a command task that will set the current level.
func (service *LevelSelectionService) SetCurrentLevelIDTask(index int) cmd.Task {
	return func(modder world.Modder) error {
		service.SetCurrentLevelID(index)
		return nil
	}
}

func (service *LevelSelectionService) currentSelection() levelSelection {
	selection, ok := service.levels[service.currentLevel]
	if !ok {
		return levelSelection{}
	}
	return *selection
}

func (service *LevelSelectionService) modifiableSelection() *levelSelection {
	selection, ok := service.levels[service.currentLevel]
	if !ok {
		selection = &levelSelection{
			objects: make(map[level.ObjectID]struct{}),
		}
		service.levels[service.currentLevel] = selection
	}
	return selection
}

func (service *LevelSelectionService) CurrentSelectedObjects() []level.ObjectID {
	// TODO: filter via provider
	return service.currentSelection().objectList()
}

func (service *LevelSelectionService) SetCurrentSelectedObjects(ids []level.ObjectID) {
	selection := service.modifiableSelection()
	for _, id := range ids {
		selection.objects[id] = struct{}{}
	}
}
