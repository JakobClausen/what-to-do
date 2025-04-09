package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case tea.KeyCtrlQ.String(), tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyEnter.String():
			selectedTask := m.table.SelectedRow()[0]
			fmt.Println(selectedTask)
			return NewTextArea(selectedTask), nil
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		baseStyle.Render(m.table.View()),
		gap,
		"(esc to quit)",
	)
}

func NewTable(store *TodoStore) model {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Date", Width: 10},
		{Title: "Done", Width: 5},
	}

	dbRows, err := store.ListTodos()

	if err != nil {
		fmt.Println("Error fetching todos:", err)
		return model{}
	}

	rows := []table.Row{}

	for _, todo := range dbRows {

		if todo == nil {
			continue
		}

		var completedStr string
		if todo.Completed {
			completedStr = "O" // Completed is true
		} else {
			completedStr = "X" // Completed is false
		}

		row := table.Row{
			todo.Title,
			todo.Description,
			completedStr,
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{t}
}
