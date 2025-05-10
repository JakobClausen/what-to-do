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

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inancgumus/screen"
)

type ListRefreshMsg struct {
	Commands []*domain.Command
}

type CompositeModel struct {
	Column             column.ColumnModel
	Lists              []*list.ListModel
	Form               form.CommandForm
	ShowForm           bool
	db                 *sqlite.Store
	deletionInProgress bool
	deletingTitle      string
	isUpdating         bool
}

func InitialModel(width int) CompositeModel {
	db, err := sqlite.New()
	if err != nil {
		panic(err)
	}

	spotlightList := list.NewListModel()
	allCommandsList := list.NewListModel()

	return CompositeModel{
		Lists:    []*list.ListModel{&spotlightList, &allCommandsList},
		Column:   column.New(width),
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

	cmds = append(cmds, m.refreshLists())

	return tea.Batch(cmds...)
}

func (m CompositeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case list.ExecuteCommandMsg:
		commandToRun := msg.Command

		return m, func() tea.Msg {
			screen.Clear()
			screen.MoveTopLeft()

			tea.Quit()

			cmd := exec.Command("bash", "-c", commandToRun)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Command execution error: %v\n", err)
			}

			os.Exit(0)
			return nil
		}
		
	case list.CopyCommandMsg:
		commandToCopy := msg.Command
		
		err := clipboard.WriteAll(commandToCopy)
		if err != nil {
			activeTab := m.Column.ActiveTab()
			return m, m.Lists[activeTab].NewStatusMessage(
				list.ErrorMessageStyle(list.FormatStatusMessage("Error copying to clipboard: "+err.Error(), false)))
		}
		
		activeTab := m.Column.ActiveTab()
		return m, m.Lists[activeTab].NewStatusMessage(
			list.StatusMessageStyle(list.FormatStatusMessage("Copied to clipboard!", true)))

	case ListRefreshMsg:
		spotlightedCmds := []*domain.Command{}
		allCmds := []*domain.Command{}

		for _, cmd := range msg.Commands {
			allCmds = append(allCmds, cmd)
			if cmd.Spotlighted {
				spotlightedCmds = append(spotlightedCmds, cmd)
			}
		}

		m.Lists[0].SetItems(spotlightedCmds)
		m.Lists[1].SetItems(allCmds)

		if m.deletionInProgress && m.deletingTitle != "" {
			activeTab := m.Column.ActiveTab()
			cmds = append(cmds, m.Lists[activeTab].NewStatusMessage(
				list.StatusMessageStyle(list.FormatStatusMessage("Deleted "+m.deletingTitle, true))))

			m.deletionInProgress = false
			m.deletingTitle = ""
		}

		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Sequence(
				CleanupTerminal,
				tea.Quit,
			)
		case "a", "n":
			if !m.ShowForm {
				m.ShowForm = true
				m.Form = form.NewCommandForm()
				return m, m.Form.Init()
			}
		}

	case list.DeleteCommandMsg:
		if m.deletionInProgress {
			return m, nil
		}

		m.deletionInProgress = true
		m.deletingTitle = msg.Title

		ctx := context.Background()
		err := m.db.Delete(ctx, msg.ID)
		if err != nil {
			fmt.Printf("Error deleting command: %v\n", err)
			m.deletionInProgress = false

			activeTab := m.Column.ActiveTab()
			return m, m.Lists[activeTab].NewStatusMessage(
				list.ErrorMessageStyle(list.FormatStatusMessage("Error: "+err.Error(), false)))
		}

		return m, m.refreshLists()

	case tea.WindowSizeMsg:
		m.Column = column.New(msg.Width)
		columnHeight := lipgloss.Height(m.Column.View())

		listMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - columnHeight,
		}

		for i := range m.Lists {
			updatedListModel, listCmd := m.Lists[i].Update(listMsg)
			if updatedModel, ok := updatedListModel.(*list.ListModel); ok {
				_ = updatedModel
			}
			cmds = append(cmds, listCmd)
		}

		return m, tea.Batch(cmds...)

	case list.UpdateCommandMsg:
		m.ShowForm = true
		m.Form = form.NewUpdateForm(msg.Command)
		m.isUpdating = true
		return m, m.Form.Init()
	}

	if m.ShowForm {
		updatedForm, cmd := m.Form.Update(msg)
		m.Form = updatedForm.(form.CommandForm)

		if m.Form.IsCompleted() {
			command := &domain.Command{
				Command:     m.Form.Command,
				Alias:       m.Form.Alias,
				Description: m.Form.Description,
				Spotlighted: m.Form.Spotlighted,
			}

			ctx := context.Background()
			var err error

			if m.Form.IsUpdating() {
				command.ID = m.Form.CommandID
				err = m.db.Update(ctx, command)
				if err != nil {
					fmt.Println("Failed to update command:", err)
				} else {
					activeTab := m.Column.ActiveTab()
					cmds = append(cmds, m.Lists[activeTab].NewStatusMessage(
						list.StatusMessageStyle(list.FormatStatusMessage("Updated "+command.Alias, true))))
					cmds = append(cmds, m.refreshLists())
				}
			} else {
				_, err = m.db.Create(ctx, command)
				if err != nil {
					fmt.Println("Failed to create command:", err)
				} else {
					cmds = append(cmds, m.refreshLists())
				}
			}

			m.ShowForm = false
			m.isUpdating = false
			return m, tea.Batch(cmds...)
		}

		if m.Form.IsCancelled() {
			m.ShowForm = false
			return m, nil
		}

		return m, cmd
	}

	updatedColumnModel, colCmd := m.Column.Update(msg)
	m.Column = updatedColumnModel.(column.ColumnModel)
	cmds = append(cmds, colCmd)

	activeIdx := m.Column.ActiveTab()
	updatedListModel, listCmd := m.Lists[activeIdx].Update(msg)
	if updatedModel, ok := updatedListModel.(*list.ListModel); ok {
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

func (m CompositeModel) refreshLists() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		commands, err := m.db.List(ctx)
		if err != nil {
			fmt.Println("Error loading commands:", err)
			return nil
		}

		return ListRefreshMsg{Commands: commands}
	}
}

func CleanupTerminal() tea.Msg {
	fmt.Print("\033[?1049l")
	fmt.Print("\033[0m")
	screen.Clear()
	screen.MoveTopLeft()
	return nil
}
