package theme

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	renderer *lipgloss.Renderer

	border     lipgloss.TerminalColor
	background lipgloss.TerminalColor
	highlight  lipgloss.TerminalColor
	error      lipgloss.TerminalColor
	body       lipgloss.TerminalColor
	accent     lipgloss.TerminalColor

	base lipgloss.Style
	form *huh.Theme
}

func BasicTheme(renderer *lipgloss.Renderer, highlight *string) Theme {
	base := Theme{
		renderer: renderer,
	}

	base.background = lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FBFCFD"}
	base.border = lipgloss.AdaptiveColor{Dark: "#3A3F42", Light: "#D7DBDF"}
	base.body = lipgloss.AdaptiveColor{Dark: "#889096", Light: "#889096"}
	base.accent = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#11181C"}
	if highlight != nil {
		base.highlight = lipgloss.Color(*highlight)
	} else {
		base.highlight = lipgloss.Color("#FF5C00")
	}
	base.error = lipgloss.Color("203")

	base.base = renderer.NewStyle().Foreground(base.body)
	base.form = HuhTheme(base)

	return base
}

func HuhTheme(theme Theme) *huh.Theme {
	var t huh.Theme

	t.FieldSeparator = theme.renderer.NewStyle().SetString("\n\n")

	f := &t.Focused
	f.Base = theme.renderer.NewStyle().
		PaddingLeft(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeft(true).
		BorderForeground(theme.accent)
	f.Title = theme.renderer.NewStyle().Foreground(theme.body)
	f.Description = theme.renderer.NewStyle().Foreground(theme.body)
	f.TextInput.Cursor = theme.renderer.NewStyle().Foreground(theme.highlight)
	f.TextInput.Placeholder = theme.renderer.NewStyle().Foreground(theme.body)
	f.TextInput.Prompt = theme.renderer.NewStyle().Foreground(theme.accent)
	f.TextInput.Text = theme.renderer.NewStyle().Foreground(theme.accent)
	f.ErrorIndicator = theme.renderer.NewStyle().Foreground(theme.error)
	f.ErrorMessage = theme.renderer.NewStyle().Foreground(theme.error)
	t.Help = help.New().Styles

	t.Blurred = copyFieldStyles(*f)
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Title.Foreground(theme.body)

	// TODO: add other huh form/input styles as needed

	return &t
}

func (b Theme) Body() lipgloss.TerminalColor {
	return b.body
}

func (b Theme) Highlight() lipgloss.TerminalColor {
	return b.highlight
}

func (b Theme) Background() lipgloss.TerminalColor {
	return b.background
}

func (b Theme) Accent() lipgloss.TerminalColor {
	return b.accent
}

func (b Theme) Base() lipgloss.Style {
	return b.base.Copy()
}

func (b Theme) TextBody() lipgloss.Style {
	return b.Base().Foreground(b.body)
}

func (b Theme) TextAccent() lipgloss.Style {
	return b.Base().Foreground(b.accent)
}

func (b Theme) TextHighlight() lipgloss.Style {
	return b.Base().Foreground(b.highlight)
}

func (b Theme) TextError() lipgloss.Style {
	return b.Base().Foreground(b.error)
}

func (b Theme) Form() *huh.Theme {
	return b.form
}

func (b Theme) Border() lipgloss.TerminalColor {
	return b.border
}
