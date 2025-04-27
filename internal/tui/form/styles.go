package form

import (
	"github.com/charmbracelet/huh"
)

// GetStyledForm returns the command form with custom styling applied
func GetStyledForm(data *CommandForm) *huh.Form {
	// Create the form
	form := data.form

	// Get the charm theme - this should be available in all versions
	theme := huh.ThemeCharm()

	// You can customize the theme here
	// theme.Focused.Title = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	// theme.Blurred.Title = lipgloss.NewStyle().Foreground(primaryColor.Faint()).Bold(true)

	// Apply the theme to the form
	form = form.WithTheme(theme)
	return form
}
