// Altered charmbracelet's viewport code to work better with an incoming stream of new
// lines to display.
// https://github.com/charmbracelet/bubbles/blob/master/viewport/viewport.go
package dynamicViewport

import (
	"fmt"
	"math"
	futils "silklight/frontend/utils"
	"silklight/styles"
	"silklight/utils"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// New returns a new model with the given width and height as well as default
// keymappings.
func New(width, height int) (m Model) {
	m.Width = width
	m.Height = height
	m.setInitialValues()
	return m
}

// For keeping track of line-breaks
type LineInfo struct {
	literal string
	lines   []string
	length  int
}

func (li *LineInfo) format(width int) {
	formatted := make([]string, 0, len(li.literal)/width+1)
	buffer := fmt.Sprintf("           %s ", styles.Separator)
	line := li.literal

	pastFirst := false
	for len(line) > width-1 {
		formatted = append(formatted, line[:width-1])
		line = line[width-1:]
		if !pastFirst {
			pastFirst = true
			width -= len(buffer)
		}
	}
	if len(line) > 0 {
		formatted = append(formatted, line)
	}

	li.lines = formatted
	for i := 1; i < len(li.lines); i++ {
		li.lines[i] = buffer + li.lines[i]
	}
	li.length = len(li.lines)
}

func length(lines []LineInfo) int {
	num := 0
	for _, li := range lines {
		num += li.length
	}
	return num
}

func subset(lines []LineInfo, start, end int) []string {
	ret := make([]string, 0, end-start)
	startLi := 0
	index := 0

	for i, li := range lines {
		index += li.length
		if index > start {
			startLi = i
			index -= li.length
			break
		}
	}

	for i := startLi; ; i++ {
		for _, line := range lines[i].lines {
			ret = append(ret, line)
			index++
			if index == end {
				goto Loop
			}
		}
	}
Loop:
	return ret
}

// Model is the Bubble Tea model for this viewport element.
type Model struct {
	Width  int
	Height int
	KeyMap KeyMap

	// Whether or not to respond to the mouse. The mouse must be enabled in
	// Bubble Tea for this to work. For details, see the Bubble Tea docs.
	MouseWheelEnabled bool

	// The number of lines the mouse wheel will scroll. By default, this is 3.
	MouseWheelDelta int

	// YOffset is the vertical scroll position.
	YOffset int

	// YPosition is the position of the viewport in relation to the terminal
	// window. It's used in high performance rendering only.
	YPosition int

	// Style applies a lipgloss style to the viewport. Realistically, it's most
	// useful for setting borders, margins and padding.
	Style lipgloss.Style

	// HighPerformanceRendering bypasses the normal Bubble Tea renderer to
	// provide higher performance rendering. Most of the time the normal Bubble
	// Tea rendering methods will suffice, but if you're passing content with
	// a lot of ANSI escape codes you may see improved rendering in certain
	// terminals with this enabled.
	//
	// This should only be used in program occupying the entire terminal,
	// which is usually via the alternate screen buffer.
	HighPerformanceRendering bool

	PrevKey string

	initialized bool
	lines       []LineInfo
}

func (m *Model) NumLines() int {
	return len(m.lines)
}

func (m *Model) NumLinesLiteral() int {
	return length(m.lines)
}

func (m *Model) PrintLiterals() string {
	ret := ""
	for _, lineInfo := range m.lines {
		if lineInfo.literal != "" {
			ret += lineInfo.literal
		} else {
			ret += "Empty literal.\n"
		}
	}
	return ret
}

func (m *Model) PrintLines() string {
	ret := ""
	for _, lineInfo := range m.lines {
		if len(lineInfo.lines) > 0 {
			ret += strings.Join(lineInfo.lines, "\n")
		} else {
			ret += "No lines!\n"
		}
	}
	return ret
}

func (m *Model) setInitialValues() {
	m.KeyMap = DefaultKeyMap()
	m.MouseWheelEnabled = true
	m.MouseWheelDelta = 3
	m.initialized = true
}

// Init exists to satisfy the tea.Model interface for composability purposes.
func (m Model) Init() tea.Cmd {
	return nil
}

// AtTop returns whether or not the viewport is in the very top position.
func (m Model) AtTop() bool {
	return m.YOffset <= 0
}

// AtBottom returns whether or not the viewport is at or past the very bottom
// position.
func (m Model) AtBottom() bool {
	return m.YOffset >= m.maxYOffset()
}

// PastBottom returns whether or not the viewport is scrolled beyond the last
// line. This can happen when adjusting the viewport height.
func (m Model) PastBottom() bool {
	return m.YOffset > m.maxYOffset()
}

// ScrollPercent returns the amount scrolled as a float between 0 and 1.
func (m Model) ScrollPercent() float64 {
	if m.Height >= len(m.lines) {
		return 1.0
	}
	y := float64(m.YOffset)
	h := float64(m.Height)
	t := float64(len(m.lines) - 1)
	v := y / (t - h)
	return math.Max(0.0, math.Min(1.0, v))
}

// SetContent set the pager's text content. For high performance rendering the
// Sync command should also be called.
func (m *Model) SetContent(s string) {
	s = strings.ReplaceAll(s, "\r\n", "\n") // normalize line endings
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		li := LineInfo{
			literal: line,
		}
		li.format(m.Width)
		m.lines = append(m.lines, li)
	}

	if m.YOffset > len(m.lines)-1 {
		m.GotoBottom()
	}
}

