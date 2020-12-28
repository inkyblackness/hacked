package resource

import "fmt"

// NotFoundError indicates the unknown ID.
type NotFoundError struct {
	ID ID
}

// Error implements the error interface.
func (err NotFoundError) Error() string {
	return fmt.Sprintf("resource with ID %v does not exist", err.ID)
}

// ErrResourceNotFound returns an error specifying the given ID doesn't
// have an associated resource.
func ErrResourceNotFound(id ID) error {
	return NotFoundError{ID: id}
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

// ErrBlockNotFound returns an error specifying the given index does not identify a valid block.
func ErrBlockNotFound(index, available int) error {
	return BlockNotFoundError{Index: index, Available: available}
}
