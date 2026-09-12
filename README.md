<img width="1238" height="725" alt="Screenshot from 2026-09-10 21-04-44" src="https://github.com/user-attachments/assets/7777ff1a-54ff-4488-9cde-304d49e2f5c9" />

# Log Error Finder

**Log Error Finder** is a terminal-based **TUI (Terminal User Interface)** application built with **Go, Bubble Tea, and Google Gemini AI**.

It helps developers inspect application logs, identify error entries, view surrounding log context, and use AI to analyze errors and suggest possible root causes and recommended actions.

Youtube Demo Link: https://youtu.be/XrsZMf9sYG8?si=KYMxjnDDMIjPwpQJ 

## Features

*  Load and analyze log files
*  Filter and identify `ERROR` log entries
*  Browse available errors from the terminal
*  Analyze errors using **Google Gemini AI**
*  Get AI-generated:

  * Root cause
  * Recommended actions
  * Supporting evidence
  * Confidence level
*  Interactive terminal UI using **Bubble Tea**
*  Keyboard-driven navigation
*  Runs directly from the terminal

## Tech Stack

| Technology        | Purpose                 |
| ----------------- | ----------------------- |
| **Go**            | Application development |
| **Bubble Tea**    | Terminal User Interface |
| **Charm**         | TUI ecosystem           |
| **Google Gemini** | AI-powered log analysis |
| **Go Modules**    | Dependency management   |

##  Application Flow

```text
                ┌─────────────────┐
                │   Log File      │
                └────────┬────────┘
                         │
                         ▼
                ┌─────────────────┐
                │  Load Log File  │
                └────────┬────────┘
                         │
                         ▼
                ┌─────────────────┐
                │ Filter ERROR    │
                │     Logs        │
                └────────┬────────┘
                         │
                         ▼
                ┌─────────────────┐
                │   Error List    │
                └────────┬────────┘
                         │
                  Select an Error
                         │
                         ▼
                ┌─────────────────┐
                │   Error Detail  │
                │ + Previous Logs │
                └────────┬────────┘
                         │
                         ▼
                ┌─────────────────┐
                │  Google Gemini  │
                │   AI Analysis   │
                └────────┬────────┘
                         │
                         ▼
          ┌──────────────────────────────┐
          │ Root Cause                   │
          │ Recommended Actions          │
          │ Evidence                     │
          │ Confidence                   │
          └──────────────────────────────┘
```

##  Prerequisites

Before running the project, make sure you have:

* Go installed
* A Google Gemini API key
* A terminal that supports the required TUI features

Check your Go installation:

```bash
go version
```

##  Configure Gemini API

The application requires a `GEMINI_API_KEY` environment variable.

### Linux / macOS

```bash
export GEMINI_API_KEY="your_api_key_here"
```

You can verify it with:

```bash
echo $GEMINI_API_KEY
```

> **Never commit your API key to GitHub.**

For permanent configuration on Linux, you can add the export to your shell configuration file such as `~/.bashrc` or `~/.zshrc`.

##  Run the Application

Clone the repository:

```bash
git clone <YOUR_GITHUB_REPOSITORY_URL>
```

Move into the project directory:

```bash
cd LogErrorFinder
```

Download dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run .
```

##  Build the Application

Create an executable:

```bash
go build -o logErrorFinder
```

Run it:

```bash
./logErrorFinder
```

You can also install it into your Go binary path:

```bash
go install .
```

Then, if your Go bin directory is in your `PATH`, you can run:

```bash
logErrorFinder
```

##  Navigation

The application uses keyboard-based navigation through the Bubble Tea TUI.

Typical controls include:

| Key            | Action    |
| -------------- | --------- |
| `↑`            | Move up   |
| `↓`            | Move down |
| `Enter`        | Select    |
| `Esc`          | Quit      |
| `Ctrl+R`       | Reset     |
| `Backspace`    | Erase     |

> Keyboard controls may vary depending on the current screen and implementation.

##  AI Log Analysis

When an error is selected, the application collects the relevant error information along with previous log entries to provide additional context.

The log context is sent to **Google Gemini**, which analyzes the information and returns structured troubleshooting insights.

The analysis focuses on:

```text
Error
  ↓
Log Context
  ↓
Gemini AI
  ↓
┌─────────────────────┐
│ Root Cause          │
│ Recommended Actions │
│ Evidence            │
│ Confidence          │
└─────────────────────┘
```

The goal is to help developers understand **what went wrong and what they can investigate next**, without manually searching through large log files.

##  Example Use Case

Imagine your application produces thousands of log entries:

```text
INFO  User request received
INFO  Fetching user profile
ERROR Database connection timeout
INFO  Retrying database connection
ERROR Failed to fetch user profile
```

Instead of manually searching through the entire log file, Log Error Finder lets you select the error and inspect its surrounding context.

Gemini can then provide an analysis such as:

```text
Root Cause:
Database connection timeout prevented the application
from retrieving the user profile.

Recommended Actions:
1. Check database availability.
2. Verify connection pool configuration.
3. Check network connectivity.
4. Review database timeout settings.

Confidence:
97%
```

##  Project Structure

```text
LogErrorFinder/
│
├── go.mod
├── go.sum
├── main.go
├── .env
├── api
│   └── gemini.go
├── model
│   └── model.go
├── update
│   └── update.go
└── view
    └── view.go
```

> The exact structure may differ depending on the current implementation.

##  Security

Do not hard-code your Gemini API key in the source code.

Don't do this:

```go
apiKey := "AIza..."
```

Use an environment variable:

```go
apiKey := os.Getenv("GEMINI_API_KEY")
```

Also make sure secrets are excluded from Git:

```gitignore
.env
*.env
```

## Why Bubble Tea?

[Bubble Tea](https://github.com/charmbracelet/bubbletea) provides a powerful framework for building terminal applications using Go.

For this project, it makes it possible to create an interactive workflow:

```text
Log File
   ↓
Error List
   ↓
Select Error
   ↓
View Context
   ↓
AI Analysis
```

without requiring a graphical desktop application.

## Possible Future Improvements

Some ideas for extending the project:

* Support multiple log formats
* Support JSON logs
* Support remote log sources
* Loki integration
* Kubernetes log analysis
* Prometheus metrics context
* Docker log support
* More AI model providers
* OpenAI / Claude integration
* Real-time log monitoring
* Export AI analysis as a report
* Search and filter improvements
* Configurable AI prompts

## Contributing

Contributions, ideas, and improvements are welcome.

1. Fork the repository
2. Create a feature branch

```bash
git checkout -b feature/my-feature
```

3. Make your changes
4. Commit your changes

```bash
git commit -m "Add my feature"
```

5. Push the branch

```bash
git push origin feature/my-feature
```

6. Open a Pull Request

## Support

If you find this project useful:

* Star the repository
* Fork the project
* Report issues
* Suggest new features
* Share it with other Go developers

## Video Demo

This project is also demonstrated in a YouTube video where I walk through the application and show how **Go, Bubble Tea, and Google Gemini AI** can be combined to build an AI-powered developer tool.

If you have ideas for the next **Go, AI, backend, or developer-tool project**, feel free to share them!

---