func (m *Model) AppendLine(s string) {
	wasAtBottom := m.AtBottom()
	s = strings.ReplaceAll(s, "\r\n", "") // normalize line ending
	s = strings.ReplaceAll(s, "\n", "")   // Do not need newlines
	li := LineInfo{
		literal: s,
	}
	li.format(m.Width)
	m.lines = append(m.lines, li)

	if m.YOffset > len(m.lines)-1 || wasAtBottom {
		m.GotoBottom()
	}
}

// maxYOffset returns the maximum possible value of the y-offset based on the
// viewport's content and set height.
func (m Model) maxYOffset() int {
	return max(0, length(m.lines)-m.Height)
}

// visibleLines returns the lines that should currently be visible in the
// viewport.
/*
func (m Model) origVisibleLines() (lines []string) {
	if len(m.lines) > 0 {
		top := max(0, m.YOffset)
		bottom := clamp(m.YOffset+m.Height, top, len(m.lines))
		lines = m.lines[top:bottom]
	}
	return lines
}
*/

func (m Model) visibleLines() (lines []string) {
	if length(m.lines) > 0 {
		top := max(0, m.YOffset)
		bottom := clamp(m.YOffset+m.Height, top, length(m.lines))
		lines = subset(m.lines, top, bottom)
	}
	return lines
}

// scrollArea returns the scrollable boundaries for high performance rendering.
func (m Model) scrollArea() (top, bottom int) {
	top = max(0, m.YPosition)
	bottom = max(top, top+m.Height)
	if top > 0 && bottom > top {
		bottom--
	}
	return top, bottom
}

// SetYOffset sets the Y offset.
func (m *Model) SetYOffset(n int) {
	m.YOffset = clamp(n, 0, m.maxYOffset())
}

// ViewDown moves the view down by the number of lines in the viewport.
// Basically, "page down".
func (m *Model) ViewDown() []string {
	if m.AtBottom() {
		return nil
	}

	m.SetYOffset(m.YOffset + m.Height)
	return m.visibleLines()
}

// ViewUp moves the view up by one height of the viewport. Basically, "page up".
func (m *Model) ViewUp() []string {
	if m.AtTop() {
		return nil
	}

	m.SetYOffset(m.YOffset - m.Height)
	return m.visibleLines()
}

// HalfViewDown moves the view down by half the height of the viewport.
func (m *Model) HalfViewDown() (lines []string) {
	if m.AtBottom() {
		return nil
	}

	m.SetYOffset(m.YOffset + m.Height/2)
	return m.visibleLines()
}

// HalfViewUp moves the view up by half the height of the viewport.
func (m *Model) HalfViewUp() (lines []string) {
	if m.AtTop() {
		return nil
	}

	m.SetYOffset(m.YOffset - m.Height/2)
	return m.visibleLines()
}

