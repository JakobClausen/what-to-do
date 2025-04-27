package list

import (
	"what-to-do/internal/domain"
)

// item represents a list item wrapping a domain.Command
type item struct {
	command *domain.Command
}

// Title returns the alias as the item title
func (i item) Title() string {
	if i.command == nil {
		return ""
	}
	return i.command.Alias
}

// Description returns the command description
func (i item) Description() string {
	if i.command == nil {
		return ""
	}
	return i.command.Description
}

// FilterValue returns the value to use for filtering
func (i item) FilterValue() string {
	return i.Title()
}

// Command returns the underlying command
func (i item) Command() *domain.Command {
	return i.command
}
