package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"test.terminal.shop/pkg/pages"
)

func main() {
	if _, err := tea.NewProgram(pages.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
