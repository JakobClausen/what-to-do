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
func NewCompositeModel(width int) CompositeModel {
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
	// Update the list model.
	updatedListModel, listCmd := m.List.Update(msg)
	// Type assert the returned tea.Model to the concrete type (list.ListModel).
	m.List = updatedListModel.(list.ListModel)

	// Update the column model.
	updatedColumnModel, colCmd := m.Column.Update(msg)
	// Type assert the returned tea.Model to the concrete type (column.ColumnModel).
	m.Column = updatedColumnModel.(column.ColumnModel)

	// Combine both commands into one using tea.Batch.
	return m, tea.Batch(listCmd, colCmd)
}

// View combines the views from both models, e.g. side by side.
func (m CompositeModel) View() string {
	// Join the views horizontally with some separation.
	return lipgloss.JoinHorizontal(lipgloss.Top, m.List.View(), m.Column.View())
}
