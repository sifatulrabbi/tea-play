package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	tickMsg time.Time
)

// ----- Styles -----
var (
	accent     = lipgloss.AdaptiveColor{Light: "#2D5BFF", Dark: "#7AA2F7"}
	subtle     = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	positive   = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#34D399"}
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)
	helpStyle  = lipgloss.NewStyle().Foreground(subtle)
	labelStyle = lipgloss.NewStyle().Foreground(positive)
	panelStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accent).Padding(1, 3).Margin(1, 0).Width(20)
)

type model struct {
	count int
	input textinput.Model
	label string
}

func New() model {
	ti := textinput.New()
	ti.Placeholder = "Type a label and press enter"
	ti.Width = 20
	ti.Focus() // focusing by default.
	return model{
		input: ti,
		label: "Counter",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.count++
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "/exit":
			return m, tea.Quit
		case "esc":
			m.input.SetValue("")
			return m, nil
		case "enter":
			switch val := strings.Trim(m.input.Value(), " \n"); val {
			case "/exit":
				return m, tea.Quit
			default:
				m.label = val
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := titleStyle.Render("⏱  Bubble Tea Counter")
	label := labelStyle.Render(m.label)

	body := fmt.Sprintf("%s: %d", label, m.count)
	panel := panelStyle.Render(body)
	controls := helpStyle.Render("Enter: set label   Esc: clear   q: quit")

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		header,
		panel,
		m.input.View(),
		controls,
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