// LineDown moves the view down by the given number of lines.
func (m *Model) LineDown(n int) (lines []string) {
	if m.AtBottom() || n == 0 {
		return nil
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we actually have left before we reach
	// the bottom.
	m.SetYOffset(m.YOffset + n)
	return m.visibleLines()
}

// LineUp moves the view down by the given number of lines. Returns the new
// lines to show.
func (m *Model) LineUp(n int) (lines []string) {
	if m.AtTop() || n == 0 {
		return nil
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we are from the top.
	m.SetYOffset(m.YOffset - n)
	return m.visibleLines()
}

// GotoTop sets the viewport to the top position.
func (m *Model) GotoTop() (lines []string) {
	if m.AtTop() {
		return nil
	}

	m.SetYOffset(0)
	return m.visibleLines()
}

// GotoBottom sets the viewport to the bottom position.
func (m *Model) GotoBottom() (lines []string) {
	m.SetYOffset(m.maxYOffset())
	return m.visibleLines()
}

// Sync tells the renderer where the viewport will be located and requests
// a render of the current state of the viewport. It should be called for the
// first render and after a window resize.
//
// For high performance rendering only.
func Sync(m Model) tea.Cmd {
	if len(m.lines) == 0 {
		return nil
	}
	top, bottom := m.scrollArea()
	return tea.SyncScrollArea(m.visibleLines(), top, bottom)
}

// ViewDown is a high performance command that moves the viewport up by a given
// numer of lines. Use Model.ViewDown to get the lines that should be rendered.
// For example:
//
//     lines := model.ViewDown(1)
//     cmd := ViewDown(m, lines)
//
func ViewDown(m Model, lines []string) tea.Cmd {
	if len(lines) == 0 {
		return nil
	}
	top, bottom := m.scrollArea()
	return tea.ScrollDown(lines, top, bottom)
}

// ViewUp is a high performance command the moves the viewport down by a given
// number of lines height. Use Model.ViewUp to get the lines that should be
// rendered.
func ViewUp(m Model, lines []string) tea.Cmd {
	if len(lines) == 0 {
		return nil
	}
	top, bottom := m.scrollArea()
	return tea.ScrollUp(lines, top, bottom)
}

// Update handles standard message-based viewport updates.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.updateAsModel(msg)
	return m, cmd
}

// Author's note: this method has been broken out to make it easier to
// potentially transition Update to satisfy tea.Model.
func (m Model) updateAsModel(msg tea.Msg) (Model, tea.Cmd) {
	if !m.initialized {
		m.setInitialValues()
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.PageDown):
			lines := m.ViewDown()
			if m.HighPerformanceRendering {
				cmd = ViewDown(m, lines)
			}

		case key.Matches(msg, m.KeyMap.PageUp):
			lines := m.ViewUp()
			if m.HighPerformanceRendering {
				cmd = ViewUp(m, lines)
			}

		case key.Matches(msg, m.KeyMap.HalfPageDown):
			lines := m.HalfViewDown()
			if m.HighPerformanceRendering {
				cmd = ViewDown(m, lines)
			}

		case key.Matches(msg, m.KeyMap.HalfPageUp):
			lines := m.HalfViewUp()
			if m.HighPerformanceRendering {
				cmd = ViewUp(m, lines)
			}

		case key.Matches(msg, m.KeyMap.Down):
			lines := m.LineDown(1)
			if m.HighPerformanceRendering {
				cmd = ViewDown(m, lines)
			}

		case key.Matches(msg, m.KeyMap.Up):
			lines := m.LineUp(1)
			if m.HighPerformanceRendering {
				cmd = ViewUp(m, lines)
			}
		}

	case tea.MouseMsg:
		if !m.MouseWheelEnabled {
			break
		}
		switch msg.Type {
		case tea.MouseWheelUp:
			lines := m.LineUp(m.MouseWheelDelta)
			if m.HighPerformanceRendering {
				cmd = ViewUp(m, lines)
			}

		case tea.MouseWheelDown:
			lines := m.LineDown(m.MouseWheelDelta)
			if m.HighPerformanceRendering {
				cmd = ViewDown(m, lines)
			}
		}
	case futils.AppendMsg:
		m.AppendLine(utils.PrependTimestamp(string(msg)))
	case tea.WindowSizeMsg:
		for i := 0; i < len(m.lines); i++ {
			m.lines[i].format(m.Width)
		}
	}

	return m, cmd
}

// View renders the viewport into a string.
func (m Model) View() string {
	if m.HighPerformanceRendering {
		// Just send newlines since we're going to be rendering the actual
		// content seprately. We still need to send something that equals the
		// height of this view so that the Bubble Tea standard renderer can
		// position anything below this view properly.
		return strings.Repeat("\n", m.Height-1)
	}

	lines := m.visibleLines()

	// Fill empty space with newlines
	extraLines := ""
	if len(lines) < m.Height {
		extraLines = strings.Repeat("\n", max(0, m.Height-len(lines)))
	}

	return m.Style.Copy().
		UnsetWidth().
		UnsetHeight().
		Render(strings.Join(lines, "\n") + extraLines)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
