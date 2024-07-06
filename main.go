package main

import (
	"fmt"
	"os"
	"time"

	"github.com/NimbleMarkets/ntcharts/barchart"
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
	TextInput  textinput.Model
	Spinner    spinner.Model
	Err        error
	Route      Route
	LLMResults LLMResults
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

func (m Model) CheckEnvVars() tea.Msg {

	requiredEnvVars := []string{"AI_API_KEY", "MODEL_NAME"}
	missingEnvVars := make([]string, 0)

	for _, requiredEnvVar := range requiredEnvVars {
		if os.Getenv(requiredEnvVar) == "" {
			missingEnvVars = append(missingEnvVars, requiredEnvVar)
		}
	}

	if len(missingEnvVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v\n\nPlease set them like this <ENV_VAR>=<VALUE> and try again.", missingEnvVars)
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.CheckEnvVars, m.Spinner.Tick, textinput.Blink)
}

type LLMResults struct {
	yes int
	no  int
}

func (m Model) AskQuestion() tea.Msg {
	time.Sleep(time.Second * 1)
	return LLMResults{yes: 68, no: 42}
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

func (m Model) View() string {

	if m.Err != nil {
		return m.ErrorView()
	}

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

	d1 := barchart.BarData{
		Label: fmt.Sprintf("Yes (%v)", m.LLMResults.yes),
		Values: []barchart.BarValue{
			{Name: "Item1", Value: float64(m.LLMResults.yes), Style: lipgloss.NewStyle().Foreground(lipgloss.Color("10"))}},
	}
	d2 := barchart.BarData{
		Label: fmt.Sprintf(" No (%v)", m.LLMResults.no),
		Values: []barchart.BarValue{
			{Name: "Item1", Value: float64(m.LLMResults.no), Style: lipgloss.NewStyle().Foreground(lipgloss.Color("9"))}},
	}

	bc := barchart.New(18, 10)
	bc.PushAll([]barchart.BarData{d1, d2})
	bc.Draw()

	return fmt.Sprint("\n", m.TextInput.Value(), "\n\n", bc.View(), "\n\n", "(r to reset, q to quit)")
}

func (m Model) ErrorView() string {
	return fmt.Sprint("Error: ", m.Err, "\n\n", "(q to quit)")
}

func main() {
	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		panic(fmt.Sprintf("Error returned from Run: %v", err))
	}
}
