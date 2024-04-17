package pages

import tea "github.com/charmbracelet/bubbletea"

type CreditCardPage struct{}

func (s *CreditCardPage) Update(m Model, _ tea.Msg) (bool, tea.Model, tea.Cmd) {
	return false, m, nil
}

func (s *CreditCardPage) Title() string { return "CreditCard" }

func (s *CreditCardPage) Render(m *Model) string {
	return "C R E D I T"
}

func NewCreditCardPage() *CreditCardPage {
	return &CreditCardPage{}
}
