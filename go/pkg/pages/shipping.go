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

    shipping.form = huh.NewForm(
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

    shipping.form.Init()

	return &shipping
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
    return s.form.View()
}
