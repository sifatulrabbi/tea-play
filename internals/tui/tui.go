package tui

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
	resultMsg string
	errMsg    string
)

// ----- Styles -----

const (
	ROOT_PADDING_X = 2
	ROOT_PADDING_Y = 1
)

var (
	accent   = lipgloss.AdaptiveColor{Light: "#2D5BFF", Dark: "#7AA2F7"}
	subtle   = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	positive = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#34D399"}

	labelSt    = lipgloss.NewStyle().Foreground(positive)
	titleSt    = lipgloss.NewStyle().Bold(true).Foreground(accent).Bold(true)
	helpSt     = lipgloss.NewStyle().Foreground(subtle)
	inputBoxSt = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(accent) //.Margin(1, 1)
	errorSt    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ef4444"))
	successSt  = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
)

type model struct {
	height int
	width  int

	input    textinput.Model
	response string

	spin           spinner.Model
	busy           bool
	status         string
	detailedStatus string
}

func New() model {
	ti := textinput.New()
	ti.Prompt = "❯ "
	ti.Placeholder = "Type a label and press enter"
	ti.Focus() // focusing by default.

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return model{
		input:          ti,
		spin:           sp,
		response:       "Counter",
		status:         "",
		detailedStatus: "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "/exit", "ctrl+c":
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
		m.status = successSt.Render("Done!")
		return m, nil

	case errMsg:
		m.busy = false
		m.response = string(msg)
		m.status = errorSt.Render("Done with error!")
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	maxContentWidth := max(m.width-(ROOT_PADDING_X*2)-2, 1)
	header := titleSt.Render("CLI Agent")
	busyLine := ""
	if m.busy {
		busyLine = fmt.Sprintf("%s  Processing...", m.spin.View())
	}
	m.input.Width = maxContentWidth
	inputField := inputBoxSt.Width(maxContentWidth).Render(m.input.View())
	controls := helpSt.Render("Enter: run async task   Esc: clear   /exit: quit")
	remainingEmptySpaces := m.height -
		ROOT_PADDING_Y*2 -
		lipgloss.Height(header) -
		lipgloss.Height(busyLine) -
		lipgloss.Height(m.status) -
		lipgloss.Height(inputField) -
		lipgloss.Height(controls)
	remainingEmptySpaces = max(remainingEmptySpaces, 1)
	messagesArea := lipgloss.NewStyle().
		Height(remainingEmptySpaces).
		Width(maxContentWidth).
		Background(lipgloss.Color("#333")).
		Render()

	finalView := lipgloss.NewStyle().
		Padding(ROOT_PADDING_Y, ROOT_PADDING_X).
		Width(max(m.width, 1)).
		Height(max(m.height, 1)).
		Render(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
			header,
			messagesArea,
			busyLine,
			m.status,
			inputField,
			controls,
		))
	return finalView
}

// func hrInAccent(width int) string {
// 	line := strings.Repeat("─", width)
// 	return lipgloss.NewStyle().
// 		Foreground(accent).
// 		Render(line)
// }

func doWork(s string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		if strings.TrimSpace(s) == "" {
			return errMsg(fmt.Errorf("input was empty").Error())
		}
		return resultMsg(strings.ToUpper(s))
	}
}

func StartProgram() {
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
