package resource

import "fmt"

// ResourceNotFoundError indicates the unknown ID.
type ResourceNotFoundError struct {
	ID ID
}

// Error implements the error interface.
func (err ResourceNotFoundError) Error() string {
	return fmt.Sprintf("resource with ID %v does not exist", err.ID)
}

// ErrResourceNotFound returns an error specifying the given ID doesn't
// have an associated resource.
func ErrResourceNotFound(id ID) error {
	return ResourceNotFoundError{ID: id}
}

// BlockNotFoundError indicates the unknown block index.
type BlockNotFoundError struct {
	Index     int
	Available int
}

// Error implements the error interface.
func (err BlockNotFoundError) Error() string {
	return fmt.Sprintf("block index wrong: %v/%v", err.Index, err.Available)
}

// ErrResourceNotFound returns an error specifying the given index does not identify a valid block.
func ErrBlockNotFound(index, available int) error {
	return BlockNotFoundError{Index: index, Available: available}
}
