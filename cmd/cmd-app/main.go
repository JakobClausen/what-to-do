package main

import (
	"fmt"
	"os"
	tui "what-to-do/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func main() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	p := tea.NewProgram(tui.InitialModel(physicalWidth))
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
