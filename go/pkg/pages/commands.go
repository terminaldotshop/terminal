package pages

import tea "github.com/charmbracelet/bubbletea"

type BeginCheckout struct { }
type StartShipping struct { }

func NewBeginCheckout() tea.Msg {
    return BeginCheckout{}
}

func NewStartShipping() tea.Msg {
    return StartShipping{}
}


