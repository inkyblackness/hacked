package texts

import (
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// TextLinesAdapter handles simple text lines.
type TextLinesAdapter struct {
	mod *model.Mod

	cp        text.Codepage
	lineCache *text.LineCache
}

// NewTextLinesAdapter returns a new instance.
func NewTextLinesAdapter(mod *model.Mod) *TextLinesAdapter {
	cp := text.DefaultCodepage()
	adapter := &TextLinesAdapter{
		mod:       mod,
		cp:        cp,
		lineCache: text.NewLineCache(cp, mod),
	}
	return adapter
}

// InvalidateResources marks all identified resources as to-be newly loaded.
func (adapter *TextLinesAdapter) InvalidateResources(ids []resource.ID) {
	adapter.lineCache.InvalidateResources(ids)
}

// Codepage returns the codepage in use.
func (adapter *TextLinesAdapter) Codepage() text.Codepage {
	return adapter.cp
}

// Line retrieves the text line of given key. Returns an empty string if not found.
func (adapter *TextLinesAdapter) Line(key resource.Key) string {
	line, err := adapter.lineCache.Line(key)
	if err != nil {
		line = ""
	}
	return line
}
