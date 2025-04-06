package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextAreaModel struct {
	textarea textarea.Model
	task     string // To store the selected task
}

func NewTextArea(task string) TextAreaModel {
	ta := textarea.New()
	ta.SetWidth(50)
	ta.SetHeight(10)
	ta.Focus()
	ta.Placeholder = "Add notes for your task..."
	ta.ShowLineNumbers = false

	return TextAreaModel{
		textarea: ta,
		task:     task,
	}
}

func (m TextAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TextAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m TextAreaModel) View() string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Task: "+m.task,
				m.textarea.View(),
			),
		)
}
