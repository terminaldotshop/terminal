package pages

import (
	"errors"
	"net/mail"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type EmailPage struct {
    form *huh.Form
    email string
}

func newEmailForm(email *EmailPage) *huh.Form {
    return huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Email").
                Value(&email.email).
                // Validating fields is easy. The form will mark erroneous fields
                // and display error messages accordingly.
                Validate(func(str string) error {
                    _, err := mail.ParseAddress(str)
                    if err != nil {
                        return errors.New("Not a valid email address")
                    }
                    return nil
                }),
            ),
        )
}

func NewEmailPage() *EmailPage {
    email := EmailPage{
        form: nil,
        email: "",
    }

    return &email
}

func (w *EmailPage) Exit(m Model) Model {
    m.email = w.email
    return m
}

func (w *EmailPage) Enter(m Model) {
    w.email = m.email
    w.form = newEmailForm(w)
    w.form.Init()
}

func (w *EmailPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
    form, cmd := w.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		w.form = f
        if w.form.State == huh.StateCompleted {
            m.email = w.email
            return true, m, NewNavigateShipping
        }
        return true, m, cmd
	}

    return false, m, nil
}

func (w *EmailPage) Title() string { return "Email" }

func (w *EmailPage) Render(m *Model) string {
    container := lipgloss.NewStyle().
        Height(m.GetMaxPageHeight())

    return container.Render(w.form.View())
}
