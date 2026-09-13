package main

import "fmt"

// taskItem adapts a Task for display in the Bubbles list component.
type taskItem struct {
	task Task
}

func (i taskItem) Title() string {
	return i.task.Title
}

func (i taskItem) Description() string {
	if i.task.Done {
		return fmt.Sprintf("#%d · done", i.task.ID)
	}

	return fmt.Sprintf("#%d", i.task.ID)
}

func (i taskItem) FilterValue() string {
	return i.task.Title
}
