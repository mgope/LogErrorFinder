package view

import (
	"fmt"
	"logerrorfinder/model"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	detailsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Width(40)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func View(m model.Model) string {

	const (
		logFilePathScreen = iota
		FromDateScreen
		ToDateScreen
		LogListScreen
		LogDetailScreen
	)
	switch m.Screen {
	case logFilePathScreen:
		return logFilePathView(m)

	case FromDateScreen:
		return fromDateView(m)

	case ToDateScreen:
		return toDateView(m)

	case LogListScreen:
		return logListView(m)

	case LogDetailScreen:
		return logDetailView(m)
	}

	return ""

}

func logDetailView(m model.Model) string {
	view := ""

	view += titleStyle.Render("Log Detail")
	view += "\n\n"

	// Selected Error
	view += titleStyle.Render("Selected Error")
	view += "\n"

	selectedErrorLog := m.ErrorLogs[m.SelectedErrorLog]
	view += valueStyle.Render(selectedErrorLog)

	view += "\n\n"

	// AI Analysis
	view += titleStyle.Render("AI Analysis")
	view += "\n\n"

	if m.LLMLoading {
		view += "Analyzing logs with Gemini...\n"
		return view
	}

	if m.LLMError != nil {
		view += "Gemini Error:\n"
		view += m.LLMError.Error()
		view += "\n"
		return view
	}

	if m.LLMResponse == nil {
		view += "No AI analysis available.\n"
		return view
	}

	// Root Cause
	view += "Root Cause:\n"
	view += valueStyle.Render(m.LLMResponse.RootCause)
	view += "\n\n"

	// Recommended Actions
	view += "Recommended Actions:\n"

	for _, action := range m.LLMResponse.RecommendedActions {
		view += "• " + action + "\n"
	}
	view += "\n\n"

	// Evidence
	view += "Evidence:\n"

	for _, evidence := range m.LLMResponse.Evidence {
		view += "• " + evidence + "\n"
	}

	view += "\n"

	// Confidence
	view += "Confidence:\n"
	view += valueStyle.Render(
		fmt.Sprintf("%.0f%%", m.LLMResponse.Confidence),
	)
	view += "\n\n"

	return view
}

func logListView(m model.Model) string {
	view := ""

	view += titleStyle.Render("Log List")
	view += "\n\n"

	view += valueStyle.Render("Log File: " + m.LogFilePath)
	view += "\n"
	view += valueStyle.Render("From: " + m.FromDate)
	view += "\n"
	view += valueStyle.Render("To: " + m.ToDate)
	view += "\n\n"

	for i, log := range m.ErrorLogs {

		if i == m.SelectedErrorLog {
			view += "> " + log + "\n"
		} else {
			view += "  " + log + "\n"
		}
	}

	return view
}

func fromDateView(m model.Model) string {
	view := ""

	view += labelStyle.Render("From Date: ")
	view += valueStyle.Render(m.FromDate)
	view += "\n\n"

	return view
}

func toDateView(m model.Model) string {
	view := ""

	view += labelStyle.Render("To Date: ")
	view += valueStyle.Render(m.ToDate)
	view += "\n\n"

	return view
}

func logFilePathView(m model.Model) string {
	view := ""
	view += labelStyle.Render("Log File Path:")
	view += valueStyle.Render(m.LogFilePath)
	view += "\n\n"
	return view
}
