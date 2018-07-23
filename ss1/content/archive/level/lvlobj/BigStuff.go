package lvlobj

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseBigStuff = interpreters.New()

var displayScenery = baseBigStuff.
	With("FrameCount", 0, 2).As(interpreters.RangedValue(0, 4)).
	With("LoopType", 2, 2).As(interpreters.EnumValue(map[uint32]string{0: "Forward", 1: "Forward/Backward", 2: "Backward", 3: "Forward/Backward"})).
	With("AlternationType", 4, 2).As(interpreters.EnumValue(map[uint32]string{0: "Don't Alternate", 3: "Alternate Randomly"})).
	With("PictureSource", 6, 2).As(interpreters.RangedValue(0, 0x01FF)).
	With("AlternateSource", 8, 2).As(interpreters.RangedValue(0, 0x01FF))

var displayControlPedestal = baseBigStuff.
	With("FrameCount", 0, 2).As(interpreters.RangedValue(0, 4)).
	With("TriggerObjectID", 2, 2).As(interpreters.ObjectID()).
	With("AlternationType", 4, 2).As(interpreters.EnumValue(map[uint32]string{0: "Don't Alternate", 3: "Alternate Randomly"})).
	With("PictureSource", 6, 2).As(interpreters.RangedValue(0, 0x01FF)).
	With("AlternateSource", 8, 2).As(interpreters.RangedValue(0, 0x01FF))

var cabinetFurniture = baseBigStuff.
	With("Object1ID", 2, 2).As(interpreters.ObjectID()).
	With("Object2ID", 4, 2).As(interpreters.ObjectID())

var texturableFurniture = baseBigStuff.
	With("TextureIndex", 6, 2).As(interpreters.RangedValue(0, 500))

var wordScenery = baseBigStuff.
	With("TextIndex", 0, 2).As(interpreters.RangedValue(0, 511)).
	With("Font", 2, 1).As(interpreters.Bitfield(map[uint32]string{
	0x0F: "Face",
	0xF0: "Size"})).
	With("Color", 4, 1).As(interpreters.RangedValue(0, 255))

var textureMapScenery = baseBigStuff.
	With("TextureIndex", 6, 2).As(interpreters.SpecialValue("LevelTexture"))

var buttonControlPedestal = baseBigStuff.
	With("TriggerObjectID", 2, 2).As(interpreters.ObjectID())

var surgicalMachine = baseBigStuff.
	With("BrokenState", 2, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "OK", 0xE7: "Broken"})).
	With("BrokenMessageIndex", 5, 1)

var securityCamera = baseBigStuff.
	With("PanningSwitch", 2, 1).As(interpreters.EnumValue(map[uint32]string{0: "Stationary", 1: "Panning"}))

var solidBridge = baseBigStuff.
	With("Size", 2, 1).As(interpreters.Bitfield(map[uint32]string{0x0F: "X", 0xF0: "Y"})).
	With("Height", 3, 1).As(interpreters.SpecialValue("ObjectHeight")).
	With("TopBottomTexture", 4, 1).As(interpreters.SpecialValue("MaterialOrLevelTexture")).
	With("SideTexture", 5, 1).As(interpreters.SpecialValue("MaterialOrLevelTexture"))

var forceBridge = baseBigStuff.
	With("Size", 2, 1).As(interpreters.Bitfield(map[uint32]string{0x0F: "X", 0xF0: "Y"})).
	With("Height", 3, 1).As(interpreters.SpecialValue("ObjectHeight")).
	With("Color", 6, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) (result string) {
		if colorText, defined := forceColors[value]; defined {
			result = fmt.Sprintf("%s", colorText)
		} else {
			result = ""
		}
		return result
	}))

func initBigStuff() interpreterRetriever {
	displays := newInterpreterLeaf(displayScenery)
	textureable := newInterpreterLeaf(texturableFurniture)

	electronics := newInterpreterEntry(baseBigStuff)
	electronics.set(6, displays)
	electronics.set(7, displays)

	furniture := newInterpreterEntry(baseBigStuff)
	furniture.set(2, newInterpreterLeaf(cabinetFurniture))
	furniture.set(5, textureable)
	furniture.set(7, textureable)
	furniture.set(8, textureable)

	surfaces := newInterpreterEntry(baseBigStuff)
	surfaces.set(3, newInterpreterLeaf(wordScenery))
	surfaces.set(6, displays)
	surfaces.set(7, newInterpreterLeaf(textureMapScenery))
	surfaces.set(8, displays)
	surfaces.set(9, displays)

	lighting := newInterpreterEntry(baseBigStuff)

	medicalEquipment := newInterpreterEntry(baseBigStuff)
	medicalEquipment.set(0, newInterpreterLeaf(buttonControlPedestal))
	medicalEquipment.set(3, newInterpreterLeaf(surgicalMachine))
	medicalEquipment.set(5, textureable)

	scienceSecurityEquipment := newInterpreterEntry(baseBigStuff)
	scienceSecurityEquipment.set(4, newInterpreterLeaf(securityCamera))
	scienceSecurityEquipment.set(6, newInterpreterLeaf(displayControlPedestal))

	gardenScenery := newInterpreterEntry(baseBigStuff)

	bridges := newInterpreterEntry(baseBigStuff)
	solidBridges := newInterpreterLeaf(solidBridge)
	forceBridges := newInterpreterLeaf(forceBridge)
	bridges.set(0, solidBridges)
	bridges.set(1, solidBridges)
	bridges.set(7, forceBridges)
	bridges.set(8, forceBridges)
	bridges.set(9, forceBridges)

	class := newInterpreterEntry(baseBigStuff)
	class.set(0, electronics)
	class.set(1, furniture)
	class.set(2, surfaces)
	class.set(3, lighting)
	class.set(4, medicalEquipment)
	class.set(5, scienceSecurityEquipment)
	class.set(6, gardenScenery)
	class.set(7, bridges)

	return class
}
