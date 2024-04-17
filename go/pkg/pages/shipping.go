package pages

import tea "github.com/charmbracelet/bubbletea"

type ShippingPage struct{}

func (s *ShippingPage) Update(m Model, _ tea.Msg) (bool, tea.Model, tea.Cmd) {
	return false, m, nil
}

func (s *ShippingPage) Title() string { return "Shipping" }

func (s *ShippingPage) Render(m *Model) string {
	return "S H I P P I N G"
}
