package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	uuid "github.com/google/uuid"
	"github.com/goombaio/namegenerator"
)

var choices = []string{"UUID", "Name", "Password"}

var mainStyle = lipgloss.NewStyle().MarginLeft(2)

const (
	UUID int = iota
	Name
	Password
)

const (
	TypeChoiceStage int = iota
	NumOfOutputsStage
	ResultStage
)

type model struct {
	cursor    int
	choice    int
	textInput textinput.Model
	err       error
	quitting  bool
	stage     int
	finished  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "5"
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 5
	ti.SetValue("5")

	return model{
		stage:     0,
		cursor:    0,
		choice:    0,
		textInput: ti,
		err:       nil,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.stage == TypeChoiceStage {
		return updateChoice(msg, m)
	}
	return updateNumOf(msg, m)
}

func updateChoice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = m.cursor
			m.stage = NumOfOutputsStage
			return m, nil

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func updateNumOf(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.stage = ResultStage
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	log.Info("End of updateNumOf")
	return m, cmd
}

func (m model) View() string {
	var s string
	if m.quitting {
		return "\n  See you later!\n\n"
	}

	log.Info(m.stage)

	switch m.stage {
	case TypeChoiceStage:
		s = typeChoicesView(m)
		return mainStyle.Render("\n" + s + "\n\n")
	case NumOfOutputsStage:
		s = numberOfChoicesView(m)
		return mainStyle.Render("\n" + s + "\n\n")
	case ResultStage:
		s = resultView(m)
		return mainStyle.Render("\n" + s + "\n\n")
	default:
		return mainStyle.Render("\n" + s + "\n\n")
	}
}

func typeChoicesView(m model) string {
	s := strings.Builder{}
	s.WriteString("What kind of ID would you like to generate??\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func numberOfChoicesView(m model) string {
	return fmt.Sprintf(
		"How many of these do you want??\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func resultView(m model) string {
	numberOfOutputs := m.textInput.Value()
	log.Info(numberOfOutputs)
	num, err := strconv.Atoi(numberOfOutputs)
	if err != nil {
		log.Error("Uh oh, issue with atoi")
		return ""
	}

	var results []string

	switch m.choice {
	case UUID:
		results = generateUUIDs(num)
	case Name:
		results = generateNames(num)
	case Password:
		results = generateSecureToken(num)
	}

	var retStr string
	for _, v := range results {
		retStr = retStr + v + "\n"
	}
	return fmt.Sprintf(retStr)
}

func generateUUIDs(n int) []string {
	var res []string

	for range n {
		resId := uuid.New()
		res = append(res, resId.String())
	}
	return res
}

func generateNames(n int) []string {
	var res []string

	for range n {
		seed := time.Now().UTC().UnixNano()
		nameGenerator := namegenerator.NewNameGenerator(seed)
		resName := nameGenerator.Generate()
		res = append(res, resName)
	}

	return res
}

func generateSecureToken(n int) []string {
	var res []string
	for range n {
		b := make([]byte, 15)
		if _, err := rand.Read(b); err != nil {
			log.Error("Error generating secure tokens")
			return res
		}
		res = append(res, hex.EncodeToString(b))
	}

	return res
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Info("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
