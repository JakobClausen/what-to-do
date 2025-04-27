package form

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type CommandForm struct {
	Command     string
	Alias       string
	Description string
	Spotlighted bool
	form        *huh.Form
}

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
