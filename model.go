package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	store    *TaskStore
	list     list.Model
	form     taskForm
	creating bool
	err      error
}

type tasksLoadedMsg struct {
	tasks []Task
	err   error
}

func (m model) handleLoadTasks(msg tasksLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, nil
	}

	items := make([]list.Item, len(msg.tasks))
	for index, task := range msg.tasks {
		items[index] = taskItem{task: task}
	}
	m.err = nil
	return m, m.list.SetItems(items)
}

func (m model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.list.SetSize(msg.Width, msg.Height)
	m.form.SetSize(msg.Width, msg.Height)
	return m, nil
}

var (
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type taskCreatedMsg struct {
	err error
}

func newModel(store *TaskStore) model {
	tasks := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	tasks.Title = "Tasks"

	return model{
		store: store,
		list:  tasks,
		form:  newTaskForm(),
	}
}

func (m model) Init() tea.Cmd {
	return loadTasksCmd(m.store)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tasksLoadedMsg:
		return m.handleLoadTasks(msg)

	case taskCreatedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		return m, loadTasksCmd(m.store)

	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)

	case taskFormSubmittedMsg:
		m.creating = false
		return m, createTaskCmd(m.store, msg.title)

	case taskFormCancelledMsg:
		m.creating = false
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || (!m.creating && key == "q") {
			return m, tea.Quit
		}

		if key == "n" {
			m.creating = true
			m.err = nil
			return m, m.form.Open()
		}
	}

	// forward message to form if we are creating
	if m.creating {
		var cmd tea.Cmd
		m.form, cmd = m.form.Update(msg)
		return m, cmd
	}

	if !m.creating {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.creating {
		return m.form.View()
	}

	view := vStack(m.list.View(), helpStyle.Render("n: new Task - q: quit"))

	if m.err != nil {
		view = vStack(view, errorStyle.Render("Error: "+m.err.Error()))
	}

	return view
}

func loadTasksCmd(store *TaskStore) tea.Cmd {
	return func() tea.Msg {
		tasks, err := store.ListTasks()
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func createTaskCmd(store *TaskStore, title string) tea.Cmd {
	return func() tea.Msg {
		_, err := store.CreateTask(title)
		return taskCreatedMsg{err: err}
	}
}
