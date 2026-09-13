package main

import (
	"github.com/charmbracelet/lipgloss"
)

func vStack(rows ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func hStack(rows ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, rows...)
}
