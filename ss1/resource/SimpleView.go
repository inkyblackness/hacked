package resource

import (
	"io"
)

type simpleView struct {
	res *Resource
}

func (viewer simpleView) Compound() bool {
	return viewer.res.Compound
}

func (viewer simpleView) ContentType() ContentType {
	return viewer.res.ContentType
}

func (viewer simpleView) Compressed() bool {
	return viewer.res.Compressed
}

func (viewer simpleView) BlockCount() int {
	return viewer.res.BlockCount()
}

func (viewer simpleView) Block(index int) (io.Reader, error) {
	return viewer.res.Block(index)
}
