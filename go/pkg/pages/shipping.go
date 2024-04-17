package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type ShippingPage struct {
    ShippingState
	form *huh.Form
}

type ShippingState struct {
	Name      string
	AddrLine1 string
	AddrLine2 string
	City      string
	State     string
	Zip       string
}

func newShippingForm(shipping *ShippingState) *huh.Form {
    return huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Name").
                Value(&shipping.Name).
                // Validating fields is easy. The form will mark erroneous fields
                // and display error messages accordingly.
                Validate(notEmpty("name")),
            huh.NewInput().
                Title("Address Line 1").
                Value(&shipping.AddrLine1).
                // Validating fields is easy. The form will mark erroneous fields
                // and display error messages accordingly.
                Validate(notEmpty("Shipping Address Line 1")),
            huh.NewInput().
                Title("Address Line 2").
                Value(&shipping.AddrLine2),
            huh.NewInput().
                Title("City").
                Value(&shipping.City).
                Validate(notEmpty("City")),
            huh.NewInput().
                Title("State").
                Value(&shipping.State).
                Validate(compose(notEmpty("State"), mustBeLen(2, "State"))),
            huh.NewInput().
                Title("Zip").
                Value(&shipping.Zip).
                Validate(notEmpty("Zip")),
            ),
        )
}

func NewShippingPage() *ShippingPage {
	shipping := ShippingPage{
        ShippingState: ShippingState{
            Name:      "",
            AddrLine1: "",
            AddrLine2: "",
            City:      "",
            State:     "",
            Zip:       "",
        },
        form:      nil,
	}

	return &shipping
}

func (s *ShippingPage) Exit(m Model) Model {
    m.shippingState = s.ShippingState
    return m
}

func (s *ShippingPage) Enter(m Model) {
    s.ShippingState = m.shippingState
    s.form = newShippingForm(&s.ShippingState)
    s.form.Init()
}

func (s *ShippingPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
    form, cmd := s.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		s.form = f
        if s.form.State == huh.StateCompleted {
            m.shippingState = s.ShippingState
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
