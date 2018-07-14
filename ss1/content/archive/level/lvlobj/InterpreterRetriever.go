package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

type interpreterRetriever interface {
	specialize(key int) interpreterRetriever
	instance(data []byte) *interpreters.Instance
}

type interpreterLeaf struct {
	desc *interpreters.Description
}

func newInterpreterLeaf(desc *interpreters.Description) *interpreterLeaf {
	return &interpreterLeaf{desc: desc}
}

func (node *interpreterLeaf) specialize(key int) interpreterRetriever {
	return node
}

func (node *interpreterLeaf) instance(data []byte) *interpreters.Instance {
	return node.desc.For(data)
}

type interpreterEntry struct {
	defaultLeaf *interpreterLeaf
	subEntries  map[int]interpreterRetriever
}

func newInterpreterEntry(defaultDesc *interpreters.Description) *interpreterEntry {
	return &interpreterEntry{
		defaultLeaf: newInterpreterLeaf(defaultDesc),
		subEntries:  make(map[int]interpreterRetriever)}
}

func (node *interpreterEntry) set(key int, sub interpreterRetriever) {
	node.subEntries[key] = sub
}

func (node *interpreterEntry) specialize(key int) (retriever interpreterRetriever) {
	retriever = node.subEntries[key]
	if retriever == nil {
		retriever = node.defaultLeaf
	}
	return
}

func (node *interpreterEntry) instance(data []byte) *interpreters.Instance {
	return node.defaultLeaf.instance(data)
}
