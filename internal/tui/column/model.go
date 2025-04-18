package column

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ColumnModel struct {
	activeTab int
	width     int
	Height    int
}

func (m ColumnModel) Init() tea.Cmd {
	return nil
}

func (m ColumnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.String() {
		case "1":
			m.activeTab = 0
		case "2":
			m.activeTab = 1
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ColumnModel) View() string {
	var doc strings.Builder

	var row string
	if m.activeTab == 0 {
		row = lipgloss.JoinHorizontal(
			lipgloss.Top,
			ActiveTab.Render("Spotlight"),
			Tab.Render("All cmds"),
		)
	} else {
		row = lipgloss.JoinHorizontal(
			lipgloss.Top,
			Tab.Render("Spotlight"),
			ActiveTab.Render("All cmds"),
		)
	}
	gap := TabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row)

	return DocStyle.MaxWidth(m.width).Render(doc.String())
}

func CreateColumnModel(width int) ColumnModel {

	return ColumnModel{
		activeTab: 0,
		width:     width,
	}
}
