package form

import (
	"errors"
	"strings"
	"what-to-do/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type CommandForm struct {
	CommandID   int64
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

	// Create a custom toggle component for spotlight selection
	spotlightOptions := []string{"Yes", "No"}
	initialSpotlight := 1 // Default to "No" (index 1)

	spotlightToggle := huh.NewSelect[int]().
		Key("spotlightIndex").
		Title("Spotlight this command?").
		Description("Spotlighted commands will appear in highlighted sections").
		Options(
			huh.NewOption(spotlightOptions[0], 0),
			huh.NewOption(spotlightOptions[1], 1),
		).
		Value(&initialSpotlight)

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

			spotlightToggle,
		),
	)

	return CommandForm{
		Command:     cmd,
		Alias:       alias,
		Description: desc,
		Spotlighted: false, // We'll set this when the form is completed
		form:        form,
	}
}

func (f *CommandForm) PopulateWithCommand(cmd *domain.Command) {
	f.Command = cmd.Command
	f.Alias = cmd.Alias
	f.Description = cmd.Description
	f.Spotlighted = cmd.Spotlighted
	f.CommandID = cmd.ID // Add this field to track the ID of command being updated
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

			// Convert the spotlight index to boolean
			spotlightIndex := f.form.GetInt("spotlightIndex")
			f.Spotlighted = spotlightIndex == 0 // 0 = Yes, 1 = No
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
