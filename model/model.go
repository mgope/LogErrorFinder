package model

import "context"

type Screen int

type GeminiAnalyzer interface {
	AnalyzeLogs(
		ctx context.Context,
		logData LogForLLM,
	) (*LLMResponse, error)
}

const (
	LogFilePathScreen Screen = iota
	FromDateScreen
	ToDateScreen
	LogListScreen
	LogDetailScreen
)

type Model struct {
	Screen           Screen
	LogFilePath      string
	FromDate         string
	ToDate           string
	ErrorLogs        []string
	SelectedErrorLog int
	PreviousLogs     []string

	LogForLLM LogForLLM

	LLMResponse   *LLMResponse
	LLMLoading    bool
	LLMError      error
	GeminiService GeminiAnalyzer
}

func NewModel(gemini GeminiAnalyzer) Model {
	return Model{
		Screen:        LogFilePathScreen,
		GeminiService: gemini,
	}
}

type SingleLog struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type LogForLLM struct {
	SelectedLog  SingleLog   `json:"selected_log"`
	PreviousLogs []SingleLog `json:"previous_logs"`
}

type LLMResponse struct {
	RootCause          string   `json:"root_cause"`
	RecommendedActions []string `json:"recommended_actions"`
	Evidence           []string `json:"evidence"`
	Confidence         float64  `json:"confidence"`
}
