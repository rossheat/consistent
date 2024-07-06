package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrMsg error

type Route int

const (
	QuestionRoute = iota
	LoadingRoute
	ResultsRoute
)

var routeName = map[Route]string{
	QuestionRoute: "question",
	LoadingRoute:  "loading",
	ResultsRoute:  "results",
}

func (s Route) String() string {
	return routeName[s]
}

type Model struct {
	TextInput textinput.Model
	Spinner   spinner.Model
	Err       error
	Route     Route
}

func InitialModel() Model {
	// TextInput
	ti := textinput.New()
	ti.Placeholder = "Is the square root of 441 greater than 20?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{Err: nil, TextInput: ti, Route: QuestionRoute, Spinner: s}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, textinput.Blink)
}

type LLMResults struct {
	yesCount int
	noCount  int
}

func (m Model) AskQuestion() tea.Msg {
	time.Sleep(time.Second * 5)
	return LLMResults{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case LLMResults:
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
	case ErrMsg:
		m.Err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	switch m.Route {
	case QuestionRoute:
		return m.QuestionView()
	case LoadingRoute:
		return m.LoadingView()
	case ResultsRoute:
		return m.ResultsView()
	default:
		return "Unknown route"
	}
}

func (m Model) QuestionView() string {
	return fmt.Sprint(
		"Ask the LLM a yes/no question:",
		"\n\n",
		m.TextInput.View(),
		"\n\n",
		"(q to quit)",
	) + "\n"
}

func (m Model) LoadingView() string {
	return fmt.Sprint("\n", m.Spinner.View(), "Asking the the LLM: ", m.TextInput.Value(), "\n\n", "(q to quit)")
}

func (m Model) ResultsView() string {
	return fmt.Sprint("Results: y:20, n:10", "\n\n", "(r to reset, q to quit)")
}

func main() {
	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		panic(fmt.Sprintf("Error returned from Run: %v", err))
	}
}
