package search

import (
	"fmt"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/muesli/termenv"
	"os"
	"strings"
)

var (
	perPage          = 10
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

type model struct {
	searchText        string
	allTemplates      []*templates.Template
	filteredTemplates []*templates.Template
	templateCursor    int
	templateSelection *templates.Template

	pager paginator.Model
}

func (m *model) updateSearchText(newSearchText string) {
	m.searchText = newSearchText

	if m.searchText == "" {
		m.filteredTemplates = m.allTemplates
	} else {
		m.filteredTemplates = make([]*templates.Template, 0)
		for _, template := range m.allTemplates {
			if strings.Contains(template.Summary, m.searchText) ||
				strings.Contains(template.Description, m.searchText) ||
				strings.Contains(strings.Join(template.Tags, ""), m.searchText) {
				// TODO: template should be a Stringer that has all this data (as markdown?)
				m.filteredTemplates = append(m.filteredTemplates, template)
			}
		}
	}

	m.pager.SetTotalPages(len(m.filteredTemplates))
	m.pager.Page = 0
	m.templateCursor = 0
}

func initialModel() *model {
	allTemplates := make([]*templates.Template, 0)
	templates.WalkTemplates(func(template *templates.Template) {
		allTemplates = append(allTemplates, template)
	})

	pager := paginator.New()
	pager.PerPage = perPage
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")
	pager.KeyMap = paginator.KeyMap{}
	m := &model{
		searchText:        "",
		allTemplates:      allTemplates,
		templateSelection: allTemplates[0],

		pager: pager,
	}
	m.updateSearchText("")
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		start, end := m.pager.GetSliceBounds(len(m.filteredTemplates))
		switch keypress := msg.String(); keypress {
		case "up":
			m.templateCursor--
			if m.templateCursor < 0 {
				m.templateCursor = len(m.filteredTemplates) - 1
				m.pager.Page = m.pager.TotalPages - 1
			}
			if m.templateCursor < start {
				m.pager.PrevPage()
			}
			m.templateSelection = m.filteredTemplates[m.templateCursor]
		case "down":
			m.templateCursor++
			if m.templateCursor >= len(m.filteredTemplates) {
				m.templateCursor = 0
				m.pager.Page = 0
			}
			if m.templateCursor >= end {
				m.pager.NextPage()
			}
			m.templateSelection = m.filteredTemplates[m.templateCursor]
		case "left":
			// TODO
		case "right":
			// TODO
		case "backspace":
			if len(m.searchText) > 0 {
				m.updateSearchText(m.searchText[:len(m.searchText)-1])
			}
		case "ctrl+c", "esc":
			m.templateSelection = nil
			return m, tea.Quit
		case "enter":
			m.templateSelection = m.filteredTemplates[m.templateCursor]
			return m, tea.Quit
		default:
			m.updateSearchText(m.searchText + keypress)
			m.templateSelection = m.filteredTemplates[m.templateCursor]
		}
	}
	return m, nil
}

func (m *model) searchView() string {
	var s strings.Builder
	s.WriteString(subduedStyle.Render("↑ / ↓ to navigate; <enter> to select and quit"))
	s.WriteString("\n\n")
	start, end := m.pager.GetSliceBounds(len(m.filteredTemplates))
	height := end - start
	for i, template := range m.filteredTemplates[start:end] {
		cursor := " "
		if i == m.templateCursor%height {
			cursor = ">"
		}

		abbreviatedSummary := template.Summary
		if len(abbreviatedSummary) > 75 {
			abbreviatedSummary = abbreviatedSummary[0:72] + "..."
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, abbreviatedSummary))
	}
	if m.pager.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", perPage-m.pager.ItemsOnPage(len(m.filteredTemplates))+1))
		s.WriteString(" " + m.pager.View())
	}
	s.WriteString("\nSearch Text: ")
	s.WriteString(m.searchText)
	s.WriteString("\n")
	return lipgloss.NewStyle().
		Width(80).
		Height(perPage * 2).
		Padding(1).
		Render(s.String())
}

func (m *model) displayView() string {
	if m.templateSelection == nil {
		return "null"
	}
	srcMd := m.templateSelection.String()
	return lipgloss.NewStyle().
		Width(80).
		Height(perPage * 2).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("227")).
		BorderLeft(true).BorderTop(false).BorderRight(false).BorderBottom(false).
		Render(srcMd)
}

func (m *model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.searchView(), m.displayView())
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
	return m.templateSelection
}
