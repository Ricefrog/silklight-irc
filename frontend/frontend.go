package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var initialModel = model{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) initInputBox() {
	ti := textinput.NewModel()
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 120

	m.inputBox.textInput = ti
}

func updateInputBoxView(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	state := &m.inputBox

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// "Are you sure you want to quit?" message.
			// Send the server a proper QUIT message.
			return m, tea.Quit
		case "enter":
			typedMsg := state.textInput.Value()
			fmt.Fprintf(os.Stderr, "MESSAGE ENTERED: %s\n", typedMsg)
			state.textInput.SetValue("")
		}
	}

	var cmd tea.Cmd
	state.textInput, cmd = state.textInput.Update(msg)
	return m, cmd
}

func inputBoxView(m model) string {
	state := m.inputBox
	return state.textInput.View()
}

func ClearScreen() {
	for i := 0; i < 200; i++ {
		fmt.Println()
	}
}

func (m model) View() string {
	return inputBoxView(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return updateInputBoxView(m, msg)
}

func Start() {
	initialModel.initInputBox()
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
