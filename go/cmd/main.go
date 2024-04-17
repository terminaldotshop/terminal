package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/terminalhq/terminal/go/pkg/assert"
	"github.com/terminalhq/terminal/go/pkg/pages"
)

func main() {

    var screen string

    flag.StringVar(&screen, "screen", "", "sets the screen to a specific screen")
    flag.Parse()

    assert.Assert(screen == "" || screen == "email" || screen == "cc" || screen == "shipping", "invalid screen jump")

    f, err := tea.LogToFile("/tmp/out", "")
    if err != nil {
        fmt.Printf("Log to file errored: %+v\n", err)
    }
    log.SetOutput(f)
    defer f.Close()

	if _, err := tea.NewProgram(pages.NewModel(screen), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
