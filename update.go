package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ErrMsg error

type LLMResults struct {
	yes int
	no  int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case ErrMsg:
		m.Err = msg
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case LLMResults:
		m.LLMResults = msg
		m.Route = ResultsRoute
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			if m.Route == QuestionRoute {
				m.Route = LoadingRoute
				return m, m.AskQuestion
			}
		case "r":
			if m.Route == ResultsRoute {
				m.TextInput.Reset()
				m.Route = QuestionRoute
				return m, nil
			}
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}
