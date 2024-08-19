package tui

import (
	"embed"
	"encoding/json"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type faqState struct {
	faqs []string
}

//go:embed faq.json
var jsonData embed.FS

type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func LoadFaqs() []FAQ {
	data, err := jsonData.ReadFile("faq.json")
	if err != nil {
		log.Fatalf("Failed to read embedded file: %s", err)
	}
	var faqs []FAQ
	if err := json.Unmarshal(data, &faqs); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	return faqs
}

func (m model) FaqSwitch() (model, tea.Cmd) {
	m = m.SwitchPage(faqPage)
	m.state.footer.commands = []footerCommand{
		{key: "↑↓", value: "scroll"},
		{key: "c", value: "cart"},
	}

	m.state.faq.faqs = []string{}
	for _, faq := range m.faqs {
		m.state.faq.faqs = append(m.state.faq.faqs, m.theme.TextAccent().Width(m.widthContent).Render(faq.Question))
		m.state.faq.faqs = append(m.state.faq.faqs, m.theme.Base().Width(m.widthContent).Render(faq.Answer))
		m.state.faq.faqs = append(m.state.faq.faqs, "")
	}

	return m, nil
}

func (m model) FaqView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.state.faq.faqs...,
	)
}
