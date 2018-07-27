package input

import "strings"

type shortcut struct {
	keyName  string
	modifier Modifier
	key      Key
}

var shortcuts = []shortcut{
	{keyName: "c", modifier: ModControl, key: KeyCopy},
	{keyName: "x", modifier: ModControl, key: KeyCut},
	{keyName: "v", modifier: ModControl, key: KeyPaste},
	{keyName: "z", modifier: ModControl, key: KeyUndo},
	{keyName: "z", modifier: ModControl.With(ModShift), key: KeyRedo},
	{keyName: "y", modifier: ModControl, key: KeyRedo},
	{keyName: "s", modifier: ModControl, key: KeySave},
}

// ResolveShortcut tries to map the given name and modifier combination to a
// known (common) shortcut key. For instance, Ctrl+C is KeyCopy.
func ResolveShortcut(keyName string, modifier Modifier) (key Key, knownKey bool) {
	lowercaseName := strings.ToLower(keyName)
	for _, entry := range shortcuts {
		if (entry.keyName == lowercaseName) && (entry.modifier == modifier) {
			knownKey = true
			key = entry.key
		}
	}

	return
}
