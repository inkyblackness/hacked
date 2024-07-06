package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseBigStuff = interpreters.New()

var multiAnimation = interpreters.New().
	With("LoopType", 0, 2).As(interpreters.Bitfield(map[uint32]string{0x01: "Forward/Backward", 0x02: "Backward"})).
	With("Alternation", 2, 2).As(interpreters.SpecialValue("MultiAnimation")).
	With("Picture", 4, 2).As(interpreters.SpecialValue("PictureSource")).
	With("Alternate", 6, 2).As(interpreters.SpecialValue("PictureSource"))

var displayScenery = baseBigStuff.
	With("FrameCount", 0, 2).As(interpreters.RangedValue(0, 8)).
	Refining("", 2, 8, multiAnimation, interpreters.Always)

var displayControlPedestal = baseBigStuff.
	With("FrameCount", 0, 2).As(interpreters.RangedValue(0, 8)).
	With("TriggerObjectID1", 2, 2).As(interpreters.ObjectID()).
	With("TriggerObjectID2", 4, 2).As(interpreters.ObjectID()).
	Refining("", 2, 8, multiAnimation, interpreters.Always)

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
	With("TriggerObjectID1", 2, 2).As(interpreters.ObjectID()).
	With("TriggerObjectID2", 4, 2).As(interpreters.ObjectID())

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
		return forceColors[value]
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
	scienceSecurityEquipment.set(5, newInterpreterLeaf(buttonControlPedestal))
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
