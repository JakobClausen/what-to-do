package list

import (
	"what-to-do/internal/domain"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

// ExecuteCommandMsg is sent when a command should be executed
type ExecuteCommandMsg struct {
	Command string
}

// UpdateCommandMsg is sent when a command should be updated
type UpdateCommandMsg struct {
	Command *domain.Command
}

// DeleteCommandMsg is sent when a command should be deleted
type DeleteCommandMsg struct {
	ID    int64
	Title string
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func NewListModel(title string) ListModel {
	delegateKeys := newDelegateKeyMap()
	items := []list.Item{}

	delegate := newItemDelegate(delegateKeys)
	groceryList := list.New(items, delegate, 0, 0)

	styledTitle := lipgloss.JoinHorizontal(
		lipgloss.Left,
		titleStyle.Render(title),
	)

	groceryList.Title = styledTitle
	groceryList.Styles.Title = lipgloss.NewStyle()

	// Keep show help but disable other features
	groceryList.SetFilteringEnabled(false)
	groceryList.SetShowFilter(false)
	groceryList.SetShowHelp(true)

	// Disable the built-in help items we don't want
	groceryList.KeyMap.ShowFullHelp.SetEnabled(false)
	groceryList.KeyMap.Filter.SetEnabled(false)

	return ListModel{
		list:         groceryList,
		delegateKeys: delegateKeys,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListModel) View() string {
	return appStyle.Render(m.list.View())
}

// NewStatusMessage creates a new status message
func (m *ListModel) NewStatusMessage(message string) tea.Cmd {
	return m.list.NewStatusMessage(message)
}

// Toggle the spinner in the list
func (m *ListModel) ToggleSpinner() tea.Cmd {
	return m.list.ToggleSpinner()
}

// SetItems updates the list with commands from the database
func (m *ListModel) SetItems(commands []*domain.Command) {
	items := make([]list.Item, 0, len(commands))
	for _, cmd := range commands {
		items = append(items, item{command: cmd})
	}
	m.list.SetItems(items)
}

// AddItem adds a command to the list
func (m *ListModel) AddItem(cmd *domain.Command) {
	m.list.InsertItem(0, item{command: cmd})
}
