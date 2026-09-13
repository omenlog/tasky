package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	formStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(40)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63"))
)

type taskForm struct {
	input  textinput.Model
	err    error
	width  int
	height int
}

type taskFormSubmittedMsg struct {
	title string
}

type taskFormCancelledMsg struct{}

func newTaskForm() taskForm {
	input := textinput.New()
	input.Placeholder = "Buy milk"
	input.Prompt = "Task title: "
	input.CharLimit = 120
	input.Width = 28

	return taskForm{input: input}
}

func (f *taskForm) SetSize(width, height int) {
	f.width = width
	f.height = height
}

func (f *taskForm) Open() tea.Cmd {
	f.err = nil
	f.input.SetValue("")
	return f.input.Focus()
}

func (f taskForm) Update(msg tea.Msg) (taskForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			f.input.Blur()
			return f, func() tea.Msg { return taskFormCancelledMsg{} }
		case "enter":
			title := strings.TrimSpace(f.input.Value())
			if title == "" {
				f.err = fmt.Errorf("a task needs a title")
				return f, nil
			}

			f.input.Blur()
			return f, func() tea.Msg { return taskFormSubmittedMsg{title: title} }
		}
	}

	var cmd tea.Cmd
	f.input, cmd = f.input.Update(msg)
	return f, cmd
}

func (f taskForm) View() string {
	content := titleStyle.Render("New task") + "\n\n" + f.input.View()
	content += "\n\n" + helpStyle.Render("Enter: save · Esc: cancel")
	if f.err != nil {
		content += "\n\n" + errorStyle.Render("Error: "+f.err.Error())
	}

	form := formStyle.Render(content)
	if f.width == 0 || f.height == 0 {
		return form
	}

	return lipgloss.Place(f.width, f.height, lipgloss.Center, lipgloss.Center, form)
}
