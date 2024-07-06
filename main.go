package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const MODEL_NAME = "Claude Sonnet 3.5"
const MODEL_INSTANCES = 50

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
	TextInput       textinput.Model
	Spinner         spinner.Model
	Err             error
	Route           Route
	LLMResults      LLMResults
	AnthropicAPIKey string
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

type AnthropicAPIKey struct {
	value string
}

func (m *Model) LoadEnvVars() tea.Msg {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("missing ANTHROPIC_API_KEY environment variable.\nPlease set it in your shell like this: ANTHROPIC_API_KEY=<value>")
	}
	return AnthropicAPIKey{apiKey}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.LoadEnvVars, m.Spinner.Tick, textinput.Blink)
}

type LLMResults struct {
	yes int
	no  int
}

type AnthropicReqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicReqBody struct {
	Model     string                `json:"model"`
	MaxTokens int                   `json:"max_tokens"`
	Messages  []AnthropicReqMessage `json:"messages"`
}

type AnthropicRespBody struct {
	Content []AnthropicRespBodyContent `json:"content"`
}

type AnthropicRespBodyContent struct {
	Text string `json:"text"`
}

type AnthropicRespContentText struct {
	Answer string `json:"answer"`
}

func SendMessage(m Model) (string, error) {

	client := &http.Client{}

	content := fmt.Sprintf("You MUST produce a correctly formatted JSON response to the following yes/no question '%v'. You can ONLY answer 'yes' or 'no'. If the 'question' is not a valid question or makes no sense, your response will be no. You MUST respond in the following JSON format: {'answer': <'yes'/'no'>}", m.TextInput.Value())

	reqBody := AnthropicReqBody{
		Model:     "claude-3-5-sonnet-20240620",
		MaxTokens: 1024,
		Messages: []AnthropicReqMessage{
			{Role: "user", Content: content},
		},
	}

	bs, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(bs))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("x-api-key", m.AnthropicAPIKey)
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response status code: %v", resp.StatusCode)
	}

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	var respBody AnthropicRespBody
	json.Unmarshal(body, &respBody)

	var respJSON AnthropicRespContentText
	json.Unmarshal([]byte(respBody.Content[0].Text), &respJSON)

	return respJSON.Answer, nil
}

func (m Model) AskQuestion() tea.Msg {
	r := LLMResults{yes: 0, no: 0}
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errChan := make(chan error, MODEL_INSTANCES)

	for i := 0; i < MODEL_INSTANCES; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 500)
		go func() {
			defer wg.Done()
			answer, err := SendMessage(m)
			if err != nil {
				errChan <- err
				return
			}

			mutex.Lock()
			defer mutex.Unlock()
			if answer == "yes" {
				r.yes++
			} else if answer == "no" {
				r.no++
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return ErrMsg(fmt.Errorf("error from the Anthropic API: %v", err))
		}
	}

	return r
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

	case AnthropicAPIKey:
		m.AnthropicAPIKey = msg.value
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
		"Ask ", MODEL_NAME, " a yes/no question:",
		"\n\n",
		m.TextInput.View(),
		"\n\n",
		"(q to quit)",
	) + "\n"
}

func (m Model) LoadingView() string {
	return fmt.Sprint("\n", m.Spinner.View(), "Asking ", MODEL_INSTANCES, " instances of ", MODEL_NAME, ":\n\"", m.TextInput.Value(), "\"\n\n", "(q to quit)")
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

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		panic(fmt.Sprintf("Error returned from Run: %v", err))
	}
}
