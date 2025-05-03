package list

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Background(lipgloss.AdaptiveColor{Light: "#EAF6F2", Dark: "#1B312C"}).
				Padding(0, 1).
				MarginTop(1).
				Bold(true).
				Render

	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#FF4672"}).
				Background(lipgloss.AdaptiveColor{Light: "#FFF0F3", Dark: "#33262A"}).
				Padding(0, 1).
				MarginTop(1).
				Bold(true).
				Render
)
