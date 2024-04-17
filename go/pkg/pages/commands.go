package pages

import tea "github.com/charmbracelet/bubbletea"

type NavigateProduct struct { }
type NavigateEmail struct { }
type NavigateShipping struct { }
type NavigateCCAddress struct { }
type NavigateCC struct { }
type NavigateConfirm struct { }
type Dialog struct { msg string }

func NewDialog(msg string) tea.Cmd {
    return func() tea.Msg {
        return Dialog{msg: msg}
    }
}

func NewNavigateConfirm() tea.Msg {
    return NavigateConfirm{}
}

func NewNavigateEmail() tea.Msg {
    return NavigateEmail{}
}

func NewNavigateShipping() tea.Msg {
    return NavigateShipping{}
}

func NewNavigateCCAddress() tea.Msg {
    return NavigateCCAddress{}
}

func NewNavigateCC() tea.Msg {
    return NavigateCC{}
}

func NewNavigateProduct() tea.Msg {
    return NavigateProduct{}
}


