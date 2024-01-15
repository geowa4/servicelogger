package search

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/muesli/termenv"
	"os"
)

var (
	verticalPadding   = 1
	horizontalPadding = 2
	paddingStyle      = lipgloss.NewStyle().Padding(verticalPadding, horizontalPadding)
)

type ListableTemplate struct {
	templates.Template
}

func NewListableTemplate(template *templates.Template) *ListableTemplate {
	lt := &ListableTemplate{}
	lt.Severity = template.Severity
	lt.ServiceName = template.ServiceName
	lt.Summary = template.Summary
	lt.Desc = template.Desc
	lt.InternalOnly = template.InternalOnly
	lt.EventStreamId = template.EventStreamId
	lt.Tags = template.Tags
	lt.SourcePath = template.SourcePath
	return lt
}

func (t *ListableTemplate) ToTemplate() *templates.Template {
	return &templates.Template{
		Severity:      t.Severity,
		ServiceName:   t.ServiceName,
		Summary:       t.Summary,
		Desc:          t.Desc,
		InternalOnly:  t.InternalOnly,
		EventStreamId: t.EventStreamId,
		Tags:          t.Tags,
		SourcePath:    t.SourcePath,
	}
}

func (t *ListableTemplate) FilterValue() string {
	return t.String()
}

func (t *ListableTemplate) Title() string {
	return t.Summary
}

func (t *ListableTemplate) Description() string {
	return t.Desc
}

type model struct {
	allTemplates      []*ListableTemplate
	templateSelection *ListableTemplate

	list list.Model

	windowWidth  int
	windowHeight int
}

func initialModel() *model {
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
	l.KeyMap.Quit.SetKeys("enter", "q")
	l.KeyMap.Quit.SetHelp("enter/q", "select/quit")

	m := &model{
		allTemplates:      allTemplates,
		templateSelection: allTemplates[0],
		list:              l,
	}
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	item := newListModel.SelectedItem()
	if template, ok := item.(*ListableTemplate); ok {
		m.templateSelection = template
	}
	return m, cmd
}

func (m *model) getPaneWidth() int {
	x, _ := paddingStyle.GetFrameSize()
	if m.windowWidth <= x {
		return 0
	}
	return (m.windowWidth - x) / 2
}

func (m *model) getPaneHeight() int {
	_, y := paddingStyle.GetFrameSize()
	if m.windowHeight <= y {
		return 0
	}
	return m.windowHeight - y
}

func (m *model) View() string {
	m.list.SetSize(m.getPaneWidth()-horizontalPadding*2, m.getPaneHeight())
	md := m.templateSelection.String()
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("notty"),
		glamour.WithWordWrap(m.getPaneWidth()-1-horizontalPadding*4),
	)
	renderedMd, err := renderer.Render(md)
	if err != nil {
		renderedMd = md
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(m.getPaneWidth()).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("227")).
			BorderLeft(false).BorderTop(false).BorderRight(true).BorderBottom(false).
			Render(
				paddingStyle.Render(m.list.View()),
			),
		paddingStyle.Width(m.getPaneWidth()).Render(renderedMd),
	)
}

func Program() *templates.Template {
	lipgloss.SetColorProfile(termenv.TrueColor)
	tm, err := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error running search program: %v\n", err)
		os.Exit(1)
	}
	m, ok := tm.(*model)
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "received unexpected model type from program: %v\n", err)
		os.Exit(1)
	}
	return m.templateSelection.ToTemplate()
}
