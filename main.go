package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
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
	Err       error
	Route     Route
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Is the square root of 441 greater than 20?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return Model{Err: nil, TextInput: ti, Route: QuestionRoute}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type LLMResults struct {
	yesCount int
	noCount  int
}

func AskQuestion() tea.Msg {
	time.Sleep(time.Second * 5)
	return LLMResults{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch m.Route {
	case QuestionRoute:
		return m.QuestionUpdate(msg)
	case LoadingRoute:
		return m.LoadingUpdate(msg)
	case ResultsRoute:
		return m.ResultsUpdate(msg)
	}
	return m, nil
}

func (m Model) QuestionUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.Route = LoadingRoute
			return m, AskQuestion
		}
	case ErrMsg:
		m.Err = msg
		return m, nil
	}
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) LoadingUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case LLMResults:
		m.Route = ResultsRoute
		return m, nil

	case ErrMsg:
		m.Err = msg
		return m, nil
	}
	return m, nil
}

func (m Model) ResultsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.Route = QuestionRoute
			m.TextInput.Reset()
			return m, nil
		}
	case ErrMsg:
		m.Err = msg
		return m, nil
	}

	return m, nil
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
	return fmt.Sprintf(
		"Ask the LLM a yes/no question:\n\n%v\n\n%v",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m Model) LoadingView() string {
	return "Loading..."
}

func (m Model) ResultsView() string {
	return fmt.Sprint("Results: y:20, n:10", "\n\n", "(Press Enter to reset)")
}

func main() {
	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		panic(fmt.Sprintf("Error returned from Run: %v", err))
	}
}