package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	todoStore, err := CreateStore()

	if err != nil {
		log.Fatalf("Failed to initialize todo store: %v", err)
	}

	defer todoStore.Close()

	m := NewTable(todoStore)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
