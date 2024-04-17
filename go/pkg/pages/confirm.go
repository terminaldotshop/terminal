package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ConfirmPage struct { }

func NewConfirmPage() *ConfirmPage {
	return &ConfirmPage{ }
}

func (c *ConfirmPage) Exit(m Model) Model { return m }
func (c *ConfirmPage) Enter(m Model) { }

func (s *ConfirmPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	return false, m, nil
}

func (s *ConfirmPage) Title() string { return "Confirmation" }

func (s *ConfirmPage) Render(m *Model) string {
    return "DONE"
}

