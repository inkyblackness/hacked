package cmd

type stackEntry struct {
	link *stackEntry
	cmd  Command
}

// Stack describes a list of commands. The stack allows to sequentially
// undo and redo stacked commands.
// It essentially stores two lists: a list of commands to undo, and
// another of commands to redo.
// Modifying stack functions will panic if they are called while already in use.
type Stack struct {
	lockedBy string
	undoList *stackEntry
	redoList *stackEntry
}

// Perform executes the given command and puts it on the stack
// if the command was successful.
// This function also clears the list of commands to be redone.
func (stack *Stack) Perform(cmd Command) error {
	stack.lock("Perform")
	defer stack.unlock()

	err := cmd.Do()
	if err != nil {
		return err
	}
	stack.undoList = &stackEntry{stack.undoList, cmd}
	stack.redoList = nil
	return nil
}

// CanUndo returns true if there is at least one more command that can be undone.
func (stack *Stack) CanUndo() bool {
	return stack.undoList != nil
}

// Undo attempts to undo the previous command on the list.
// If there is no further command to undo, nothing happens.
// An error is returned if the command failed. In this case, the stack is
// unchanged and a further attempt to undo will try the same command again.
func (stack *Stack) Undo() error {
	stack.lock("Undo")
	defer stack.unlock()

	if stack.undoList == nil {
		return nil
	}
	entry := stack.undoList
	err := entry.cmd.Undo()
	if err != nil {
		return err
	}
	stack.undoList = entry.link
	entry.link = stack.redoList
	stack.redoList = entry
	return nil
}

// CanRedo returns true if there is at least one more command that can be redone.
func (stack *Stack) CanRedo() bool {
	return stack.redoList != nil
}

// Redo attempts to perform the next command on the redo list.
// If there is no further command to redo, nothing happens.
// An error is returned if the command failed. In this case, the stack is
// unchanged and a further attempt to redo will try the same command again.
func (stack *Stack) Redo() error {
	stack.lock("Redo")
	defer stack.unlock()

	if stack.redoList == nil {
		return nil
	}
	entry := stack.redoList
	err := entry.cmd.Do()
	if err != nil {
		return err
	}
	stack.redoList = entry.link
	entry.link = stack.undoList
	stack.undoList = entry
	return nil
}

func (stack *Stack) lock(by string) {
	if stack.lockedBy != "" {
		panic("Stack already in use by <" + stack.lockedBy + ">")
	}
	stack.lockedBy = by
}

func (stack *Stack) unlock() {
	stack.lockedBy = ""
}
