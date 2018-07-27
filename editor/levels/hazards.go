package levels

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

type hazardFormatter func(int) string

func rawHazardValue(int) string {
	return "raw: %d"
}

func lbpHazardValue(value int) string {
	return fmt.Sprintf("%3.1f LBP - raw: %%d", float32(value)/2)
}

func percentHazardValue(value int) string {
	return fmt.Sprintf("%v%%%% - raw: %%d", value*25)
}

type ceilingHazardInfo struct {
	title             string
	radiationRegister bool
	formatter         hazardFormatter
}

var ceilingHazards = []ceilingHazardInfo{
	{title: "Off", radiationRegister: false, formatter: rawHazardValue},
	{title: "Radiation", radiationRegister: true, formatter: lbpHazardValue},
}

func currentCeilingHazard(param *level.Parameters) ceilingHazardInfo {
	radiationRegister := param.RadiationRegister > 1
	if radiationRegister {
		return ceilingHazards[1]
	}
	return ceilingHazards[0]
}

type floorHazardInfo struct {
	title             string
	biohazardRegister bool
	isGravity         bool
	formatter         hazardFormatter
}

var floorHazards = []floorHazardInfo{
	{title: "Off", biohazardRegister: false, isGravity: false, formatter: rawHazardValue},
	{title: "Gravity", biohazardRegister: false, isGravity: true, formatter: percentHazardValue},
	{title: "Biohazard", biohazardRegister: true, isGravity: false, formatter: lbpHazardValue},
}

func currentFloorHazard(param *level.Parameters) floorHazardInfo {
	biohazardRegister := param.BiohazardRegister > 1
	isGravity := param.FloorHazardIsGravity != 0
	if isGravity {
		return floorHazards[1]
	} else if biohazardRegister {
		return floorHazards[2]
	}
	return floorHazards[0]
}
