package pages

import tea "github.com/charmbracelet/bubbletea"

type EmailPage struct { }

func (w *EmailPage) Update(m Model, _ tea.Msg) (bool, tea.Model, tea.Cmd) {
    return false, m, nil
}

func (w *EmailPage) Title() string { return "Lets Start With The Email" }

func (w *EmailPage) Render(m *Model) string {
    return "git commit, git push"

}
