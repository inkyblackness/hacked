package edit

import (
	"sort"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

type levelSelection struct {
	tiles   map[level.TilePosition]struct{}
	objects map[level.ObjectID]struct{}
}

func (selection levelSelection) tileList() []level.TilePosition {
	list := make([]level.TilePosition, 0, len(selection.tiles))
	for pos := range selection.tiles {
		list = append(list, pos)
	}
	sort.Slice(list, func(a, b int) bool {
		posA := list[a]
		posB := list[b]
		return (posA.Y < posB.Y) || ((posA.Y == posB.Y) && (posA.X < posB.X))
	})
	return list
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
	IsTileOnMap(levelIndex int, pos level.TilePosition) bool
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
	return *service.cleanSelection(service.currentLevel, selection)
}

func (service *LevelSelectionService) modifiableSelection() *levelSelection {
	selection, ok := service.levels[service.currentLevel]
	if !ok {
		selection = &levelSelection{
			tiles:   make(map[level.TilePosition]struct{}),
			objects: make(map[level.ObjectID]struct{}),
		}
		service.levels[service.currentLevel] = selection
	}
	return service.cleanSelection(service.currentLevel, selection)
}

func (service *LevelSelectionService) cleanSelection(levelIndex int, selection *levelSelection) *levelSelection {
	return service.cleanSelectionTiles(levelIndex, service.cleanSelectionObjects(levelIndex, selection))
}

func (service *LevelSelectionService) cleanSelectionTiles(levelIndex int, selection *levelSelection) *levelSelection {
	var toDelete []level.TilePosition
	for pos := range selection.tiles {
		if !service.provider.IsTileOnMap(levelIndex, pos) {
			toDelete = append(toDelete, pos)
		}
	}
	for _, pos := range toDelete {
		delete(selection.tiles, pos)
	}
	return selection
}

func (service *LevelSelectionService) cleanSelectionObjects(levelIndex int, selection *levelSelection) *levelSelection {
	var toDelete []level.ObjectID
	for id := range selection.objects {
		if !service.provider.IsObjectInUse(levelIndex, id) {
			toDelete = append(toDelete, id)
		}
	}
	for _, id := range toDelete {
		delete(selection.objects, id)
	}
	return selection
}

// NumberOfSelectedTiles returns the number of currently selected tiles in the current level.
func (service *LevelSelectionService) NumberOfSelectedTiles() int {
	return len(service.currentSelection().tiles)
}

// IsTileSelected returns true if the given tile is currently selected.
func (service *LevelSelectionService) IsTileSelected(pos level.TilePosition) bool {
	_, selected := service.currentSelection().tiles[pos]
	return selected
}

// CurrentSelectedTiles returns the list of currently selected tiles in the current level.
func (service *LevelSelectionService) CurrentSelectedTiles() []level.TilePosition {
	return service.currentSelection().tileList()
}

// SetCurrentSelectedTiles sets the list of currently selected tiles in the current level.
func (service *LevelSelectionService) SetCurrentSelectedTiles(list []level.TilePosition) {
	selection := service.modifiableSelection()
	for pos := range selection.tiles {
		delete(selection.tiles, pos)
	}
	for _, pos := range list {
		selection.tiles[pos] = struct{}{}
	}
}

// AddCurrentSelectedTiles adds the given tile IDs to the list of currently selected in the current level.
func (service *LevelSelectionService) AddCurrentSelectedTiles(list []level.TilePosition) {
	selection := service.modifiableSelection()
	for _, pos := range list {
		selection.tiles[pos] = struct{}{}
	}
}

// RemoveCurrentSelectedTiles removes the given tile IDs from the list of currently selected in the current level.
func (service *LevelSelectionService) RemoveCurrentSelectedTiles(list []level.TilePosition) {
	selection := service.modifiableSelection()
	for _, pos := range list {
		delete(selection.tiles, pos)
	}
}

// ToggleTileSelection toggles the selection of the given tile IDs in the current level.
func (service *LevelSelectionService) ToggleTileSelection(list []level.TilePosition) {
	selection := service.modifiableSelection()
	for _, pos := range list {
		if _, selected := selection.tiles[pos]; selected {
			delete(selection.tiles, pos)
		} else {
			selection.tiles[pos] = struct{}{}
		}
	}
}

// NumberOfSelectedObjects returns the number of currently selected objects in the current level.
func (service *LevelSelectionService) NumberOfSelectedObjects() int {
	return len(service.currentSelection().objects)
}

// CurrentSelectedObjects returns the list of currently selected objects in the current level.
func (service *LevelSelectionService) CurrentSelectedObjects() []level.ObjectID {
	return service.currentSelection().objectList()
}

// SetCurrentSelectedObjects sets the list of currently selected objects in the current level.
func (service *LevelSelectionService) SetCurrentSelectedObjects(ids []level.ObjectID) {
	selection := service.modifiableSelection()
	for id := range selection.objects {
		delete(selection.objects, id)
	}
	for _, id := range ids {
		selection.objects[id] = struct{}{}
	}
}

// AddCurrentSelectedObjects adds the given object IDs to the list of currently selected in the current level.
func (service *LevelSelectionService) AddCurrentSelectedObjects(ids []level.ObjectID) {
	selection := service.modifiableSelection()
	for _, id := range ids {
		selection.objects[id] = struct{}{}
	}
}

// RemoveCurrentSelectedObjects removes the given object IDs from the list of currently selected in the current level.
func (service *LevelSelectionService) RemoveCurrentSelectedObjects(ids []level.ObjectID) {
	selection := service.modifiableSelection()
	for _, id := range ids {
		delete(selection.objects, id)
	}
}

// ToggleObjectSelection toggles the selection of the given object IDs in the current level.
func (service *LevelSelectionService) ToggleObjectSelection(ids []level.ObjectID) {
	selection := service.modifiableSelection()
	for _, id := range ids {
		if _, selected := selection.objects[id]; selected {
			delete(selection.objects, id)
		} else {
			selection.objects[id] = struct{}{}
		}
	}
}
