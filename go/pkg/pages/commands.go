package pages

import tea "github.com/charmbracelet/bubbletea"

type BeginCheckout struct { }
type StartShipping struct { }
type StartCC struct { }

func NewBeginCheckout() tea.Msg {
    return BeginCheckout{}
}

func NewStartShipping() tea.Msg {
    return StartShipping{}
}

func NewStartCC() tea.Msg {
    return StartCC{}
}


