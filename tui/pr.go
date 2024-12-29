package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sivchari/gh-aggregate-issue/gh"
)

type PRModel struct {
	list   []*pr
	cursor int
	err    error
}

type pr struct {
	number  int
	title   string
	url     string
	checked bool
}

var _ tea.Model = (*PRModel)(nil)

func NewPRModel(prs []*gh.PullRequest) *PRModel {
	var list []*pr
	for _, p := range prs {
		list = append(list, &pr{
			number: p.Number,
			title:  p.Title,
			url:    p.URL,
		})
	}
	return &PRModel{
		list: list,
	}
}

func (m *PRModel) Value() []*gh.PullRequest {
	var prs []*gh.PullRequest
	for _, p := range m.list {
		if p.checked {
			prs = append(prs, &gh.PullRequest{
				Number: p.number,
				Title:  p.title,
				URL:    p.url,
			})
		}
	}
	return prs
}

func (m *PRModel) Init() tea.Cmd {
	if len(m.list) == 0 {
		return tea.Quit
	}
	return nil
}

func (m *PRModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.list[m.cursor].checked = !m.list[m.cursor].checked
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
			if m.cursor == -1 {
				m.cursor = 0
			}
		case tea.KeyDown:
			if m.cursor < len(m.list)-1 {
				m.cursor++
			}
			if m.cursor == len(m.list) {
				m.cursor = len(m.list) - 1
			}
		}
	case error:
		m.err = msg
		return m, nil
	}
	return m, cmd
}

func (m *PRModel) View() string {
	var s string
	for i, p := range m.list {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		checked := " "
		if p.checked {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] #%d %s\n", cursor, checked, p.number, p.title)
	}
	s += "(Press Enter to toggle, ↑↓ to navigate, ESC to proceed)\n"
	return s
}
