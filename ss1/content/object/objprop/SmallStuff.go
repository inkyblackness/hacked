package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var cyberItems = interpreters.New().
	Refining("ColorScheme", 0, 6, cyberColorScheme, interpreters.Always)

func initSmallStuff() {
	objClass := object.Class(8)

	setSpecific(objClass, 5, cyberItems)
}
