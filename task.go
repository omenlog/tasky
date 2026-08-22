package main

import (
	"database/sql"
	"fmt"
)

// Task is one item in the task list.
type Task struct {
	ID    int64
	Title string
	Done  bool
}

// TaskStore contains all database operations for tasks.
type TaskStore struct {
	db *sql.DB
}

func (s *TaskStore) ListTasks() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, title, done FROM tasks ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (s *TaskStore) CreateTask(title string) (Task, error) {
	result, err := s.db.Exec("INSERT INTO tasks (title) VALUES (?)", title)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf("get new task ID: %w", err)
	}

	return Task{ID: id, Title: title}, nil
}

func (s *TaskStore) ToggleTask(id int64) error {
	result, err := s.db.Exec("UPDATE tasks SET done = NOT done WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("toggle task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check toggled task: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("toggle task: task %d does not exist", id)
	}

	return nil
}
