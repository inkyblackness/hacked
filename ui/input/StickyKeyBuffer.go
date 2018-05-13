package input

// StickyKeyListener is the listener interface for receiving key events.
type StickyKeyListener interface {
	// Key is called for a typed key
	Key(key Key, modifier Modifier)
	// Modifier is called when the currently active modifier changed.
	Modifier(modifier Modifier)
}

// StickyKeyBuffer is a buffer to keep track of several identically named keys.
// Keys can be reported being pressed or released. Their state will be forwarded
// to a StickyKeyListener instance. If a specific key is reported to be pressed
// more than once, the listener will have received the down state only once.
type StickyKeyBuffer struct {
	pressedKeys     map[Key]int
	pressedModifier map[Modifier]int
	activeModifier  Modifier
	listener        StickyKeyListener
}

// NewStickyKeyBuffer returns a new instance of a sticky key buffer.
func NewStickyKeyBuffer(listener StickyKeyListener) *StickyKeyBuffer {
	buffer := &StickyKeyBuffer{
		pressedKeys:     make(map[Key]int),
		pressedModifier: make(map[Modifier]int),
		activeModifier:  ModNone,
		listener:        listener}

	return buffer
}

// ActiveModifier returns the currently pressed modifier set.
func (buffer *StickyKeyBuffer) ActiveModifier() Modifier {
	return buffer.activeModifier
}

// KeyDown registers a pressed key state. Multiple down states can be
// registered for the same key and result in only one key event.
func (buffer *StickyKeyBuffer) KeyDown(key Key, modifier Modifier) {
	keyAsModifier := key.AsModifier()

	buffer.setActiveModifier(modifier.With(keyAsModifier))
	if keyAsModifier == ModNone {
		oldCount := buffer.pressedKeys[key]

		buffer.pressedKeys[key] = oldCount + 1
		if oldCount == 0 {
			buffer.listener.Key(key, modifier)
		}
	} else {
		buffer.pressedModifier[keyAsModifier]++
	}
}

// KeyUp registers a released key state. Multiple up states can be registered
// for the same key, as long as enough down states were reported.
func (buffer *StickyKeyBuffer) KeyUp(key Key, modifier Modifier) {
	keyAsModifier := key.AsModifier()

	if keyAsModifier == ModNone {
		oldCount := buffer.pressedKeys[key]

		buffer.setActiveModifier(modifier)
		if oldCount > 0 {
			buffer.pressedKeys[key] = oldCount - 1
		}
	} else {
		oldCount := buffer.pressedModifier[keyAsModifier]

		if oldCount > 0 {
			buffer.pressedModifier[keyAsModifier] = oldCount - 1
			if oldCount == 1 {
				buffer.setActiveModifier(modifier.Without(keyAsModifier))
			}
		}
	}
}

// ReleaseAll notifies the listener of the reset of all modifiers.
// Key states are reset to accept new down states.
func (buffer *StickyKeyBuffer) ReleaseAll() {
	buffer.pressedKeys = make(map[Key]int)
	buffer.setActiveModifier(ModNone)
}

func (buffer *StickyKeyBuffer) setActiveModifier(modifier Modifier) {
	if buffer.activeModifier != modifier {
		for mod := range buffer.pressedModifier {
			if !modifier.Has(mod) {
				buffer.pressedModifier[mod] = 0
			}
		}
		buffer.activeModifier = modifier
		buffer.listener.Modifier(buffer.activeModifier)
	}
}
