package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"what-to-do/internal/domain"
	"what-to-do/internal/store/sqlite"
	"what-to-do/internal/tui/column"
	"what-to-do/internal/tui/form"
	"what-to-do/internal/tui/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListRefreshMsg is sent when the lists need to be refreshed with new data
type ListRefreshMsg struct {
	Commands []*domain.Command
}

type CompositeModel struct {
	Column   column.ColumnModel
	Lists    []*list.ListModel
	Form     form.CommandForm
	ShowForm bool
	db       *sqlite.Store // Database reference
}

func InitialModel(width int) CompositeModel {
	db, err := sqlite.New()
	if err != nil {
		panic(err)
	}

	spotlightList := list.NewListModel("Spotlight Items")
	allCommandsList := list.NewListModel("All Commands")

	return CompositeModel{
		Lists:    []*list.ListModel{&spotlightList, &allCommandsList}, // Store pointers
		Column:   column.CreateColumnModel(width),
		Form:     form.NewCommandForm(),
		ShowForm: false,
		db:       db,
	}
}

func (m CompositeModel) Init() tea.Cmd {
	cmds := []tea.Cmd{m.Column.Init()}

	for i := range m.Lists {
		cmds = append(cmds, m.Lists[i].Init())
	}

	// Add an initial list refresh command
	cmds = append(cmds, m.refreshLists())

	return tea.Batch(cmds...)
}

func (m CompositeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle specific message types first
	switch msg := msg.(type) {
	case list.ExecuteCommandMsg:
		// Clear screen and quit app
		tea.ClearScreen()

		// Execute the command in the shell
		cmd := exec.Command("bash", "-c", msg.Command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// First quit the application
		return m, tea.Sequence(
			tea.Quit,
			func() tea.Msg {
				// Then execute the command
				cmd.Run()
				return nil
			},
		)
	case ListRefreshMsg:
		spotlightedCmds := []*domain.Command{}
		allCmds := []*domain.Command{}

		for _, cmd := range msg.Commands {
			allCmds = append(allCmds, cmd)
			if cmd.Spotlighted {
				spotlightedCmds = append(spotlightedCmds, cmd)
			}
		}

		// Update lists with the commands
		m.Lists[0].SetItems(spotlightedCmds) // Spotlight list
		m.Lists[1].SetItems(allCmds)         // All commands list

		return m, nil

	case tea.KeyMsg:
		if !m.ShowForm && msg.String() == "n" {
			m.ShowForm = true
			m.Form = form.NewCommandForm() // Reset form
			return m, m.Form.Init()
		}

	case tea.WindowSizeMsg:
		m.Column = column.CreateColumnModel(msg.Width)
		columnHeight := lipgloss.Height(m.Column.View())

		listMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - columnHeight,
		}

		for i := range m.Lists {
			updatedListModel, listCmd := m.Lists[i].Update(listMsg)
			if updatedModel, ok := updatedListModel.(*list.ListModel); ok {
				// No need to reassign since we're using pointers already
				// The model is updated in-place
				_ = updatedModel // Just to avoid unused variable warning
			}
			cmds = append(cmds, listCmd)
		}

		return m, tea.Batch(cmds...)
	}

	// Handle form updates when form is showing
	if m.ShowForm {
		updatedForm, cmd := m.Form.Update(msg)
		m.Form = updatedForm.(form.CommandForm)

		if m.Form.IsCompleted() {
			// Convert form data to domain.Command
			command := &domain.Command{
				Command:     m.Form.Command,
				Alias:       m.Form.Alias,
				Description: m.Form.Description,
				Spotlighted: m.Form.Spotlighted,
			}

			// Save to database
			ctx := context.Background()
			_, err := m.db.Create(ctx, command)
			if err != nil {
				// Handle error (you might want to display this to the user)
				fmt.Println("Failed to save command:", err)
			} else {
				// Command saved successfully, refresh lists
				cmds = append(cmds, m.refreshLists())
			}

			m.ShowForm = false
			return m, tea.Batch(cmds...)
		}

		if m.Form.IsCancelled() {
			m.ShowForm = false
			return m, nil
		}

		return m, cmd
	}

	// Handle regular column and list updates
	updatedColumnModel, colCmd := m.Column.Update(msg)
	m.Column = updatedColumnModel.(column.ColumnModel)
	cmds = append(cmds, colCmd)

	activeIdx := m.Column.ActiveTab()
	updatedListModel, listCmd := m.Lists[activeIdx].Update(msg)
	if updatedModel, ok := updatedListModel.(*list.ListModel); ok {
		// No need to reassign since we're using pointers
		_ = updatedModel
	}
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m CompositeModel) View() string {
	if m.ShowForm {
		return m.Form.View()
	}

	activeIdx := m.Column.ActiveTab()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Column.View(),
		m.Lists[activeIdx].View(),
	)
}

// refreshLists creates a command to refresh list data from the database
func (m CompositeModel) refreshLists() tea.Cmd {
	return func() tea.Msg {
		// Get all commands from database
		ctx := context.Background()
		commands, err := m.db.List(ctx)
		if err != nil {
			fmt.Println("Error loading commands:", err)
			return nil
		}

		// Return a message with the commands
		return ListRefreshMsg{Commands: commands}
	}
}
