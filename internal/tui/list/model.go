package list

import (
	"what-to-do/internal/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

type ExecuteCommandMsg struct {
	Command string
}

type UpdateCommandMsg struct {
	Command *domain.Command
}

type DeleteCommandMsg struct {
	ID    int64
	Title string
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

	groceryList.SetFilteringEnabled(false)
	groceryList.SetShowFilter(false)
	groceryList.SetShowHelp(true)
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

func (m *ListModel) NewStatusMessage(message string) tea.Cmd {
	return m.list.NewStatusMessage(message)
}

func (m *ListModel) ToggleSpinner() tea.Cmd {
	return m.list.ToggleSpinner()
}

func (m *ListModel) SetItems(commands []*domain.Command) {
	items := make([]list.Item, 0, len(commands))
	for _, cmd := range commands {
		items = append(items, item{command: cmd})
	}
	m.list.SetItems(items)
}

func (m *ListModel) AddItem(cmd *domain.Command) {
	m.list.InsertItem(0, item{command: cmd})
}
