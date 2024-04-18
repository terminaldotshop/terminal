package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/terminalhq/terminal/go/pkg/api"
)

type CreditCardPage struct {
	card api.CreditCard
	form *huh.Form

	differentBillingAddress bool
}

func newCreditCardForm(differentBillingAddress *bool, creditCard *api.CreditCard) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name On CC").
				Value(&creditCard.Name).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(notEmpty("name")),
			huh.NewInput().
				Title("Credit Card").
				Value(&creditCard.Number).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(ccnValidator),
			huh.NewInput().
				Title("CVC").
				Value(&creditCard.CVC).
				Validate(compose(
					compose(
						notEmpty("CVC"),
						isDigits("CVC")),
					withinLen(3, 4, "CVC"),
				)),
			huh.NewInput().
				Title("ExpMonth").
				Value(&creditCard.ExpMonth).
				Validate(compose(
					compose(notEmpty("ExpMonth"), isDigits("ExpMonth")),
					mustBeLen(2, "ExpMonth"))),
			huh.NewInput().
				Title("ExpYear").
				Value(&creditCard.ExpYear).
				Validate(compose(
					compose(notEmpty("ExpYear"), isDigits("ExpYear")),
					mustBeLen(2, "ExpYear"))),
			huh.NewConfirm().
				Title("Is Shipping Address Different From Billing?").
				Value(differentBillingAddress),
		),
	)
}

func NewCreditCardPage() *CreditCardPage {
	creditCard := CreditCardPage{
		card: api.NewDefaultCreditCard(),
		form: nil,
	}

	return &creditCard
}

func (c *CreditCardPage) Exit(m Model) Model {
	m.creditCard = c.card
	return m
}

func (c *CreditCardPage) Enter(m Model) {
	c.card = m.creditCard
	c.form = newCreditCardForm(&c.differentBillingAddress, &c.card)
	c.form.Init()
}

func (c *CreditCardPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
		if c.form.State == huh.StateCompleted {
			m.creditCard = c.card
			m.differentBillingAddress = c.differentBillingAddress

			nextView := NewNavigateCCAddress
			if !c.differentBillingAddress {
				m.billingAddress = m.shippingAddress
				nextView = NewNavigateConfirm
			}
			return true, m, nextView
		}
		return true, m, cmd
	}

	return false, m, nil
}

func (s *CreditCardPage) Title() string { return "CreditCard" }

func (s *CreditCardPage) Render(m *Model) string {
	return RenderSplitView(*m, s.form.View(),
		lipgloss.JoinVertical(
			0,
			RenderEmail(*m),
			RenderShipping(m.theme, m.shippingAddress, "Shipping"),
		))
}
