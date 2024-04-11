package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	figure "github.com/common-nighthawk/go-figure"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()
	sshPort := os.Getenv("SSH_PORT")
	httpPort := os.Getenv("HTTP_PORT")
	if sshPort == "" {
		sshPort = "2222"
	}
	if httpPort == "" {
		httpPort = "8000"
	}
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("0.0.0.0", sshPort)),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	log.Info("Starting SSH server", "port", sshPort)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			cancel()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.terminal.shop", http.StatusFound)
	})

	// Listen on port 80
	go func() {
		defer cancel()
		err := http.ListenAndServe(":"+httpPort, nil)
		if err != nil {
			log.Fatal("ListenAndServe error:", err)
		}
	}()

	<-ctx.Done()
	s.Shutdown(ctx)
	slog.Info("Shutting down server")
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	renderer := bubbletea.MakeRenderer(s)

	m := model{
		timer:    timer.NewWithInterval(time.Hour*24, time.Millisecond),
		term:     pty.Term,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		renderer: renderer,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	timer     timer.Model
	term      string
	width     int
	height    int
	fontIndex int
	renderer  *lipgloss.Renderer
}

type tickMessage struct{}

var tickEverySecond = tea.Tick(time.Millisecond, func(time.Time) tea.Msg {
	return tickMessage{}
})

func (m model) Init() tea.Cmd {
	return tickEverySecond
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMessage:
		return m, tickEverySecond
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.fontIndex++
			if len(fonts) <= m.fontIndex {
				m.fontIndex = 0
			}
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

var date = time.Date(2024, time.April, 18, 0, 0, 0, 0, time.UTC)

func (m model) View() string {
	timeDiff := date.Sub(time.Now())
	days := int(timeDiff.Hours() / 24)
	hours := int(timeDiff.Hours()) % 24
	minutes := int(timeDiff.Minutes()) % 60
	seconds := int(timeDiff.Seconds()) % 60
	// milliseconds := int(timeDiff.Milliseconds()) % 1000
	font := fonts[m.fontIndex]

	result := strings.Join(
		[]string{
			figure.NewFigure(
				fmt.Sprintf("%02dd %02dh", days, hours),
				font, true).String(),
			figure.NewFigure(
				fmt.Sprintf("%02dm %02ds", minutes, seconds),
				font, true).String(),
			// figure.NewFigure(
			// 	fmt.Sprintf("%03dms", milliseconds),
			// 	font, true).String(),
		},
		"\n",
	)

	lines := strings.Split(result, "\n")
	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	return m.renderer.NewStyle().PaddingLeft((m.width - width) / 2).PaddingTop((m.height - height) / 2).Render(result)
}

var fonts = []string{
	"eftipiti",
	"dotmatrix",
	"3-d",
	"3x5",
	"5lineoblique",
	"acrobatic",
	"alligator",
	"alligator2",
	"alphabet",
	"avatar",
	"banner",
	"banner3-D",
	"banner3",
	"banner4",
	"barbwire",
	"basic",
	"bell",
	"big",
	"bigchief",
	"binary",
	"block",
	"bubble",
	"bulbhead",
	"calgphy2",
	"caligraphy",
	"catwalk",
	"chunky",
	"coinstak",
	"colossal",
	"computer",
	"contessa",
	"contrast",
	"cosmic",
	"cosmike",
	"cricket",
	"cursive",
	"cyberlarge",
	"cybermedium",
	"cybersmall",
	"diamond",
	"digital",
	"doh",
	"doom",
	"drpepper",
	"eftichess",
	"eftifont",
	"eftirobot",
	"eftitalic",
	"eftiwall",
	"eftiwater",
	"epic",
	"fender",
	"fourtops",
	"fuzzy",
	"goofy",
	"gothic",
	"graffiti",
	"hollywood",
	"invita",
	"isometric1",
	"isometric2",
	"isometric3",
	"isometric4",
	"italic",
	"ivrit",
	"jazmine",
	"jerusalem",
	"katakana",
	"kban",
	"larry3d",
	"lcd",
	"lean",
	"letters",
	"linux",
	"lockergnome",
	"madrid",
	"marquee",
	"maxfour",
	"mike",
	"mini",
	"mirror",
	"mnemonic",
	"morse",
	"moscow",
	"nancyj-fancy",
	"nancyj-underlined",
	"nancyj",
	"nipples",
	"ntgreek",
	"o8",
	"ogre",
	"pawp",
	"peaks",
	"pebbles",
	"pepper",
	"poison",
	"puffy",
	"pyramid",
	"rectangles",
	"relief",
	"relief2",
	"rev",
	"roman",
	"rot13",
	"rounded",
	"rowancap",
	"rozzo",
	"runic",
	"runyc",
	"sblood",
	"script",
	"serifcap",
	"shadow",
	"short",
	"slant",
	"slide",
	"slscript",
	"small",
	"smisome1",
	"smkeyboard",
	"smscript",
	"smshadow",
	"smslant",
	"smtengwar",
	"speed",
	"stampatello",
	"standard",
	"starwars",
	"stellar",
	"stop",
	"straight",
	"tanja",
	"tengwar",
	"term",
	"thick",
	"thin",
	"threepoint",
	"ticks",
	"ticksslant",
	"tinker-toy",
	"tombstone",
	"trek",
	"tsalagi",
	"twopoint",
	"univers",
	"usaflag",
	"wavy",
	"weird",
}
