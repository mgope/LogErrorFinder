package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"logerrorfinder/api"
	"logerrorfinder/model"
	"logerrorfinder/update"
	"logerrorfinder/view"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.Background()

	geminiService, err := api.NewGeminiService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	m := model.NewModel(geminiService)

	p := tea.NewProgram(initialModel{
		model: m,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type initialModel struct {
	model model.Model
}

func (m initialModel) Init() tea.Cmd {
	return nil
}

func (m initialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := update.Update(m.model, msg)
	m.model = updatedModel

	return m, cmd
}

func (m initialModel) View() string {
	return view.View(m.model)
}
