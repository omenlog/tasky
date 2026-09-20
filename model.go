package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type model struct {
	store           *TaskStore
	list            list.Model
	form            taskForm
	creating        bool
	deleting        bool
	taskForDeletion *Task
	err             error
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

type taskDeletedMsg struct {
	err error
}

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

	case taskDeletedMsg:
		m.deleting = false
		if msg.err != nil {
			m.err = msg.err
			m.taskForDeletion = nil
			return m, nil
		}

		m.taskForDeletion = nil
		return m, loadTasksCmd(m.store)

	case taskFormCancelledMsg:
		m.creating = false
		return m, nil

	case tea.KeyMsg:
		key := msg.String()
		listShortcutsEnabled := !m.creating && !m.list.SettingFilter()

		if key == "ctrl+c" || (listShortcutsEnabled && key == "q") {
			return m, tea.Quit
		}

		if m.taskForDeletion != nil {
			switch key {
			case "y":
				if !m.deleting {
					m.deleting = true
					return m, deleteTaskCmd(m.store, m.taskForDeletion.ID)
				}

			case "n", "esc":
				m.taskForDeletion = nil
				return m, nil
			}

			return m, nil
		}

		// delete a task
		if key == "d" && listShortcutsEnabled {
			item, ok := m.list.SelectedItem().(taskItem)
			if !ok {
				fmt.Println("Error getting current task")
				return m, nil
			}

			m.taskForDeletion = &item.task
			m.err = nil

			return m, nil
		}

		// create new task
		if key == "n" && listShortcutsEnabled {
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
	if m.taskForDeletion != nil {
		content := vStack(
			titleStyle.Render("Delete task ?"),
			"",
			fmt.Sprintf("%q?", m.taskForDeletion.Title),
			"",
			helpStyle.Render("y: delete  n/Esc: cancel"),
		)

		return formStyle.Render(content)
	}

	if m.creating {
		return m.form.View()
	}

	view := vStack(m.list.View(), helpStyle.Render("n: new Task - q: quit - d: delete task"))

	if m.err != nil {
		view = vStack(view, errorStyle.Render("Error: "+m.err.Error()))
	}

	return view
}

func deleteTaskCmd(store *TaskStore, taskID int64) tea.Cmd {
	return func() tea.Msg {
		err := store.DeleteTask(taskID)
		return taskDeletedMsg{err: err}
	}
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
