package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sivchari/gh-aggregate-issue/gh"
)

type IssueModel struct {
	list   []*issue
	cursor int
	err    error
}

type issue struct {
	number  int
	title   string
	url     string
	checked bool
}

var _ tea.Model = (*IssueModel)(nil)

func NewIssueModel(prs []*gh.Issue) *IssueModel {
	var list []*issue
	for _, i := range prs {
		list = append(list, &issue{
			number: i.Number,
			title:  i.Title,
			url:    i.URL,
		})
	}
	return &IssueModel{
		list: list,
	}
}

func (m *IssueModel) Value() []*gh.Issue {
	var issuess []*gh.Issue
	for _, i := range m.list {
		if i.checked {
			issuess = append(issuess, &gh.Issue{
				Number: i.number,
				Title:  i.title,
				URL:    i.url,
			})
		}
	}
	return issuess
}

func (m *IssueModel) Init() tea.Cmd {
	if len(m.list) == 0 {
		return tea.Quit
	}
	return nil
}

func (m *IssueModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *IssueModel) View() string {
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
