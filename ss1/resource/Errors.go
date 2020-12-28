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

// ErrNotFound returns an error specifying the given ID doesn't
// have an associated resource.
func ErrNotFound(id ID) error {
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

// WrongTypeError indicates a resource of unexpected type.
type WrongTypeError struct {
	Key      Key
	Expected ContentType
}

// Error implements the error interface.
func (err WrongTypeError) Error() string {
	return fmt.Sprintf("resource %v does not have expected type %v", err.Key, err.Expected)
}

// ErrWrongType returns an error specifying a resource was of unexpected type.
func ErrWrongType(key Key, expected ContentType) error {
	return WrongTypeError{Key: key, Expected: expected}
}
