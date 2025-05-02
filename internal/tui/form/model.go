package form

import (
	"errors"
	"strings"
	"what-to-do/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type CommandForm struct {
	Command     string
	Alias       string
	Description string
	Spotlighted bool
	form        *huh.Form
	CommandID   int64 // Add this field to track which command is being updated
}

// NewCommandForm creates a new empty form
func NewCommandForm() CommandForm {
	return createForm("", "", "", false, 0)
}

// NewUpdateForm creates a form pre-filled with command data
func NewUpdateForm(cmd *domain.Command) CommandForm {
	return createForm(cmd.Command, cmd.Alias, cmd.Description, cmd.Spotlighted, cmd.ID)
}

// createForm is a helper to create a form with the given values
func createForm(command, alias, desc string, spotlight bool, id int64) CommandForm {
	var cmd = command
	var al = alias
	var description = desc
	var spotlighted = spotlight

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
				Value(&al).
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
				Value(&description),

			huh.NewConfirm().
				Key("spotlighted").
				Title("Spotlight this command?").
				Description("Spotlighted commands will appear in highlighted sections").
				Value(&spotlighted),
		),
	)

	return CommandForm{
		Command:     cmd,
		Alias:       al,
		Description: description,
		Spotlighted: spotlighted,
		form:        form,
		CommandID:   id,
	}
}

func (f CommandForm) Init() tea.Cmd {
	return f.form.Init()
}

func (f CommandForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)

	if newForm, ok := form.(*huh.Form); ok {
		f.form = newForm

		if f.form.State == huh.StateCompleted {
			f.Command = f.form.GetString("command")
			f.Alias = f.form.GetString("alias")
			f.Description = f.form.GetString("description")
			f.Spotlighted = f.form.GetBool("spotlighted")
		}
	}

	return f, cmd
}

func (f CommandForm) View() string {
	return f.form.View()
}

func (f CommandForm) IsCompleted() bool {
	return f.form.State == huh.StateCompleted
}

func (f CommandForm) IsCancelled() bool {
	return f.form.State == huh.StateAborted
}

// IsUpdating returns whether the form is for updating an existing command
func (f CommandForm) IsUpdating() bool {
	return f.CommandID > 0
}
