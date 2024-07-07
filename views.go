package main

import (
	"fmt"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

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
		"\n",
		"Ask ", config.model, " a yes/no question:",
		"\n\n",
		m.TextInput.View(),
		"\n\n",
		"(esc to quit)",
	) + "\n"
}

func (m Model) LoadingView() string {
	return fmt.Sprint("\n", m.Spinner.View(), "Asking ", config.instances, " instances of ", config.model, ":\n\n\"", m.TextInput.Value(), "\"\n\n", "(esc to quit)")
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

	return fmt.Sprint("\n", m.TextInput.Value(), "\n\n", bc.View(), "\n\n", "(r to reset, esc to quit)")
}

func (m Model) ErrorView() string {
	return fmt.Sprint("Error: ", m.Err, "\n\n", "(esc to quit)")
}
