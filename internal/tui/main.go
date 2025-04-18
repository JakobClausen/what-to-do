package tui

import (
	"what-to-do/internal/tui/column"
	"what-to-do/internal/tui/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CompositeModel struct {
	Column column.ColumnModel
	Lists  []list.ListModel
}

func InitialModel(width int) CompositeModel {

	spotlightList := list.NewListModel("Spotlight Items")

	allCommandsList := list.NewListModel("All Commands")

	return CompositeModel{
		Lists:  []list.ListModel{spotlightList, allCommandsList},
		Column: column.CreateColumnModel(width),
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

	activeIdx := m.Column.ActiveTab()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Column.View(),
		m.Lists[activeIdx].View(),
	)
}
