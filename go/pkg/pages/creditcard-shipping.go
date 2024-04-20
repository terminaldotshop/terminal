package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminalhq/terminal/go/pkg/api"
)

type CreditCardAddress struct {
	address api.Address
	form    *huh.Form
}

func NewCreditCardAddress() *CreditCardAddress {
	creditCard := CreditCardAddress{
		address: api.NewAddress("", "", "", "", "", "", ""),
		form:    nil,
	}

	return &creditCard
}

func (c *CreditCardAddress) Exit(m Model) Model {
	m.billingAddress = c.address
	return m
}

func (c *CreditCardAddress) Enter(m Model) {
	c.address = m.billingAddress
	c.form = newShippingForm(&c.address)
	c.form.Init()
}

func (s *CreditCardAddress) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	form, cmd := s.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		s.form = f
		if s.form.State == huh.StateCompleted {
			cardparams := api.NewCardParams(m.creditCard, m.billingAddress)

			var err error
			m.stripeCardToken, err = api.StripeCreditCard(cardparams)
			if err != nil {
				return true, m, NewDialog(fmt.Sprintf("Credit Card Creation Failed: %s", err))
			}

			return true, m, NewNavigateConfirm
		}
		return true, m, cmd
	}

	return false, m, nil
}

func (s *CreditCardAddress) Title() string { return "CC - Address" }

func (s *CreditCardAddress) Render(m *Model) string {
	return RenderSplitView(*m, s.form.View(),
		lipgloss.JoinVertical(0,
			RenderEmail(*m),
			RenderShipping(m.theme, m.shippingAddress, "Shipping"),
			RenderCreditCard(m.theme, m.creditCard)),
	)
}
