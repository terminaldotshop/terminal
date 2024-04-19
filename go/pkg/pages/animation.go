package pages

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type AnimationPage struct {
	frameIdx int
	frames   []string
}

func NewAnimationPage() *AnimationPage {
	frameFiles, _ := fs.ReadDir(frameDir, "frames")
	frameFilesSorted := []string{}
	for _, f := range frameFiles {
		log.Warn("name:", "name", f.Name())
		if f.IsDir() {
			continue
		}

		frameFilesSorted = append(frameFilesSorted, f.Name())
	}

	sort.Slice(frameFilesSorted, func(i, j int) bool {
		return frameFilesSorted[i] < frameFilesSorted[j]
	})

	frames := []string{}
	for _, name := range frameFilesSorted {
		text, err := frameDir.ReadFile("frames/" + name)
		if err != nil {
			log.Fatal(err)
		}

		frames = append(frames, string(text))
	}

	return &AnimationPage{
		frameIdx: 0,
		frames:   frames,
	}
}

type NextFrameMsg struct{}

func (c *AnimationPage) Enter(m Model)      {}
func (c *AnimationPage) Exit(m Model) Model { return m }

func (s *AnimationPage) Update(m Model, msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg.(type) {
	case NextFrameMsg:
		s.frameIdx = s.frameIdx + 1
		return true, m, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
			return NextFrameMsg{}
		})
	}

	return false, m, nil
}

func (s *AnimationPage) Title() string { return "Animation" }

func (s *AnimationPage) Render(m *Model) string {

	container := lipgloss.NewStyle().
		Height(m.GetMaxPageHeight()).
		Width(m.GetMaxPageWidth()).
		PaddingLeft((m.GetMaxPageWidth() - 20) / 2).
		AlignVertical(lipgloss.Center)

	frameIdx := s.frameIdx % len(s.frames)
	text := s.frames[frameIdx]

	message := "Loading"
	length := len(message)
	messages := []string{}
	for idx := range message {
		messages = append(messages,
			fmt.Sprintf("%s%s",
				message[0:idx],
				strings.Repeat(" ", length-idx),
			),
		)
	}

	messages = append(messages, message)

	for idx := range message {
		messages = append(messages,
			fmt.Sprintf("%s%s",
				message[0:length-idx],
				strings.Repeat(" ", idx),
			),
		)
	}

	message = messages[s.frameIdx%len(messages)]

	lines := strings.Split(text, "\n")
	lines[10] = lines[10][0:4] + message + lines[10][4+len(message):]
	text = strings.Join(lines, "\n")

	return container.Render(
		fmt.Sprintf(
			"%d\n%s", s.frameIdx, text,
		),
	)
}
