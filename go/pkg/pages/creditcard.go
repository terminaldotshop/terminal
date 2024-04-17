package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type CreditCardPage struct {
	CreditCardState
	form *huh.Form
}

type CreditCardState struct {
	Name string

	CC       string
	CVC      string
	ExpMonth string
	ExpYear  string

    // Store a shipping state?
}

func NewCreditCardPage() *CreditCardPage {
	creditCard := CreditCardPage{
		CreditCardState: CreditCardState{
			Name:     "",
			CC:       "",
			CVC:      "",
			ExpMonth: "",
			ExpYear:  "",
		},
		form: nil,
	}

	creditCard.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name On CC").
				Value(&creditCard.Name).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(notEmpty("name")),
			huh.NewInput().
				Title("Credit Card").
				Value(&creditCard.CC).
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
		),
	)

	creditCard.form.Init()

	return &creditCard
}

func (s *CreditCardPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	form, cmd := s.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		s.form = f
		if s.form.State == huh.StateCompleted {
			m.creditCardState = s.CreditCardState
			return true, m, NewStartCC
		}
		return true, m, cmd
	}

	return false, m, nil
}

func (s *CreditCardPage) Title() string { return "CreditCard" }

func (s *CreditCardPage) Render(m *Model) string {
	return s.form.View()
}
