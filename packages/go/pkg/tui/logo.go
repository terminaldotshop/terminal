package tui

func (m model) LogoView() string {
	return m.theme.TextAccent().Bold(true).Render("terminal") + m.CursorView()
}
