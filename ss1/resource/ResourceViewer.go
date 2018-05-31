package resource

import (
	"io"
)

type resourceViewer struct {
	res *Resource
}

func (viewer resourceViewer) Compound() bool {
	return viewer.res.Compound
}

func (viewer resourceViewer) ContentType() ContentType {
	return viewer.res.ContentType
}

func (viewer resourceViewer) Compressed() bool {
	return viewer.res.Compressed
}

func (viewer resourceViewer) BlockCount() int {
	return viewer.res.BlockCount()
}

func (viewer resourceViewer) Block(index int) (io.Reader, error) {
	return viewer.res.Block(index)
}
