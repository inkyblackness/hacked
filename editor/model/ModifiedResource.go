package model

import (
	"bytes"
	"fmt"
	"github.com/inkyblackness/hacked/ss1/resource"
	"io"
)

type modifiedResource struct {
	compound    bool
	contentType resource.ContentType
	compressed  bool
	blockCount  int
	blocks      map[int][]byte
}

func (res modifiedResource) Compound() bool {
	return res.compound
}

func (res modifiedResource) ContentType() resource.ContentType {
	return res.contentType
}

func (res modifiedResource) Compressed() bool {
	return res.compressed
}

func (res modifiedResource) BlockCount() int {
	return res.blockCount
}

func (res modifiedResource) Block(index int) (io.Reader, error) {
	if !res.isBlockIndexValid(index) {
		return nil, fmt.Errorf("block index wrong: %v/%v", index, res.blockCount)
	}
	return bytes.NewReader(res.blocks[index]), nil
}

func (res *modifiedResource) setBlocks(data [][]byte) {
	res.blocks = make(map[int][]byte)
	for index, blockData := range data {
		res.setBlock(index, blockData)
	}
}

func (res *modifiedResource) setBlock(index int, data []byte) {
	if index < 0 {
		return
	}
	if len(data) == 0 {
		res.delBlock(index)
	} else {
		res.blocks[index] = data
		if index >= res.blockCount {
			res.blockCount = index + 1
		}
	}
}

func (res *modifiedResource) delBlock(index int) {
	if !res.isBlockIndexValid(index) {
		return
	}
	delete(res.blocks, index)
	for (res.blockCount > 0) && !res.hasBlock(res.blockCount-1) {
		res.blockCount--
	}
}

func (res modifiedResource) isBlockIndexValid(index int) bool {
	return (index >= 0) && (index < res.blockCount)
}

func (res modifiedResource) hasBlock(index int) bool {
	_, existing := res.blocks[index]
	return existing
}
