package labels

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/geowa4/servicelogger/pkg/templates"
	"os"
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

type model struct {
	searchText        string
	templates         []*templates.Template
	templateCursor    int
	templateSelection *templates.Template
}

func initialModel() model {
	allTemplates := make([]*templates.Template, 0)
	templates.WalkTemplates(func(template *templates.Template) {
		allTemplates = append(allTemplates, template)
	})
	return model{
		searchText:        "",
		templates:         allTemplates,
		templateCursor:    0,
		templateSelection: nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			// The "up" and "k" keys move the cursor up
			if m.templateCursor > 0 {
				m.templateCursor--
			}
		case "down", "j":
			// The "down" and "j" keys move the cursor down
			if m.templateCursor < len(m.templates)-1 {
				m.templateCursor++
			}
		case "enter", " ":
			m.templateSelection = m.templates[m.templateCursor]
			//return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := ""
	for i, template := range m.templates {
		if i > 10 {
			break
		}
		cursor := " "
		if m.templateCursor == i {
			cursor = ">"
		}
		checked := " "
		if m.templateSelection != nil && template.Summary == m.templateSelection.Summary {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, template.Summary)
	}
	s += "\nPress q to quit.\n"
	return s
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
