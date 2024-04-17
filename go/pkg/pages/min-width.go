package pages

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

const MIN_WIDTH = 100
const MIN_HEIGHT = 28

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
        Render(fmt.Sprintf("You have %dx%d but we require a minimum of %dx%d", m.height, m.width, MIN_HEIGHT, MIN_WIDTH))

    return lipgloss.JoinVertical(0, title, desc)
}


