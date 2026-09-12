package update

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"logerrorfinder/model"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type LLMResultMsg struct {
	Response *model.LLMResponse
	Err      error
}

func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {

	const (
		logFilePathScreen = iota
		FromDateScreen
		ToDateScreen
		LogListScreen
		LogDetailScreen
	)

	switch msg := msg.(type) {

	case LLMResultMsg:

		m.LLMLoading = false

		if msg.Err != nil {
			m.LLMError = msg.Err
			return m, nil
		}

		m.LLMResponse = msg.Response

		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case "enter":

			if m.Screen == logFilePathScreen {

				m.Screen = FromDateScreen

			} else if m.Screen == FromDateScreen {

				m.Screen = ToDateScreen

			} else if m.Screen == ToDateScreen {

				logs, err := ReadLogs(
					m.LogFilePath,
					m.FromDate,
					m.ToDate,
				)

				if err != nil {
					log.Fatal(err)
				}

				m.ErrorLogs = logs
				m.SelectedErrorLog = 0
				m.Screen = LogListScreen

			} else if m.Screen == LogListScreen {

				selectedError := m.ErrorLogs[m.SelectedErrorLog]

				parts := strings.Split(selectedError, " ")

				if len(parts) < 4 {
					return m, nil
				}

				selectedErrorLog := model.SingleLog{
					Timestamp: parts[0] + " " + parts[1],
					Level:     parts[2],
					Message:   strings.Join(parts[3:], " "),
				}

				logs, err := GetPreviousLogs(
					m.LogFilePath,
					selectedError,
				)

				if err != nil {
					log.Fatal(err)
				}

				var previousLogs []model.SingleLog

				for _, logLine := range logs {

					parts := strings.Split(logLine, " ")

					if len(parts) < 4 {
						continue
					}

					previousLogs = append(previousLogs, model.SingleLog{
						Timestamp: parts[0] + " " + parts[1],
						Level:     parts[2],
						Message:   strings.Join(parts[3:], " "),
					})
				}

				m.LogForLLM = model.LogForLLM{
					SelectedLog:  selectedErrorLog,
					PreviousLogs: previousLogs,
				}

				// Move to detail screen FIRST
				m.Screen = LogDetailScreen

				// Start Gemini analysis
				m.LLMLoading = true
				m.LLMError = nil

				return m, AnalyzeLogsCmd(
					m.GeminiService,
					m.LogForLLM,
				)
			}

		case "esc":
			return m, tea.Quit

		case "backspace":

			if m.Screen == logFilePathScreen && len(m.LogFilePath) > 0 {
				m.LogFilePath = m.LogFilePath[:len(m.LogFilePath)-1]

			} else if m.Screen == FromDateScreen && len(m.FromDate) > 0 {
				m.FromDate = m.FromDate[:len(m.FromDate)-1]

			} else if m.Screen == ToDateScreen && len(m.ToDate) > 0 {
				m.ToDate = m.ToDate[:len(m.ToDate)-1]
			}

		case "up":

			if m.Screen == LogListScreen {
				if m.SelectedErrorLog > 0 {
					m.SelectedErrorLog--
				}
			}

		case "down":

			if m.Screen == LogListScreen {
				if m.SelectedErrorLog < len(m.ErrorLogs)-1 {
					m.SelectedErrorLog++
				}
			}
		case "ctrl+r": // Reset the application
			m = model.NewModel(m.GeminiService)

		default:

			if m.Screen == logFilePathScreen {
				m.LogFilePath += msg.String()

			} else if m.Screen == FromDateScreen {
				m.FromDate += msg.String()

			} else if m.Screen == ToDateScreen {
				m.ToDate += msg.String()
			}
		}
	}

	return m, nil
}

func ReadLogs(logFilePath, fromDate, toDate string) ([]string, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)

		if len(fields) < 4 {
			continue
		}

		logDate := fields[0] + " " + fields[1]
		level := fields[2]

		if logDate < fromDate || logDate > toDate {
			continue
		}

		if level != "ERROR" {
			continue
		}

		logs = append(logs, line)
	}

	return logs, scanner.Err()
}

func GetPreviousLogs(logFilePath string, selectedError string) ([]string, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allLogs []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty or invalid log lines
		fields := strings.Fields(line)

		if len(fields) < 4 {
			continue
		}

		allLogs = append(allLogs, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Find selected error
	selectedIndex := -1

	for i, line := range allLogs {
		if line == selectedError {
			selectedIndex = i
			break
		}
	}

	if selectedIndex == -1 {
		return nil, fmt.Errorf("selected error not found")
	}

	// Get previous 10 logs
	start := selectedIndex - 10

	if start < 0 {
		start = 0
	}

	return allLogs[start:selectedIndex], nil
}

func AnalyzeLogsCmd(gemini model.GeminiAnalyzer, logData model.LogForLLM) tea.Cmd {

	return func() tea.Msg {

		response, err := gemini.AnalyzeLogs(
			context.Background(),
			logData,
		)

		return LLMResultMsg{
			Response: response,
			Err:      err,
		}
	}
}
