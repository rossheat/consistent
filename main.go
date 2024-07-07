package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var config Config

func main() {

	config = NewConfig()

	if config.debug {
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
