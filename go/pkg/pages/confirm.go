package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmPage struct{}

func NewConfirmPage() *ConfirmPage {
	return &ConfirmPage{}
}

func (c *ConfirmPage) Exit(m Model) Model { return m }
func (c *ConfirmPage) Enter(m Model)      {}

func (s *ConfirmPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	return false, m, nil
}

func (s *ConfirmPage) Title() string { return "Confirmation" }

func (s *ConfirmPage) Render(m *Model) string {
	container := lipgloss.NewStyle().
		Height(m.GetMaxPageHeight()).
		PaddingLeft(2)

	lines := []string{
		RenderEmail(*m),
		RenderShipping(m.theme, m.shippingAddress, "Shipping Address"),
		RenderCreditCard(m.theme, m.creditCard),
	}

	billingAddressTitle := "Billing Address"
	if m.differentBillingAddress {
		lines = append(lines, RenderShipping(m.theme, m.billingAddress, billingAddressTitle))
	} else {
		lines = append(lines, RenderSameShipping(m.theme, billingAddressTitle))
	}

	return container.Render(
		lipgloss.JoinVertical(
			0,
			lines...))
}
