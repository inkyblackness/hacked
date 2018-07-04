package level_test

import (
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloorCeilingInfoPermutations(t *testing.T) {
	for height := level.TileHeightUnit(0); height < level.TileHeightUnitMax; height++ {
		for rotations := -4; rotations < 6; rotations++ {
			for _, hazard := range []bool{false, true} {
				floor1 := level.FloorInfo(0).WithAbsoluteHeight(height).WithTextureRotations(rotations).WithHazard(hazard)
				floor2 := level.FloorInfo(0).WithTextureRotations(rotations).WithHazard(hazard).WithAbsoluteHeight(height)
				floor3 := level.FloorInfo(0).WithHazard(hazard).WithAbsoluteHeight(height).WithTextureRotations(rotations)

				assert.Equal(t, floor1, floor2, "Floor1/2 mismatch")
				assert.Equal(t, floor1, floor3, "Floor1/3 mismatch")
				assert.Equal(t, floor2, floor3, "Floor2/3 mismatch")

				assert.Equal(t, height, floor1.AbsoluteHeight(), "floor.AbsoluteHeight mismatch")
				assert.Equal(t, (4+rotations)%4, floor1.TextureRotations(), "floor.TextureRotations mismatch")
				assert.Equal(t, hazard, floor1.HasHazard(), "floor.HasHazard mismatch")

				ceiling1 := level.CeilingInfo(0).WithAbsoluteHeight(height + 1).WithTextureRotations(rotations).WithHazard(hazard)
				ceiling2 := level.CeilingInfo(0).WithTextureRotations(rotations).WithHazard(hazard).WithAbsoluteHeight(height + 1)
				ceiling3 := level.CeilingInfo(0).WithHazard(hazard).WithAbsoluteHeight(height + 1).WithTextureRotations(rotations)

				assert.Equal(t, ceiling1, ceiling2, "Floor1/2 mismatch")
				assert.Equal(t, ceiling1, ceiling3, "Floor1/3 mismatch")
				assert.Equal(t, ceiling2, ceiling3, "Floor2/3 mismatch")

				assert.Equal(t, height+1, ceiling1.AbsoluteHeight(), "ceiling.AbsoluteHeight mismatch")
				assert.Equal(t, (4+rotations)%4, ceiling1.TextureRotations(), "ceiling.TextureRotations mismatch")
				assert.Equal(t, hazard, ceiling1.HasHazard(), "ceiling.HasHazard mismatch")
			}
		}
	}
}
