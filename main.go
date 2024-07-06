package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type ErrMsg error

type Model struct {
	TextInput textinput.Model
	Err       error
	Question  string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Is the square root of 441 greater than 20?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return Model{Err: nil, TextInput: ti}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:

		}

	case ErrMsg:
		m.Err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Ask the LLM a yes/no question:\n\n%v\n\n%v",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}

func main() {
	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		panic(fmt.Sprintf("Error returned from Run: %v", err))
	}
}
