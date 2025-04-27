package internal

import (
	"what-to-do/internal/store/sqlite"
	"what-to-do/internal/tui/column"
	"what-to-do/internal/tui/form"
	"what-to-do/internal/tui/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CompositeModel struct {
	Column   column.ColumnModel
	Lists    []list.ListModel
	Form     form.CommandForm
	ShowForm bool
}

func InitialModel(width int) CompositeModel {
	_, err := sqlite.New()

	if err != nil {
		panic(err)
	}

	spotlightList := list.NewListModel("Spotlight Items")
	allCommandsList := list.NewListModel("All Commands")

	return CompositeModel{
		Lists:    []list.ListModel{spotlightList, allCommandsList},
		Column:   column.CreateColumnModel(width),
		Form:     form.NewCommandForm(),
		ShowForm: false,
	}
}

func (m CompositeModel) Init() tea.Cmd {
	cmds := []tea.Cmd{m.Column.Init()}

	for i := range m.Lists {
		cmds = append(cmds, m.Lists[i].Init())
	}

	return tea.Batch(cmds...)
}

func (m CompositeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle showing form when "n" is pressed
	if keyMsg, ok := msg.(tea.KeyMsg); ok && !m.ShowForm {
		if keyMsg.String() == "n" {
			m.ShowForm = true
			m.Form = form.NewCommandForm() // Reset form
			return m, m.Form.Init()
		}
	}

	// If form is active, delegate updates to the form
	if m.ShowForm {
		updatedForm, cmd := m.Form.Update(msg)
		m.Form = updatedForm.(form.CommandForm)

		// Check if form is completed or cancelled
		if m.Form.IsCompleted() {
			// Process the completed form
			// Example: Save to database or add to lists
			if m.Form.Spotlighted {
				// Add to spotlight list
				// m.Lists[0].AddItem(m.Form)
			}
			// Add to all commands list
			// m.Lists[1].AddItem(m.Form)

			m.ShowForm = false
			return m, nil
		}

		if m.Form.IsCancelled() {
			m.ShowForm = false
			return m, nil
		}

		return m, cmd
	}

	// Normal update flow for non-form state
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Column = column.CreateColumnModel(msg.Width)
		columnHeight := lipgloss.Height(m.Column.View())

		listMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - columnHeight,
		}

		for i := range m.Lists {
			updatedListModel, listCmd := m.Lists[i].Update(listMsg)
			m.Lists[i] = updatedListModel.(list.ListModel)
			cmds = append(cmds, listCmd)
		}

		return m, tea.Batch(cmds...)
	}

	updatedColumnModel, colCmd := m.Column.Update(msg)
	m.Column = updatedColumnModel.(column.ColumnModel)
	cmds = append(cmds, colCmd)

	activeIdx := m.Column.ActiveTab()
	updatedListModel, listCmd := m.Lists[activeIdx].Update(msg)
	m.Lists[activeIdx] = updatedListModel.(list.ListModel)
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m CompositeModel) View() string {
	// Show form when active
	if m.ShowForm {
		return m.Form.View()
	}

	// Otherwise show the normal view
	activeIdx := m.Column.ActiveTab()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Column.View(),
		m.Lists[activeIdx].View(),
	)
}
