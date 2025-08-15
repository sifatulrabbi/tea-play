package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	tickMsg   time.Time
	resultMsg string
	errMsg    string
)

// ----- Styles -----
var (
	accent       = lipgloss.AdaptiveColor{Light: "#2D5BFF", Dark: "#7AA2F7"}
	subtle       = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	positive     = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#34D399"}
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(accent)
	helpStyle    = lipgloss.NewStyle().Foreground(subtle)
	labelStyle   = lipgloss.NewStyle().Foreground(positive)
	panelStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accent).Padding(1, 3).Margin(1, 0).Width(20)
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ef4444"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
)

type model struct {
	count    int
	input    textinput.Model
	response string

	spin   spinner.Model
	busy   bool
	status string
}

func New() model {
	ti := textinput.New()
	ti.Placeholder = "Type a label and press enter"
	ti.Width = 20
	ti.Focus() // focusing by default.

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return model{
		input:    ti,
		spin:     sp,
		response: "Counter",
		status:   "",
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
			if m.busy {
				return m, nil
			}

			val := strings.Trim(m.input.Value(), " \n")
			m.input.Reset()
			switch val {
			case "/exit":
				return m, tea.Quit
			}
			m.busy = true
			m.status = "Invoking LLM model..."
			return m, tea.Batch(m.spin.Tick, doWork(val))
		}

	case spinner.TickMsg:
		if m.busy {
			var cmd tea.Cmd
			m.spin, cmd = m.spin.Update(msg)
			return m, cmd
		}
		return m, nil

	case resultMsg:
		m.busy = false
		m.response = string(msg)
		m.status = successStyle.Render("Done!")
		return m, nil

	case errMsg:
		m.busy = false
		m.response = string(msg)
		m.status = errorStyle.Render("Done with error!")
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := titleStyle.Render("🧵 Async Demo — spinner + command")
	body := fmt.Sprintf("Count: %d\n%s", m.count, m.response)
	panel := panelStyle.Render(body)

	busyLine := ""
	if m.busy {
		busyLine = fmt.Sprintf("%s  Processing...", m.spin.View())
	}

	controls := helpStyle.Render("Enter: run async task   Esc: clear   /exit: quit")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n\n%s\n",
		header,
		m.input.View(),
		panel,
		busyLine,
		m.status,
		controls,
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func doWork(s string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		if strings.TrimSpace(s) == "" {
			return errMsg(fmt.Errorf("input was empty").Error())
		}
		return resultMsg(strings.ToUpper(s))
	}
}

func main() {
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
