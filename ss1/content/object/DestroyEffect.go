package object

// DestroyEffect is a packed type describing effects during object destruction.
type DestroyEffect byte

const (
	destroyEffectValueMask     byte = 0x1F
	destroyEffectSoundMask     byte = 0x60
	destroyEffectExplosionMask byte = 0x80

	// DestroyEffectValueLimit is the maximum the value can take.
	DestroyEffectValueLimit = destroyEffectValueMask
)

// WithValue returns a new destroy effect with given value.
func (effect DestroyEffect) WithValue(value byte) DestroyEffect {
	return DestroyEffect((byte(effect) & ^destroyEffectValueMask) | value)
}

// Value returns the effect value.
func (effect DestroyEffect) Value() byte {
	return byte(effect) & destroyEffectValueMask
}

// WithSound returns a new destroy effect that has sound playback marked.
func (effect DestroyEffect) WithSound(value bool) DestroyEffect {
	soundValue := byte(0)
	if value {
		soundValue = 0x20
	}
	return DestroyEffect((byte(effect) & ^destroyEffectSoundMask) | soundValue)
}

// PlaySound returns true if a sound shall be played.
func (effect DestroyEffect) PlaySound() bool {
	return (byte(effect) & destroyEffectSoundMask) != 0
}

// WithExplosion returns a new destroy effect that has explosion display marked.
func (effect DestroyEffect) WithExplosion(value bool) DestroyEffect {
	explosionValue := byte(0)
	if value {
		explosionValue = 0x80
	}
	return DestroyEffect((byte(effect) & ^destroyEffectExplosionMask) | explosionValue)
}

// ShowExplosion returns true if an explosion shall be shown.
func (effect DestroyEffect) ShowExplosion() bool {
	return (byte(effect) & destroyEffectExplosionMask) != 0
}
