package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type MinWidthPage struct { }

func (w *MinWidthPage) Title() string { return "Minimum Width Required" }

func (w *MinWidthPage) Render(m *Model) string {

    title := m.theme.ActiveTitleForeground().
        AlignHorizontal(lipgloss.Center).
		Width(m.width).
        AlignVertical(lipgloss.Center).
		MarginBottom(2).
        Render("Minimum Width Required")

    desc := m.theme.ActiveDescForeground().
        AlignHorizontal(lipgloss.Center).
		Width(m.width).
        Render(fmt.Sprintf("You have %dx%d but we require a minimum of 100x30", m.width, m.height))

    return lipgloss.JoinVertical(0, title, desc)
}


