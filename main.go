package main

import (
	"log"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused = (m.focused + 1) % 3
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused = (m.focused - 1) % 3
	}
}

type Model struct {
	loaded   bool
	focused  status
	quitting bool
	lists    []list.Model
	err      error
}

func New() *Model {
	return &Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		}

	}

	var cmd tea.Cmd

	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if !m.loaded {
		return "loading..."
	}

	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneView := m.lists[done].View()

	switch m.focused {
	case inProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			focusedStyle.Render(doneView),
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultList.SetShowHelp(false)

	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[todo].Title = "to do"
	m.lists[todo].SetItems([]list.Item{
		Task{title: "title 1 goes here", description: "description 1 goes here", status: todo},
		Task{title: "title 2 goes here", description: "description 2 goes here", status: todo},
		Task{title: "title 3 goes here", description: "description 3 goes here", status: todo},
	})

	m.lists[inProgress].Title = "in progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{title: "title 1 goes here", description: "description 1 goes here", status: inProgress},
		Task{title: "title 2 goes here", description: "description 2 goes here", status: inProgress},
		Task{title: "title 3 goes here", description: "description 3 goes here", status: inProgress},
	})

	m.lists[done].Title = "done"
	m.lists[done].SetItems([]list.Item{
		Task{title: "title 1 goes here", description: "description 1 goes here", status: done},
		Task{title: "title 2 goes here", description: "description 2 goes here", status: done},
		Task{title: "title 3 goes here", description: "description 3 goes here", status: done},
	})
}

type Task struct {
	title       string
	description string
	status      status
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t Task) FilterValue() string {
	return t.title
}

type status int

const (
	todo status = iota
	inProgress
	done
)

const divisor = 4

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)
