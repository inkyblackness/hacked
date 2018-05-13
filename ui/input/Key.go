package input

// Key describes a named key on the keyboard. These are keys which are
// unspecific to layout or language, or are universal.
// Or, described in another way: keys that don't end up as printable characters.
type Key int

// Constants for commonly known named keys.
const (
	KeyEnter     = Key(300)
	KeyEscape    = Key(301)
	KeyBackspace = Key(302)
	KeyTab       = Key(303)

	KeyDown  = Key(310)
	KeyLeft  = Key(311)
	KeyRight = Key(312)
	KeyUp    = Key(313)

	KeyDelete   = Key(320)
	KeyEnd      = Key(321)
	KeyHome     = Key(322)
	KeyInsert   = Key(323)
	KeyPageDown = Key(324)
	KeyPageUp   = Key(325)

	KeyAlt     = Key(330)
	KeyControl = Key(331)
	KeyShift   = Key(332)
	KeySuper   = Key(333)

	KeyPause       = Key(340)
	KeyPrintScreen = Key(341)
	KeyCapsLock    = Key(342)
	KeyScrollLock  = Key(343)

	KeyF1  = Key(351)
	KeyF10 = Key(360)
	KeyF11 = Key(361)
	KeyF12 = Key(362)
	KeyF13 = Key(363)
	KeyF14 = Key(364)
	KeyF15 = Key(365)
	KeyF16 = Key(366)
	KeyF17 = Key(367)
	KeyF18 = Key(368)
	KeyF19 = Key(369)
	KeyF2  = Key(352)
	KeyF20 = Key(370)
	KeyF21 = Key(371)
	KeyF22 = Key(372)
	KeyF23 = Key(373)
	KeyF24 = Key(374)
	KeyF25 = Key(375)
	KeyF3  = Key(353)
	KeyF4  = Key(354)
	KeyF5  = Key(355)
	KeyF6  = Key(356)
	KeyF7  = Key(357)
	KeyF8  = Key(358)
	KeyF9  = Key(359)

	KeyCopy  = Key(380)
	KeyCut   = Key(381)
	KeyPaste = Key(382)

	KeyUndo = Key(390)
	KeyRedo = Key(391)
)

var keyToModifier = map[Key]Modifier{
	KeyShift:   ModShift,
	KeyControl: ModControl,
	KeyAlt:     ModAlt,
	KeySuper:   ModSuper}

// AsModifier returns the modifier equivalent for the key - if applicable.
func (key Key) AsModifier() Modifier {
	mod, isModifier := keyToModifier[key]

	if !isModifier {
		mod = ModNone
	}

	return mod
}
