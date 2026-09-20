# Task Manager

A small terminal task manager written in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Tasks are stored locally in SQLite. At the moment is a **Work In Progress** project

## Requirements

- Go 1.25 or newer

## Run

```bash
go run .
```

On first run, the application creates its SQLite database at `db/tasks.db`.

## Controls

| Key | Action |
| --- | --- |
| `n` | Create a task |
| `Enter` | Save a new task |
| `Esc` | Cancel task creation or deletion |
| `d` | Delete the selected task |
| `y` | Confirm task deletion |
| `n` | Cancel task deletion |
| `q` | Quit from the task list |
| `Ctrl+C` | Quit |

The task list uses the standard Bubble Tea list key bindings for navigation and filtering.

## Build

```bash
go build -o task-manager .
./task-manager
```
