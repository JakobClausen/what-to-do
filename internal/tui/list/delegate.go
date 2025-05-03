package list

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				// Get the selected item
				if i, ok := m.SelectedItem().(item); ok {
					// Get the actual command to execute
					cmd := i.Command()
					if cmd != nil {
						// Return a command that will send a message to run this command
						return func() tea.Msg {
							return ExecuteCommandMsg{Command: cmd.Command}
						}
					}
				}
				return nil

			case key.Matches(msg, keys.remove):
				item, ok := m.SelectedItem().(item)
				if !ok {
					return nil
				}

				title = item.Title()
				commandID := item.ID()

				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}

				// Return a command that will send a message to delete this from the database
				return tea.Sequence(
					m.NewStatusMessage(StatusMessageStyle(FormatStatusMessage("Deleting "+title+"...", true))),
					tea.Tick(time.Millisecond*800, func(time.Time) tea.Msg {
						return DeleteCommandMsg{
							ID:    commandID,
							Title: title,
						}
					}),
				)

			case key.Matches(msg, keys.update):
				// Get the selected item
				item, ok := m.SelectedItem().(item)
				if !ok {
					return nil
				}

				// Get the command to update
				cmd := item.Command()
				if cmd != nil {
					// Return a command that will send a message to update this command
					return func() tea.Msg {
						return UpdateCommandMsg{Command: cmd}
					}
				}
				return nil
			}
		}

		return nil
	}

	// Use the keys.ShortHelp() and keys.FullHelp() methods directly
	// This ensures we get all keys in the correct order
	d.ShortHelpFunc = func() []key.Binding {
		return keys.ShortHelp()
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return keys.FullHelp()
	}

	return d
}

func FormatStatusMessage(message string, success bool) string {
	icon := "✓"
	if !success {
		icon = "✗"
	}

	return fmt.Sprintf(" %s %s ", icon, message)
}

type delegateKeyMap struct {
	choose key.Binding
	add    key.Binding
	update key.Binding
	remove key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.add,
		d.update,
		d.remove,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.add,
			d.update,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "execute"),
		),
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		update: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
