package list

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#000000"}).
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

	// Base styles with consistent padding
	ItemTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#FFFFFF"}).
			Bold(true).
			Width(30).
			Padding(0, 1)

	ItemDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#AAAAAA"}).
			Italic(true).
			Padding(0, 1)

	// Selected styles with arrow
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#0066CC", Dark: "#5AA9FF"}).
				Bold(true).
				Width(32). // Increased width to accommodate arrow
				Padding(0, 1)

	// Description style - just normal style
	SelectedDescStyle = ItemDescStyle.Copy()
)
