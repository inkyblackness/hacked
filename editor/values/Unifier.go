package values

// Unifier can be used to verify whether a list of values has equal entries.
type Unifier struct {
	state unifierState
}

// NewUnifier returns a new instance.
func NewUnifier() Unifier {
	return Unifier{state: unifierInitState{}}
}

// UnifierFor returns a unifier for a single value.
func UnifierFor(value interface{}) Unifier {
	unifier := NewUnifier()
	unifier.Add(value)
	return unifier
}

// Add unifies the given value to the current state.
// The given value must be comparable.
func (u *Unifier) Add(value interface{}) {
	u.state = u.state.add(value)
}

// Unified returns the result of the unification.
// If all values that were added to the unifier were equal, then the first
// added value will be returned. Otherwise, nil will be returned.
func (u *Unifier) Unified() interface{} {
	return u.state.unified()
}

// IsUnique returns true if the unifier has received only equal values.
func (u Unifier) IsUnique() bool {
	return u.state.isUnique()
}

// IsEmpty returns true if the unifier has not received any value.
func (u Unifier) IsEmpty() bool {
	return (u.state == nil) || u.state.isEmpty()
}

type unifierState interface {
	add(value interface{}) unifierState
	unified() interface{}
	isUnique() bool
	isEmpty() bool
}

type unifierInitState struct{}

func (state unifierInitState) add(value interface{}) unifierState {
	return unifierMatchedState{value: value}
}

func (state unifierInitState) unified() interface{} {
	return nil
}

func (state unifierInitState) isUnique() bool {
	return false
}

func (state unifierInitState) isEmpty() bool {
	return true
}

type unifierMatchedState struct {
	value interface{}
}

func (state unifierMatchedState) add(value interface{}) unifierState {
	if state.value == value {
		return state
	}
	return unifierMismatchedState{}
}

func (state unifierMatchedState) unified() interface{} {
	return state.value
}

func (state unifierMatchedState) isUnique() bool {
	return true
}

func (state unifierMatchedState) isEmpty() bool {
	return false
}

type unifierMismatchedState struct{}

func (state unifierMismatchedState) add(value interface{}) unifierState {
	return state
}

func (state unifierMismatchedState) unified() interface{} {
	return nil
}

func (state unifierMismatchedState) isUnique() bool {
	return false
}

func (state unifierMismatchedState) isEmpty() bool {
	return false
}
