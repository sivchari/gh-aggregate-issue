package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TitleModel struct {
	model textinput.Model
	err   error
}

var _ tea.Model = (*TitleModel)(nil)

func NewTitle() *TitleModel {
	model := textinput.New()
	model.Placeholder = "Aggregate Issue"
	model.Focus()
	model.Width = 20
	return &TitleModel{
		model: model,
	}
}

func (m *TitleModel) Value() string {
	return m.model.Value()
}

func (m *TitleModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TitleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m *TitleModel) View() string {
	return fmt.Sprintf(
		"Enter a aggregate issue title\n\n%s\n\n%s",
		m.model.View(),
		"(Press ESC to quit)",
	)
}
