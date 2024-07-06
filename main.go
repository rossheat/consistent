package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	placeholder string
}

func InitialModel() Model {
	return Model{placeholder: "Placeholder value"}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "Example view\n"
	s += fmt.Sprintf("\n%v\n", m.placeholder)
	s += "\nPress 'q' to quit.\n"
	return s
}

func main() {
	program := tea.NewProgram(InitialModel())
	if _, err := program.Run(); err != nil {
		panic(fmt.Sprintf("Error returned from Run: %v", err))
	}
}
