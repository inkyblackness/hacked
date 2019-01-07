package values

import (
	"fmt"
	"sort"
	"strings"

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

	return simplifier
}
