package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

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
	ti.Placeholder = "Can a set contain itself as an element?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 110

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{Err: nil, TextInput: ti, Route: QuestionRoute, Spinner: s}
}
