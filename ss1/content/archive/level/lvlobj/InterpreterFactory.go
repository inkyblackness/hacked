package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// InterpreterFactory returns an interpreter instance according to a specific object triple.
type InterpreterFactory func(object.Triple, []byte) *interpreters.Instance
