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
    var dialog string

    flag.StringVar(&screen, "screen", "", "sets the screen to a specific screen")
    flag.StringVar(&dialog, "dialog", "", "adds a dialog to the model")
    flag.Parse()

    // TODO: Sucks
    assert.Assert(screen == "" ||
        screen == "email" ||
        screen == "cc" ||
        screen == "shipping" ||
        screen == "cc-addr" ||
        screen == "confirm", "invalid screen jump")

    f, err := tea.LogToFile("/tmp/out", "")
    if err != nil {
        fmt.Printf("Log to file errored: %+v\n", err)
    }
    log.SetOutput(f)
    defer f.Close()

    model := pages.NewModel(screen)
    if len(dialog) > 0 {
        model.Dialog = &dialog
    }

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
