package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/terminalhq/terminal/go/pkg/api"
)

type ShippingPage struct {
	address api.Address
	form    *huh.Form
}

func newShippingForm(address *api.Address) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&address.Name).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(notEmpty("name")),
			huh.NewInput().
				Title("Address Line 1").
				Value(address.AddrLine1).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(notEmpty("Shipping Address Line 1")),
			huh.NewInput().
				Title("Address Line 2").
				Value(address.AddrLine2),
			huh.NewInput().
				Title("City").
				Value(address.City).
				Validate(notEmpty("City")),
			huh.NewInput().
				Title("State").
				Value(address.State).
				Validate(compose(notEmpty("State"), mustBeLen(2, "State"))),
			huh.NewInput().
				Title("Zip").
				Value(address.Zip).
				Validate(notEmpty("Zip")),
			huh.NewInput().
				Title("Country").
				Value(address.Country).
				Validate(notEmpty("Country")),
		),
	)
}

func NewShippingPage() *ShippingPage {
	shipping := ShippingPage{
		address: api.Address{},
		form:    nil,
	}

	return &shipping
}

func (s *ShippingPage) Exit(m Model) Model {
	m.shippingAddress = s.address
	return m
}

func (s *ShippingPage) Enter(m Model) {
	s.address = m.shippingAddress
	s.form = newShippingForm(&s.address)
	s.form.Init()
}

func (s *ShippingPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	form, cmd := s.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		s.form = f
		if s.form.State == huh.StateCompleted {
			m.shippingAddress = s.address
			return true, m, NewNavigateCC
		}
		return true, m, cmd
	}

	return false, m, nil
}

func (s *ShippingPage) Title() string { return "Shipping" }

func (s *ShippingPage) Render(m *Model) string {
	return RenderSplitView(*m, s.form.View(), RenderEmail(*m))
}
