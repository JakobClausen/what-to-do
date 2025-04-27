package form

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// CommandForm holds the information about a bash command
type CommandForm struct {
	Command     string
	Alias       string
	Description string
	Spotlighted bool
	form        *huh.Form
}

// NewCommandForm creates a new command form
func NewCommandForm() CommandForm {
	var cmd string
	var alias string
	var desc string
	var spotlight bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("command").
				Title("Bash Command").
				Placeholder("Enter a bash command").
				Value(&cmd).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("Command cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Key("alias").
				Title("Command Alias").
				Placeholder("Enter a short alias for this command").
				Value(&alias).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("Alias cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Key("description").
				Title("Description (optional)").
				Placeholder("What does this command do?").
				CharLimit(200).
				Value(&desc),

			huh.NewConfirm().
				Key("spotlighted").
				Title("Spotlight this command?").
				Description("Spotlighted commands will appear in highlighted sections").
				Value(&spotlight),
		),
	)

	return CommandForm{
		Command:     cmd,
		Alias:       alias,
		Description: desc,
		Spotlighted: spotlight,
		form:        form,
	}
}

// Init initializes the form
func (f CommandForm) Init() tea.Cmd {
	return f.form.Init()
}

// Update handles form events
func (f CommandForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)

	// Update the form reference
	if newForm, ok := form.(*huh.Form); ok {
		f.form = newForm

		// When form completes, store the values
		if f.form.State == huh.StateCompleted {
			// Get the values directly without dereferencing
			f.Command = f.form.GetString("command")
			f.Alias = f.form.GetString("alias")
			f.Description = f.form.GetString("description")
			f.Spotlighted = f.form.GetBool("spotlighted")
		}
	}

	return f, cmd
}

// View renders the form
func (f CommandForm) View() string {
	return f.form.View()
}

// IsCompleted returns whether the form has completed
func (f CommandForm) IsCompleted() bool {
	return f.form.State == huh.StateCompleted
}

// IsCancelled returns whether the form was cancelled
func (f CommandForm) IsCancelled() bool {
	return f.form.State == huh.StateAborted
}
