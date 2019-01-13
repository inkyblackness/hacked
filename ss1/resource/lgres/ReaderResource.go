package lgres

import "github.com/inkyblackness/hacked/ss1/resource"

type readerResource struct {
	resource.Properties
	blockReader
}

func (rr readerResource) Compound() bool {
	return rr.Properties.Compound
}

func (rr readerResource) ContentType() resource.ContentType {
	return rr.Properties.ContentType
}

func (rr readerResource) Compressed() bool {
	return rr.Properties.Compressed
}
