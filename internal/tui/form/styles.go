package form

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func GetStyledForm() *huh.Theme {
	theme := huh.ThemeCharm()

	theme.Group = huh.GroupStyles{
		Base: lipgloss.NewStyle().
			Padding(1, 0),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EB5160")).
			Bold(true),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 0, 1, 2),
	}

	// Style the input fields
	theme.Focused.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EB5160")).
		Bold(true)

	theme.Blurred.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	theme.Blurred.UnselectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	return theme
}
