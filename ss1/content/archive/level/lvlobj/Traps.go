package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlobj/actions"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlobj/conditions"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseTraps = interpreters.New()

var repulsor = baseTraps.
	With("Ignored0000", 0, 1).As(interpreters.SpecialValue("Ignored")).
	With("StartHeightFraction", 10, 2).As(interpreters.RangedValue(0, 65536)).
	With("StartHeight", 12, 2).As(interpreters.RangedValue(0, 31)).
	With("EndHeightFraction", 14, 2).As(interpreters.RangedValue(0, 65536)).
	With("EndHeight", 16, 2).As(interpreters.RangedValue(0, 32)).
	With("Flags", 18, 4).As(interpreters.Bitfield(map[uint32]string{0x00000001: "Disabled", 0x00000008: "Strong"}))

var aiHint = baseTraps.
	With("Ignored0000", 0, 1).As(interpreters.SpecialValue("Ignored")).
	With("NextObjectIndex", 6, 2).As(interpreters.ObjectID()).
	With("TriggerObjectFlag", 18, 2).As(interpreters.EnumValue(map[uint32]string{0: "Off", 1: "On"})).
	With("TriggerObjectIndex", 20, 2).As(interpreters.ObjectID())

var baseTrigger = baseTraps.
	Refining("Action", 0, 22, actions.Unconditional(), interpreters.Always)

var gameVariableTrigger = baseTrigger.
	Refining("Condition", 2, 4, conditions.GameVariable(), interpreters.Always)

var puzzleData = interpreters.New()

var nullTrigger = baseTraps.
	Refining("Action", 0, 22, actions.Unconditional().
		Refining("PuzzleData", 6, 16, puzzleData, func(inst *interpreters.Instance) bool { return inst.Get("Type") == 0 }),
		interpreters.Always).
	Refining("Condition", 2, 4, conditions.GameVariable(),
		func(inst *interpreters.Instance) bool { return inst.Refined("Action").Get("Type") != 0 })

var deathWatchTrigger = baseTrigger.
	With("ConditionType", 5, 1).As(interpreters.EnumValue(map[uint32]string{0: "Object Type", 1: "Object ID"})).
	Refining("TypeCondition", 2, 4, conditions.ObjectType(), func(inst *interpreters.Instance) bool {
		return inst.Get("ConditionType") == 0
	}).
	Refining("IndexCondition", 2, 4, conditions.ObjectID(), func(inst *interpreters.Instance) bool {
		return inst.Get("ConditionType") == 1
	})

var ecologyTrigger = baseTrigger.
	Refining("TypeCondition", 2, 4, conditions.ObjectType(), interpreters.Always).
	With("ConditionLimit", 5, 1)

var mapNote = baseTraps.
	With("EntryOffset", 18, 4)

var musicVoodoo = baseTraps.
	With("MusicFlavour", 6, 1).As(interpreters.RangedValue(0, 4))

func initTraps() interpreterRetriever {

	gameVariableTriggers := newInterpreterLeaf(gameVariableTrigger)

	trigger := newInterpreterEntry(baseTraps)
	trigger.set(0, gameVariableTriggers) // tile entry trigger
	trigger.set(1, newInterpreterLeaf(nullTrigger))
	trigger.set(2, gameVariableTriggers) // floor trigger
	trigger.set(3, gameVariableTriggers) // player death trigger
	trigger.set(4, newInterpreterLeaf(deathWatchTrigger))
	trigger.set(7, newInterpreterLeaf(aiHint))
	trigger.set(8, gameVariableTriggers) // level entry trigger
	trigger.set(10, newInterpreterLeaf(repulsor))
	trigger.set(11, newInterpreterLeaf(ecologyTrigger))
	trigger.set(12, gameVariableTriggers) // shodan trigger

	mapMarker := newInterpreterEntry(baseTraps)
	mapMarker.set(3, newInterpreterLeaf(mapNote))
	mapMarker.set(4, newInterpreterLeaf(musicVoodoo))

	class := newInterpreterEntry(baseTraps)
	class.set(0, trigger)
	class.set(2, mapMarker)

	return class
}
