package form

import (
	"errors"
	"os"
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
	CommandID   int64
	form        *huh.Form
}

func NewCommandForm() CommandForm {
	return createForm("", "", "", false, 0)
}

func NewUpdateForm(cmd *domain.Command) CommandForm {
	return createForm(cmd.Command, cmd.Alias, cmd.Description, cmd.Spotlighted, cmd.ID)
}

func createForm(command, alias, desc string, spotlight bool, id int64) CommandForm {
	var cmd = command
	var al = alias
	var description = desc
	var spotlighted = spotlight

	editor := getPreferredEditor()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Key("command").
				Title("Bash Command").
				Placeholder("Enter a bash command").Editor(editor).
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

			huh.NewInput().
				Key("description").
				Title("Description (optional)").
				Placeholder("What does this command do?").
				CharLimit(200).
				Value(&description),

			huh.NewConfirm().
				Key("spotlighted").
				Title("Spotlight this command?").
				Value(&spotlighted),
		),
	)

	form = form.WithTheme(GetStyledForm())

	return CommandForm{
		Command:     cmd,
		Alias:       al,
		Description: description,
		Spotlighted: spotlighted,
		CommandID:   id,
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

func (f CommandForm) IsUpdating() bool {
	return f.CommandID > 0
}

func getPreferredEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Default to vim but fallback to nano if needed
		return "vim"
	}
	return editor
}
