package levels

import "github.com/inkyblackness/hacked/ss1/content/archive/level"

// ObjectSelectionSetEvent notifies about the current set of selected objects.
type ObjectSelectionSetEvent struct {
	objects []level.ObjectID
}

// ObjectSelectionAddEvent notifies about added objects to the current selection.
type ObjectSelectionAddEvent struct {
	objects []level.ObjectID
}

// ObjectSelectionRemoveEvent notifies about removed objects from the current selection.
type ObjectSelectionRemoveEvent struct {
	objects []level.ObjectID
}
