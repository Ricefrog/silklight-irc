package frontend

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
	viewPort       viewport.Model
	textBox        textinput.Model
	messages       string
	currentChannel string
	state          int // value specifies what is currently highlighted
	selectMode     bool
	width          int
	height         int
}

var borderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

// Send to the update function to add more content to the viewport.
type AppendMsg string

func (m MainModel) appendMsgCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return AppendMsg(message)
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) initTextBox() {
	ti := textinput.NewModel()
	ti.CharLimit = 450
	ti.Width = 97

	m.textBox = ti
}

func (m *MainModel) initViewport(width, height int) {
	m.viewPort = viewport.New(width, height)
	m.viewPort.SetContent(m.messages)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// selection mode, then options for when you've selected
	// viewport or text
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.selectMode {
			switch msg.String() {
			case "ctrl+c":
				// "Are you sure you want to quit?" message.
				// Send the server a proper QUIT message.
				return m, tea.Quit
			case "k", "up":
				if m.state > 0 {
					m.state--
				}
			case "j", "down":
				if m.state < 1 {
					m.state++
				}
			case "enter":
				m.selectMode = false
				if m.state == 1 {
					m.textBox.Focus()
				}
			}
		} else {
			switch m.state {
			case 1:
				switch msg.String() {
				case "ctrl+c":
					// "Are you sure you want to quit?" message.
					// Send the server a proper QUIT message.
					return m, tea.Quit
				case "enter":
					//typedMsg := m.textBox.Value()
					//fmt.Fprintf(os.Stderr, "MESSAGE ENTERED: %s\n", typedMsg)
					m.textBox.SetValue("")
				case "esc":
					m.selectMode = true
					m.textBox.Blur()
				}
			case 0:
				switch msg.String() {
				case "ctrl+c":
					// "Are you sure you want to quit?" message.
					// Send the server a proper QUIT message.
					return m, tea.Quit
				case "esc":
					m.selectMode = true
				}
			}
		}
	case AppendMsg:
		var b strings.Builder
		fmt.Fprintf(&b, "%s%s", m.messages, msg)
		m.messages = b.String()
		m.viewPort.SetContent(m.messages)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.viewPort.Width = msg.Width - m.viewPort.Style.GetHorizontalFrameSize() - 3
		//m.viewPort.Width = msg.Width - 3
		m.textBox.Width = msg.Width
		m.viewPort.Height = msg.Height - m.viewPort.Style.GetVerticalFrameSize() - 8
		//m.viewPort.Height = msg.Height - 2 - 4

		m.viewPort.SetContent(m.messages)
		return m, nil
	}

	if !m.selectMode {
		switch m.state {
		case 0:
			m.viewPort, cmd = m.viewPort.Update(msg)
			cmds = append(cmds, cmd)
		case 1:
			m.textBox, cmd = m.textBox.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var b strings.Builder

	vpStyle := borderStyle.Copy().Width(m.width - 3)
	tbStyle := borderStyle.Copy().Width(m.width - 3)

	var hlColor lipgloss.Color
	if m.selectMode {
		hlColor = lipgloss.Color("6")
	} else {
		hlColor = lipgloss.Color("5")
	}

	switch m.state {
	case 0:
		vpStyle = vpStyle.BorderForeground(hlColor)
	case 1:
		tbStyle = tbStyle.BorderForeground(hlColor)
	}

	m.viewPort.Style = vpStyle

	fmt.Fprintf(&b, m.viewPort.View()+"\n")
	fmt.Fprintf(&b, tbStyle.Render(m.textBox.View()))

	return b.String()
}
