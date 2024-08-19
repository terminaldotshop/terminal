package theme

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
)

// copy returns a copy of a TextInputStyles with all children styles copied.
func copyTextStyles(t huh.TextInputStyles) huh.TextInputStyles {
	return huh.TextInputStyles{
		Cursor:      t.Cursor.Copy(),
		Placeholder: t.Placeholder.Copy(),
		Prompt:      t.Prompt.Copy(),
		Text:        t.Text.Copy(),
	}
}

// copy returns a copy of a FieldStyles with all children styles copied.
func copyFieldStyles(f huh.FieldStyles) huh.FieldStyles {
	return huh.FieldStyles{
		Base:           f.Base.Copy(),
		Title:          f.Title.Copy(),
		Description:    f.Description.Copy(),
		ErrorIndicator: f.ErrorIndicator.Copy(),
		ErrorMessage:   f.ErrorMessage.Copy(),
		SelectSelector: f.SelectSelector.Copy(),
		// NextIndicator:       f.NextIndicator.Copy(),
		// PrevIndicator:       f.PrevIndicator.Copy(),
		Option: f.Option.Copy(),
		// Directory:           f.Directory.Copy(),
		// File:                f.File.Copy(),
		MultiSelectSelector: f.MultiSelectSelector.Copy(),
		SelectedOption:      f.SelectedOption.Copy(),
		SelectedPrefix:      f.SelectedPrefix.Copy(),
		UnselectedOption:    f.UnselectedOption.Copy(),
		UnselectedPrefix:    f.UnselectedPrefix.Copy(),
		FocusedButton:       f.FocusedButton.Copy(),
		BlurredButton:       f.BlurredButton.Copy(),
		TextInput:           copyTextStyles(f.TextInput),
		Card:                f.Card.Copy(),
		NoteTitle:           f.NoteTitle.Copy(),
		Next:                f.Next.Copy(),
	}
}

func copy(t huh.Theme) huh.Theme {
	return huh.Theme{
		Form:           t.Form.Copy(),
		Group:          t.Group.Copy(),
		FieldSeparator: t.FieldSeparator.Copy(),
		Blurred:        copyFieldStyles(t.Blurred),
		Focused:        copyFieldStyles(t.Focused),
		Help: help.Styles{
			Ellipsis:       t.Help.Ellipsis.Copy(),
			ShortKey:       t.Help.ShortKey.Copy(),
			ShortDesc:      t.Help.ShortDesc.Copy(),
			ShortSeparator: t.Help.ShortSeparator.Copy(),
			FullKey:        t.Help.FullKey.Copy(),
			FullDesc:       t.Help.FullDesc.Copy(),
			FullSeparator:  t.Help.FullSeparator.Copy(),
		},
	}
}
