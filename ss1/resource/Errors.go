package resource

import "fmt"

// ErrResourceDoesNotExist returns an error specifying the given ID doesn't
// have an associated resource.
func ErrResourceDoesNotExist(id ID) error {
	return fmt.Errorf("resource with ID %v does not exist", id)
}
