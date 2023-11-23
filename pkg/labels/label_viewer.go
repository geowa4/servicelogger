package labels

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/geowa4/servicelogger/pkg/templates"
	"os"
	"strings"
)

func FindFilesWithTags() map[string][]*templates.Template {
	tagMap := map[string][]*templates.Template{}
	templates.WalkTemplates(func(template *templates.Template) {
		for _, tag := range template.Tags {
			tagMap[tag] = append(tagMap[tag], template)
		}
		if template.Tags == nil {
			tagMap["untagged"] = append(tagMap["untagged"], template)
		}
	})
	return tagMap
}

var (
	perPage          = 30
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

type model struct {
	searchText        string
	templates         []*templates.Template
	templatesCount    int
	templateCursor    int
	templateSelection *templates.Template

	pager paginator.Model
}

func initialModel() model {
	allTemplates := make([]*templates.Template, 0)
	templates.WalkTemplates(func(template *templates.Template) {
		allTemplates = append(allTemplates, template)
	})
	numTemplates := len(allTemplates)

	pager := paginator.New()
	pager.SetTotalPages((numTemplates + perPage - 1) / perPage)
	pager.PerPage = perPage
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")
	pager.KeyMap = paginator.KeyMap{}
	pager.Page = 0
	return model{
		searchText:        "",
		templates:         allTemplates,
		templatesCount:    numTemplates,
		templateCursor:    0,
		templateSelection: nil,

		pager: pager,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		start, end := m.pager.GetSliceBounds(m.templatesCount)
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.templateCursor--
			if m.templateCursor < 0 {
				m.templateCursor = len(m.templates) - 1
				m.pager.Page = m.pager.TotalPages - 1
			}
			if m.templateCursor < start {
				m.pager.PrevPage()
			}
		case "down":
			m.templateCursor++
			if m.templateCursor >= len(m.templates) {
				m.templateCursor = 0
				m.pager.Page = 0
			}
			if m.templateCursor >= end {
				m.pager.NextPage()
			}
		case "enter":
			m.templateSelection = m.templates[m.templateCursor]
			//return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	start, end := m.pager.GetSliceBounds(len(m.templates))
	height := end - start
	for i, template := range m.templates[start:end] {
		cursor := " "
		if i == m.templateCursor%height {
			cursor = ">"
		}
		checked := " "
		if m.templateSelection != nil && template.Summary == m.templateSelection.Summary {
			checked = "x"
		}
		s.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, template.Summary))
	}
	if m.pager.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", perPage-m.pager.ItemsOnPage(len(m.templates))+1))
		s.WriteString(" " + m.pager.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, subduedStyle.Render("↑ / ↓ to navigate; <enter> to select; q to quit\n"), s.String(), "\nSearch Text")
}

func SearchProgram() {
	tm, err := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error running search program: %v\n", err)
		os.Exit(1)
	}
	m, ok := tm.(model)
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "received unexpected model type from program: %v\n", err)
		os.Exit(1)
	}
	templateJson, err := json.Marshal(m.templateSelection)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error printing selected template: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(templateJson))
}
