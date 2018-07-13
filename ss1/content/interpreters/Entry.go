package interpreters

// FieldRange is a function specializing the range of a field to a simplifier.
type FieldRange func(*Simplifier) bool

type entry struct {
	start int
	count int
	via   FieldRange
}

func (e *entry) describe(simplifier *Simplifier) {
	if (e.via == nil) || !e.via(simplifier) {
		simplifier.rawValue(e)
	}
}
