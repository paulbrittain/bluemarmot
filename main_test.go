package main

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestProgramCanQuit(t *testing.T) {
	m := initialModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestProgramCanGenerateUUIDs(t *testing.T) {
	m := initialModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// some change

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestProgramCanGenerateNames(t *testing.T) {
	m := initialModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(tea.KeyDown)
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
func TestProgramCanGeneratePasswords(t *testing.T) {
	m := initialModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(tea.KeyDown)
	tm.Send(tea.KeyDown)

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
