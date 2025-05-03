package list

import (
	"what-to-do/internal/domain"
)

type item struct {
	command *domain.Command
}

func (i item) Title() string {
	if i.command == nil {
		return ""
	}
	return i.command.Alias
}

func (i item) ID() int64 {
	return i.command.ID
}

func (i item) Description() string {
	if i.command == nil {
		return ""
	}
	return ItemDescStyle.Render(i.command.Description)
}

func (i item) FilterValue() string {
	return i.command.Alias
}

func (i item) Command() *domain.Command {
	return i.command
}
