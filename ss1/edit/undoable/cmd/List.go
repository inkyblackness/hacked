package cmd

// List is a sequence of commands, which are executed as one step.
type List []Command

// Do performs the entries in the list in ascending order.
// If an entry returns an error, the iteration is aborted and that error is returned.
func (list List) Do(trans Transaction) error {
	for _, entry := range list {
		err := entry.Do(trans)
		if err != nil {
			return err
		}
	}
	return nil
}

// Undo performs the entries in the list in descending order.
// If an entry returns an error, the iteration is aborted and that error is returned.
func (list List) Undo(trans Transaction) error {
	for i := len(list) - 1; i >= 0; i-- {
		err := list[i].Undo(trans)
		if err != nil {
			return err
		}
	}
	return nil
}
