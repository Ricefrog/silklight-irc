package frontend

import (
	"fmt"
	"net"
	"strings"

	"silklight/frontend/dynamicViewport"
	futils "silklight/frontend/utils"
	"silklight/irc"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
	//viewPort       viewport.Model
	ViewPort       dynamicViewport.Model
	textBox        textinput.Model
	messages       string
	CurrentChannel string
	state          int // value specifies what is currently highlighted
	selectMode     bool
	width          int
	height         int
	Conn           net.Conn
	ClientName     string
}

var borderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

func (m MainModel) appendMsgCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return futils.AppendMsg(message)
	}
}

func sendMessageCmd(conn net.Conn, channel, message, clientName string) tea.Cmd {
	return func() tea.Msg {
		if strings.HasPrefix(message, "#raw ") {
			message = message[5:]
			irc.SendString(conn, message)
			return futils.AppendMsg(fmt.Sprintf("<%s> Sent raw msg: %s", clientName, message))
		}
		irc.SendMessage(conn, channel, message)
		return futils.AppendMsg(fmt.Sprintf("<%s> %s", clientName, message))
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) initTextBox() {
	ti := textinput.New()
	ti.CharLimit = 450
	ti.Width = 97
	ti.Prompt = fmt.Sprintf("%s >> ", m.ClientName)

	m.textBox = ti
}

func (m *MainModel) initViewport(width, height int) {
	//m.viewPort = viewport.New(width, height)
	m.ViewPort = dynamicViewport.New(width, height)
	m.ViewPort.SetContent(m.messages)
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
				futils.Quit = true
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
			case 1: // textbox
				switch msg.String() {
				case "ctrl+c":
					// "Are you sure you want to quit?" message.
					// Send the server a proper QUIT message.
					return m, tea.Quit
				case "enter":
					message := m.textBox.Value()
					m.textBox.SetValue("")
					return m, sendMessageCmd(m.Conn, m.CurrentChannel, message, m.ClientName)
				case "esc":
					m.selectMode = true
					m.textBox.Blur()
				}
			case 0: // viewport
				switch msg.String() {
				case "ctrl+c":
					// "Are you sure you want to quit?" message.
					// Send the server a proper QUIT message.
					return m, tea.Quit
				case "esc":
					m.selectMode = true
				case "g":
					if m.ViewPort.PrevKey == "g" {
						m.ViewPort.GotoTop()
						m.ViewPort.PrevKey = ""
					} else {
						m.ViewPort.PrevKey = "g"
					}
				case "G":
					m.ViewPort.GotoBottom()
				}
			}
		}
	case futils.AppendMsg:
		m.ViewPort, cmd = m.ViewPort.Update(msg)
		return m, cmd
	case futils.AppendMsgToTextBox:
		m.textBox.SetValue(string(msg) + fmt.Sprintf(" Width: %d\n", m.ViewPort.Width))
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.ViewPort.Width = msg.Width - m.ViewPort.Style.GetHorizontalFrameSize() - 3
		m.textBox.Width = msg.Width
		m.ViewPort.Height = msg.Height - m.ViewPort.Style.GetVerticalFrameSize() - 8
		m.ViewPort, cmd = m.ViewPort.Update(msg)
		return m, cmd
	}

	if !m.selectMode {
		switch m.state {
		case 0:
			m.ViewPort, cmd = m.ViewPort.Update(msg)
			cmds = append(cmds, cmd)
		case 1:
			m.textBox, cmd = m.textBox.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func ViewWithBuilder(s1, s2 string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n%s", s1, s2)
	return b.String()
}

func (m MainModel) View() string {
	vpStyle := borderStyle.Copy().Width(m.ViewPort.Width).Height(m.ViewPort.Height)
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

	vp := vpStyle.Render(m.ViewPort.View())
	tb := tbStyle.Render(m.textBox.View())

	return lipgloss.JoinVertical(lipgloss.Left, vp, tb)
	//return ViewWithBuilder(vp, tb)
}
