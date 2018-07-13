package interpreters

// Predicate is a function that specifies whether a certain description
// is active within another.
type Predicate func(inst *Instance) bool

// Always is a predicate that returns true for every call.
func Always(inst *Instance) bool {
	return true
}

// Never is a predicate that returns false for every call.
func Never(inst *Instance) bool {
	return false
}

// Description is the meta information about a data interpreter. A description
// is an immutable object that provides new descriptions with any modification.
type Description struct {
	fields      map[string]*entry
	refinements map[string]*refinement
	lastField   *entry
}

var empty = newDescription()

// New returns an empty description.
func New() *Description {
	return empty
}

func newDescription() *Description {
	return &Description{
		fields:      make(map[string]*entry),
		refinements: make(map[string]*refinement)}
}

func (desc *Description) clone() *Description {
	cloned := newDescription()

	for key, e := range desc.fields {
		cloned.fields[key] = e
	}
	for key, r := range desc.refinements {
		cloned.refinements[key] = r
	}

	return cloned
}

// For returns an instance of an interpreter with the given data slice
// as backing buffer.
func (desc *Description) For(data []byte) *Instance {
	return &Instance{desc: desc, data: data}
}

// With extends the given description with a new field.
// The returned description is a new, separated object from the originating one.
func (desc *Description) With(key string, byteStart int, byteCount int) *Description {
	cloned := desc.clone()
	cloned.lastField = &entry{start: byteStart, count: byteCount}
	cloned.fields[key] = cloned.lastField

	return cloned
}

// As sets the value range for the previously added field.
func (desc *Description) As(fieldRange FieldRange) *Description {
	if desc.lastField == nil {
		panic("No field active")
	}
	desc.lastField.via = fieldRange
	return desc
}

// Refining adds another description within the given one. The provided predicate
// determines whether the refined description is active.
func (desc *Description) Refining(key string, byteStart int, byteCount int, refined *Description, predicate Predicate) *Description {
	cloned := desc.clone()
	cloned.refinements[key] = &refinement{
		entry:     entry{start: byteStart, count: byteCount},
		desc:      refined,
		predicate: predicate}

	return cloned
}
