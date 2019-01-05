package values

import (
	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// ObjectTypeControlRenderer renders a control to select an object.
type ObjectTypeControlRenderer struct {
	// Meta describes the layout of classes, subclasses, and types.
	Meta object.PropertiesTable
	// TextCache is used to retrieve the name of a type.
	TextCache *text.Cache
}

// Render creates the controls according to the given parameters.
func (renderer ObjectTypeControlRenderer) Render(readOnly, multiple bool, label string,
	classUnifier Unifier, typeUnifier Unifier,
	tripleResolver func(Unifier) object.Triple,
	changeHandler func(object.Triple)) {
	switch {
	case classUnifier.IsUnique():
		class := tripleResolver(classUnifier).Class
		triples := renderer.Meta.TriplesInClass(class)
		selectedIndex := -1
		if typeUnifier.IsUnique() {
			triple := tripleResolver(typeUnifier)
			for index, availableTriple := range triples {
				if availableTriple == triple {
					selectedIndex = index
				}
			}
			if selectedIndex < 0 {
				selectedIndex = len(triples)
				triples = append(triples, triple)
			}
		}
		RenderUnifiedCombo(readOnly, multiple, label, typeUnifier,
			func(u Unifier) int { return selectedIndex },
			func(value int) string {
				triple := triples[value]
				return renderer.tripleName(triple)
			},
			len(triples),
			func(newValue int) {
				triple := triples[newValue]
				changeHandler(triple)
			})
	case multiple:
		imgui.LabelText(label, "(multiple classes)")
	default:
		imgui.LabelText(label, "")
	}
}

func (renderer ObjectTypeControlRenderer) tripleName(triple object.Triple) string {
	suffix := "???"
	linearIndex := renderer.Meta.TripleIndex(triple)
	if (linearIndex >= 0) && (renderer.TextCache != nil) {
		key := resource.KeyOf(ids.ObjectLongNames, resource.LangDefault, linearIndex)
		objName, err := renderer.TextCache.Text(key)
		if err == nil {
			suffix = objName
		}
	}
	return triple.String() + ": " + suffix
}
