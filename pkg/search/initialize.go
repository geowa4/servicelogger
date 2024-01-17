package search

import (
	"encoding/json"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/geowa4/servicelogger/pkg/templates"
)

type Model struct {
	allTemplates      []*ListableTemplate
	templateSelection *ListableTemplate

	list list.Model

	windowWidth  int
	windowHeight int
}

func NewModel() *Model {
	allTemplates := make([]*ListableTemplate, 0)
	internalTemplate := &ListableTemplate{}
	internalTemplate.SourcePath = ""
	err := json.Unmarshal([]byte(`{
			"severity": "Info",
			"service_name": "SREManualAction",
			"summary": "INTERNAL ONLY, DO NOT SHARE WITH CUSTOMER",
			"description": "${MESSAGE}",
			"internal_only": true
		}`), internalTemplate)
	if err == nil {
		allTemplates = append(allTemplates, internalTemplate)
	}
	templates.WalkTemplates(func(template *templates.Template) {
		allTemplates = append(allTemplates, NewListableTemplate(template))
	})

	items := make([]list.Item, len(allTemplates))
	for i, t := range allTemplates {
		items[i] = t
	}
	d := list.NewDefaultDelegate()
	l := list.New(items, d, 0, 0)
	l.Title = "Service Log Search"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)
	l.InfiniteScrolling = true
	l.KeyMap.Quit.SetKeys("enter", "q")
	l.KeyMap.Quit.SetHelp("enter/q", "select/quit")

	m := &Model{
		allTemplates:      allTemplates,
		templateSelection: allTemplates[0],
		list:              l,
	}
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}
