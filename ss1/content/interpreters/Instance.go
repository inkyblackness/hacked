package interpreters

import (
	"bytes"
	"sort"
)

// Instance is one instantiated interpreter for a data block,
// based on a predetermined description.
type Instance struct {
	desc *Description
	data []byte
}

// Raw returns the raw array.
func (inst *Instance) Raw() []byte {
	return inst.data
}

// Undefined returns a byte array that specifies which bits were not
// specified. The array has the same length as the underlying buffer and
// has all bits set to 1 where no field or active refinement is present.
func (inst *Instance) Undefined() []byte {
	mask := bytes.Repeat([]byte{0xFF}, len(inst.data))

	for _, e := range inst.desc.fields {
		if inst.isValidRange(e) {
			copy(mask[e.start:e.start+e.count], bytes.Repeat([]byte{0x00}, e.count))
		}
	}
	for key, r := range inst.desc.refinements {
		if r.predicate(inst) {
			refined := inst.Refined(key)
			subMask := refined.Undefined()

			for index, m := range subMask {
				mask[r.start+index] &= m
			}
		}
	}

	return mask
}

// Keys returns an array of all registered keys, sorted by start index
func (inst *Instance) Keys() []string {
	return sortKeys(inst.desc.fields)
}

// ActiveRefinements returns an array of all keys where the corresponding refinement
// is active (The corresponding predicate returns true).
func (inst *Instance) ActiveRefinements() (keys []string) {
	entries := make(map[string]*entry)
	for key, r := range inst.desc.refinements {
		entries[key] = &r.entry
	}
	sortedKeys := sortKeys(entries)
	for _, key := range sortedKeys {
		if inst.desc.refinements[key].predicate(inst) {
			keys = append(keys, key)
		}
	}

	return
}

// Get returns the value associated with the given key. Should there be no
// value for the requested key, the function returns 0.
func (inst *Instance) Get(key string) uint32 {
	e := inst.desc.fields[key]
	value := uint32(0)

	if e != nil && inst.isValidRange(e) {
		for i := 0; i < e.count; i++ {
			value = (value << 8) | uint32(inst.data[e.start+e.count-1-i])
		}
	}

	return value
}

// Describe returns the description of a value key.
func (inst *Instance) Describe(key string, simplifier *Simplifier) {
	e := inst.desc.fields[key]

	if e != nil && inst.isValidRange(e) {
		e.describe(simplifier)
	}
}

// Set stores the provided value with the given key. Should there be no
// registration for the key, the function does nothing.
func (inst *Instance) Set(key string, value uint32) {
	e := inst.desc.fields[key]

	if e != nil && inst.isValidRange(e) {
		for i := 0; i < e.count; i++ {
			inst.data[e.start+i] = byte(value >> uint32(i*8))
		}
	}
}

// Refined returns an instance that was nested according to the description.
// This method returns an instance for registered refinements even if their predicate
// would not specify them being active.
// Should the refinement not exist, an instance without any fields will be returned.
func (inst *Instance) Refined(key string) (refined *Instance) {
	r := inst.desc.refinements[key]

	if r != nil && inst.isValidRange(&r.entry) {
		refined = &Instance{
			desc: r.desc,
			data: inst.data[r.start : r.start+r.count]}
	} else {
		refined = New().For(nil)
	}

	return
}

func (inst *Instance) isValidRange(e *entry) bool {
	available := len(inst.data)

	return (e.start + e.count) <= available
}

func sortKeys(entries map[string]*entry) []string {
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(a, b int) bool {
		entryA := entries[keys[a]]
		entryB := entries[keys[b]]
		return entryA.start < entryB.start
	})
	return keys
}
