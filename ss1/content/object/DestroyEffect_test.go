package object_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/object"
)

func TestDestroyEffectValueRange(t *testing.T) {
	for _, showExplosion := range []bool{false, true} {
		for _, playSound := range []bool{false, true} {
			for value := byte(0); value <= object.DestroyEffectValueLimit; value++ {
				effect := object.DestroyEffect(0).WithExplosion(showExplosion).WithSound(playSound).WithValue(value)

				assert.Equal(t, showExplosion, effect.ShowExplosion(), "Explosion bit wrong")
				assert.Equal(t, playSound, effect.PlaySound(), "Sound value wrong")
				assert.Equal(t, value, effect.Value(), "Value wrong")
			}
		}
	}
}
