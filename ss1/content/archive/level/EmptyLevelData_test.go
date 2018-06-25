package level_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"

	"github.com/stretchr/testify/assert"
)

func TestEmptyLevelDataSizesForDefault(t *testing.T) {
	param := level.EmptyLevelParameters{
		Cyberspace:  false,
		MapModifier: func(level.TileMap) {},
	}
	levelData := level.EmptyLevelData(param)

	// For now all the expected length values are deliberately written as numbers and not from the constants.
	// This is to better catch any changes against the default (vanilla) behaviour.

	assert.Equal(t, 4, len(levelData[lvlids.MapVersionNumber]), "MapVersionNumber")
	assert.Equal(t, 4, len(levelData[lvlids.ObjectVersionNumber]), "ObjectVersionNumber")
	assert.Equal(t, 58, len(levelData[lvlids.Information]), "Information")

	assert.Equal(t, 64*64*16, len(levelData[lvlids.TileMap]), "TileMap")
	assert.Equal(t, 8*64, len(levelData[lvlids.Schedules]), "Schedules")
	assert.Equal(t, 2*54, len(levelData[lvlids.TextureAtlas]), "TextureAtlas")
	assert.Equal(t, 27*872, len(levelData[lvlids.ObjectMasterTable]), "ObjectMasterTable")
	assert.Equal(t, 10*1600, len(levelData[lvlids.ObjectCrossRefTable]), "ObjectCrossRefTable")

	// leaving the object class tables

	assert.Equal(t, 4, len(levelData[lvlids.SavefileVersion]), "SavefileVersion")
	assert.Equal(t, 1, len(levelData[lvlids.Unused41]), "Unused41")

	assert.Equal(t, 7*4, len(levelData[lvlids.TextureAnimations]), "TextureAnimations")
	assert.Equal(t, 2*8, len(levelData[lvlids.SurveillanceSources]), "SurveillanceSources")
	assert.Equal(t, 2*8, len(levelData[lvlids.SurveillanceSurrogates]), "SurveillanceSurrogates")
	assert.Equal(t, 94, len(levelData[lvlids.Parameters]), "Parameters")
	assert.Equal(t, 0x0800, len(levelData[lvlids.MapNotes]), "MapNotes")
	assert.Equal(t, 4, len(levelData[lvlids.MapNotesPointer]), "MapNotesPointer")

	assert.Equal(t, 0x30, len(levelData[lvlids.Unknown48]), "Unknown48")
	assert.Equal(t, 0x1C0, len(levelData[lvlids.Unknown49]), "Unknown49")
	assert.Equal(t, 2, len(levelData[lvlids.Unknown50]), "Unknown50")

	assert.Equal(t, 15*64, len(levelData[lvlids.LoopConfiguration]), "LoopConfiguration")

	assert.Equal(t, 2, len(levelData[lvlids.Unknown52]), "Unknown52")

	assert.Equal(t, 0x40, len(levelData[lvlids.HeightSemaphores]), "HeightSemaphores")
}
