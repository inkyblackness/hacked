package interpreters

type refinement struct {
	entry

	desc      *Description
	predicate Predicate
}
