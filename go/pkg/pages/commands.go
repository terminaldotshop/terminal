package pages

import tea "github.com/charmbracelet/bubbletea"

type NavigateProduct struct { }
type NavigateEmail struct { }
type NavigateShipping struct { }
type NavigateCC struct { }

func NewNavigateEmail() tea.Msg {
    return NavigateEmail{}
}

func NewNavigateShipping() tea.Msg {
    return NavigateShipping{}
}

func NewNavigateCC() tea.Msg {
    return NavigateCC{}
}

func NewNavigateProduct() tea.Msg {
    return NavigateProduct{}
}


