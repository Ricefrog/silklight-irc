package frontend

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var initialModel = mainModel{}

func fillerMessages() string {
	var b strings.Builder

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			fmt.Fprintf(&b, "%d", i%10)
		}
		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}

func ClearScreen() {
	for i := 0; i < 200; i++ {
		fmt.Println()
	}
}

func Start() {
	initialModel.selectMode = true
	initialModel.messages = fillerMessages()
	initialModel.initViewport(10, 10)
	initialModel.initTextBox()

	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
