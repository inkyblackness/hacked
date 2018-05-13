package input

type shortcut struct {
	keyName  string
	modifier Modifier
	key      Key
}

var shortcuts = []shortcut{
	{"c", ModControl, KeyCopy},
	{"x", ModControl, KeyCut},
	{"v", ModControl, KeyPaste},
	{"z", ModControl, KeyUndo},
	{"Z", ModControl.With(ModShift), KeyRedo},
	{"z", ModControl.With(ModShift), KeyRedo},
	{"y", ModControl, KeyRedo}}

// ResolveShortcut tries to map the given name and modifier combination to a
// known (common) shortcut key. For instance, Ctrl+C is KeyCopy.
func ResolveShortcut(keyName string, modifier Modifier) (key Key, knownKey bool) {
	for _, entry := range shortcuts {
		if (entry.keyName == keyName) && (entry.modifier == modifier) {
			knownKey = true
			key = entry.key
		}
	}

	return
}
