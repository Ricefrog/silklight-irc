package frontend

import "github.com/charmbracelet/bubbles/textinput"

type ViewPort struct {
	messages string
}

type InputBox struct {
	textInput textinput.Model
}

type model struct {
	viewPort ViewPort
	inputBox InputBox
}
