package tui

import (
	"what-to-do/internal/tui/column"
	"what-to-do/internal/tui/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CompositeModel holds both the list model and the column model.
type CompositeModel struct {
	List   list.ListModel
	Column column.ColumnModel
}

// NewCompositeModel initializes both models using the given width.
func InitialModel(width int) CompositeModel {
	return CompositeModel{
		List:   list.NewListModel(),
		Column: column.CreateColumnModel(width),
	}
}

// Init batches the initial commands from both submodels.
func (m CompositeModel) Init() tea.Cmd {
	return tea.Batch(m.List.Init(), m.Column.Init())
}

// Update passes the message to both submodels and updates their state.
func (m CompositeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// First update the column with the full width
		m.Column = column.CreateColumnModel(msg.Width)
		columnHeight := lipgloss.Height(m.Column.View())

		// Then update the list with adjusted height
		listMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - columnHeight, // -2 for padding/spacing
		}

		updatedListModel, listCmd := m.List.Update(listMsg)
		m.List = updatedListModel.(list.ListModel)
		cmds = append(cmds, listCmd)

		return m, tea.Batch(cmds...)
	}

	// Handle other message types as before
	updatedListModel, listCmd := m.List.Update(msg)
	m.List = updatedListModel.(list.ListModel)

	updatedColumnModel, colCmd := m.Column.Update(msg)
	m.Column = updatedColumnModel.(column.ColumnModel)

	return m, tea.Batch(listCmd, colCmd)
}

// View combines the views from both models, e.g. side by side.
func (m CompositeModel) View() string {

	// Join components with proper spacing
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Column.View(),
		m.List.View(),
	)
}
