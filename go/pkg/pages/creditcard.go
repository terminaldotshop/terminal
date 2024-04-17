package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type CreditCardPage struct {
	form *huh.Form
	info *CreditCardInfo
}

type CreditCardInfo struct{}

func (s *CreditCardPage) Update(m Model, _ tea.Msg) (bool, tea.Model, tea.Cmd) {
	return false, m, nil
}

func (s *CreditCardPage) Title() string          { return "CreditCard" }
func (s *CreditCardPage) Render(m *Model) string { return "C R E D I T" }

func NewCreditCardPage() *CreditCardPage {
	info := CreditCardInfo{}
	page := CreditCardPage{
		info: &info,
	}

	page.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Credit Card"),
			// Value(&email.email)
			// Validating fields is easy. The form will mark erroneous fields
			// and display error messages accordingly.
			// Validate(func(str string) error {
			// 	_, err := mail.ParseAddress(str)
			// 	if err != nil {
			// 		return errors.New("Not a valid email address")
			// 	}
			// 	return nil
			// }),
		),
	)
	return &CreditCardPage{}
}
