package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"

	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/terminalhq/terminal/go/pkg/pages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

//go:embed react-miami.txt
var reactMiami string

type PasswordState int

const (
	PasswordSkip PasswordState = iota
	PasswordPossible
	PasswordWaiting
	PasswordAccepted
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
		wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
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

func newPasswordForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("password").
				Password(true).
				Description("Secret Password for Beta Testers"),
		),
	)

	return form

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

	width := pty.Window.Width
	height := pty.Window.Height

	publicKey := s.PublicKey()
	hash := md5.Sum(publicKey.Marshal())
	stringKey := hex.EncodeToString(hash[:])

	model, err := pages.NewModel(renderer, width, height, stringKey)
	if err != nil {
		log.Fatal(err)
	}

	// usePages := os.Getenv("REACT_MIAMI") != "" || true
	coffeePassword := os.Getenv("COFFEE_PASSWORD")

	var passwordForm *huh.Form
	passwordState := PasswordSkip
	if coffeePassword != "" {
		passwordState = PasswordPossible
		passwordForm = newPasswordForm()
	}

	return sshModel{
		failed:   err != nil,
		usePages: false,
		model:    model,

		renderer: renderer,
		width:    width,
		height:   height,

		// Simple password stuff
		passwordTruth: coffeePassword,
		passwordState: passwordState,
		passwordForm:  passwordForm,
	}, []tea.ProgramOption{tea.WithAltScreen()}
}

type sshModel struct {
	failed   bool
	usePages bool
	loggedIn bool
	model    pages.Model

	renderer *lipgloss.Renderer
	width    int
	height   int

	// Simple password_truth stuff
	passwordState PasswordState
	passwordTruth string
	passwordForm  *huh.Form
}

func (m sshModel) Init() tea.Cmd {
	if m.passwordForm != nil {
		return tea.Batch(m.model.Init(), m.passwordForm.Init())
	}

	return m.model.Init()
}

type setPasswordValue struct{ value string }

func (m sshModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.passwordState == PasswordWaiting {
		form, cmd := m.passwordForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.passwordForm = f
		}

		if m.passwordForm.State == huh.StateCompleted {
			passwordValue := m.passwordForm.GetString("password")
			if passwordValue == m.passwordTruth {
				m.passwordState = PasswordAccepted
				m.usePages = true
			} else {
				m.passwordState = PasswordPossible
				m.passwordForm = newPasswordForm()
				m.passwordForm.Init()
			}
		}

		return m, cmd
	}

	if m.usePages {
		return m.model.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+p":
			if m.passwordState == PasswordPossible {
				m.passwordState = PasswordWaiting
			}
		}
	}

	return m, nil
}

func (m sshModel) View() string {
	if m.passwordState == PasswordWaiting {
		return m.passwordForm.View()
	}

	if m.failed {
		return "Sorry, this is definitely @ThePrimeagen's fault. Come back later"
	}

	if m.usePages {
		return m.model.View()
	}

	minWidth := 170
	minheight := 40
	text := "GET YOUR COFFEE AT\n" + reactMiami + "\nCome Next Week to Order 'Online'. Come to the booth on Floor 10 to order in person!"
	if m.width < minWidth || m.height < minheight {
		text = `Get Your Coffee at React Miami!
Come Next Week to Order 'Online'.
Come to the booth on Floor 10 to order in person!
(Zoom out to see the whole message)`
	}

	return m.renderer.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render(text)
}
