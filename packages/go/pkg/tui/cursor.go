package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type cursorState struct {
	visible bool
}

type CursorTickMsg struct{}

func (m model) CursorInit() tea.Cmd {
	return tea.Every(time.Millisecond*700, func(t time.Time) tea.Msg {
		return CursorTickMsg{}
	})
}

func (m model) CursorUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case CursorTickMsg:
		m.state.cursor.visible = !m.state.cursor.visible
		return m, tea.Every(time.Millisecond*700, func(t time.Time) tea.Msg {
			return CursorTickMsg{}
		})
	}
	return m, nil
}

func (m model) CursorView() string {
	if m.state.cursor.visible {
		return m.theme.Base().Background(m.theme.Highlight()).Render(" ")
	} else {
		return m.theme.Base().Render(" ")
	}
}
