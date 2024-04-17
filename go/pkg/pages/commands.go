package pages

import tea "github.com/charmbracelet/bubbletea"

type BeginCheckout struct { }

func NewBeginCheckout() tea.Msg {
    return BeginCheckout{}
}

