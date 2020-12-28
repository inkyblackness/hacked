package values

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// StandardSimplifier returns a simplifier with common property controls.
func StandardSimplifier(readOnly bool, multiple bool, fullKey string, unifier Unifier,
	updater func(func(uint32) uint32), objTypeRenderer ObjectTypeControlRenderer) *interpreters.Simplifier {
	keys := strings.Split(fullKey, ".")
	key := keys[len(keys)-1]
	label := key + "###" + fullKey

	simplifier := interpreters.NewSimplifier(func(minValue, maxValue int64, formatter interpreters.RawValueFormatter) {
		RenderUnifiedSliderInt(readOnly, multiple, label, unifier,
			func(u Unifier) int {
				unifiedValue := u.Unified().(int32)
				if (minValue == -1) && (maxValue == 0x7FFF) {
					unifiedValue = int32(int16(unifiedValue))
				}
				return int(unifiedValue)
			},
			func(value int) string {
				result := formatter(value)
				if len(result) == 0 {
					result = "%d"
				} else {
					result += "  - raw: %d"
				}
				return result
			},
			int(minValue), int(maxValue),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})

	simplifier.SetEnumValueHandler(func(enumValues map[uint32]string) {
		valueKeys := make([]uint32, 0, len(enumValues))
		for valueKey := range enumValues {
			valueKeys = append(valueKeys, valueKey)
		}
		sort.Slice(valueKeys, func(indexA, indexB int) bool { return valueKeys[indexA] < valueKeys[indexB] })

		RenderUnifiedCombo(readOnly, multiple, label, unifier,
			func(u Unifier) int {
				unifiedValue := uint32(u.Unified().(int32))
				for index, valueKey := range valueKeys {
					if valueKey == unifiedValue {
						return index
					}
				}
				return -1
			},
			func(index int) string {
				if index < 0 {
					return ""
				}
				return enumValues[valueKeys[index]]
			},
			len(valueKeys),
			func(newIndex int) {
				updater(func(oldValue uint32) uint32 { return valueKeys[newIndex] })
			})
	})

	simplifier.SetBitfieldHandler(func(maskNames map[uint32]string) {
		masks := make([]uint32, 0, len(maskNames))
		for mask := range maskNames {
			masks = append(masks, mask)
		}
		sort.Slice(masks, func(indexA, indexB int) bool { return masks[indexA] < masks[indexB] })

		addMaskedItem := func(mask uint32) {
			maxValue := mask
			shift := 0
			maskedLabel := key + "." + maskNames[mask] + "###" + label + "-" + maskNames[mask]

			for (maxValue & 1) == 0 {
				shift++
				maxValue >>= 1
			}

			if maxValue == 1 {
				booleanUnifier := NewUnifier()
				if unifier.IsUnique() {
					booleanUnifier.Add((uint32(unifier.Unified().(int32)) & mask) != 0)
				}
				RenderUnifiedCheckboxCombo(readOnly, multiple, maskedLabel, booleanUnifier,
					func(newValue bool) {
						updater(func(oldValue uint32) uint32 {
							result := oldValue & ^mask
							if newValue {
								result |= mask
							}
							return result
						})
					})
			} else {
				RenderUnifiedSliderInt(readOnly, multiple, maskedLabel, unifier,
					func(u Unifier) int { return int((uint32(u.Unified().(int32)) & mask) >> uint32(shift)) },
					func(value int) string { return "%d" },
					0, int(maxValue),
					func(newValue int) {
						updater(func(oldValue uint32) uint32 {
							return (oldValue & ^mask) | (uint32(newValue) << uint32(shift))
						})
					})
			}
		}

		for _, mask := range masks {
			addMaskedItem(mask)
		}
	})

	simplifier.SetRotationHandler(func(minValue, maxValue int64) {
		calcValuePercent := func(value int) float64 {
			valueRange := (maxValue - minValue) + 1
			return float64(value) / float64(valueRange)
		}
		RenderUnifiedSliderInt(readOnly, multiple, label, unifier,
			func(u Unifier) int {
				unifiedValue := u.Unified().(int32)
				return int(unifiedValue)
			},
			func(value int) string {
				valuePercent := calcValuePercent(value)
				result := fmt.Sprintf("%3.02f°  - raw: %%d", valuePercent*360.0)
				return result
			},
			int(minValue), int(maxValue),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})

		if (imgui.IsItemFocused() || imgui.IsItemActive() || imgui.IsItemHovered()) && unifier.IsUnique() {
			valuePercent := calcValuePercent(int(unifier.Unified().(int32)))
			valueRadian := valuePercent * (math.Pi * 2)

			lineColor := imgui.Packed(color.RGBA{R: 0x21, G: 0xFF, B: 0x43, A: 0xFF}) // color copied from style.
			lineThickness := float32(4.0)
			imgui.BeginTooltip()
			imgui.Dummy(imgui.Vec2{X: imgui.TextLineHeightWithSpacing() * 2, Y: imgui.TextLineHeightWithSpacing() * 2})
			dl := imgui.WindowDrawList()
			winTopLeft := imgui.WindowPos()
			winSize := imgui.WindowSize()
			center := winTopLeft.Plus(winSize.Times(0.5))
			circleRadius := ((winSize.X / 2) * 5) / 6
			dl.AddCircleV(center, circleRadius, lineColor, 100, lineThickness)
			// This code is specialized for Hacker-Yaw. To re-use it for other rotations, the target needs to be provided.
			targetX := 1.0
			targetY := 0.0
			dl.AddLineV(center, imgui.Vec2{
				X: center.X + float32((float64(circleRadius)*targetX)*math.Cos(valueRadian)-(float64(circleRadius)*targetY)*math.Sin(valueRadian)),
				Y: center.Y - float32((float64(circleRadius)*targetY)*math.Cos(valueRadian)+(float64(circleRadius)*targetX)*math.Sin(valueRadian)),
			}, lineColor, lineThickness)
			imgui.EndTooltip()
		}
	})

	simplifier.SetSpecialHandler("ObjectTriple", func() {
		var classNames [object.ClassCount]string
		for index, class := range object.Classes() {
			classNames[index] = class.String()
		}
		tripleResolver := func(u Unifier) object.Triple { return object.TripleFromInt(int(u.Unified().(int32))) }
		RenderUnifiedCombo(readOnly, multiple, key+"-Class###"+fullKey+"-Class", unifier,
			func(u Unifier) int {
				triple := tripleResolver(u)
				return int(triple.Class)
			},
			func(value int) string { return fmt.Sprintf("%2d: %v", value, object.Class(value)) },
			object.ClassCount,
			func(newValue int) {
				triple := object.TripleFrom(newValue, 0, 0)
				updater(func(oldValue uint32) uint32 { return uint32(triple.Int()) })
			})

		objTypeRenderer.Render(readOnly, multiple, key+"###"+fullKey+"-Type", unifier, unifier,
			tripleResolver,
			func(newValue object.Triple) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue.Int()) })
			})
	})

	simplifier.SetSpecialHandler("Mistake", func() {})
	simplifier.SetSpecialHandler("Ignored", func() {})
	simplifier.SetSpecialHandler("Unknown", func() {})
	simplifier.SetSpecialHandler("Internal", func() {})

	return simplifier
}
