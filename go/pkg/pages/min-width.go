package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const MIN_WIDTH = 120
const MIN_HEIGHT = 30

type MinWidthPage struct{}

func (w *MinWidthPage) Update(m Model, _ tea.Msg) (bool, tea.Model, tea.Cmd) {
	return false, m, nil
}

func (w *MinWidthPage) Exit(m Model) Model { return m }
func (w *MinWidthPage) Enter(m Model)      {}
func (w *MinWidthPage) Title() string      { return "Minimum Width Required" }

func (w *MinWidthPage) Render(m *Model) string {

	height := m.GetMaxPageHeight()

	minWidthContainer := lipgloss.NewStyle().
		Height(height).
		AlignVertical(lipgloss.Center)

	title := m.theme.ActiveTitle().
		AlignHorizontal(lipgloss.Center).
		Width(m.width).
		AlignVertical(lipgloss.Center).
		MarginBottom(2).
		Render("Minimum Width Required")

	desc := m.theme.ActiveDescription().
		AlignHorizontal(lipgloss.Center).
		Width(m.width).
		Render(fmt.Sprintf("You have %dx%d but we require a minimum of %dx%d", m.height, m.width, MIN_HEIGHT, MIN_WIDTH))

	return minWidthContainer.Render(lipgloss.JoinVertical(0, title, desc))
}
