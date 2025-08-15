package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	tickMsg time.Time
)

type model struct {
	count int
	input textinput.Model
	label string
}

func New() model {
	ti := textinput.New()
	ti.Placeholder = "Enter label"
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
		case "q":
			return m, tea.Quit
		case "esc":
			m.input.SetValue("")
			return m, nil
		case "enter":
			m.label = m.input.Value()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s: %d\n\n%s\n\nPress Enter to set label • Esc to clear • q to quit\n",
		m.label, m.count, m.input.View(),
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
