package editor

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/muesli/termenv"
	"os"
	"strings"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

type model struct {
	baseTemplate     *templates.Template
	filledInTemplate *templates.Template
	variables        []string
	inputs           []textinput.Model
	focusIndex       int
	cursorMode       cursor.Mode
}

func initialModel(template *templates.Template) *model {
	variables := template.GetVariables()
	inputs := make([]textinput.Model, len(variables))
	for i, match := range variables {
		input := textinput.New()
		input.Cursor.Style = cursorStyle
		input.CharLimit = 80
		input.Placeholder = match
		if i == 0 {
			input.Focus()
			input.PromptStyle = focusedStyle
			input.TextStyle = focusedStyle
		}
		inputs[i] = input
	}
	m := &model{
		baseTemplate: template,
		filledInTemplate: &templates.Template{
			Severity:      template.Severity,
			ServiceName:   template.ServiceName,
			Summary:       template.Summary,
			Description:   template.Description,
			InternalOnly:  template.InternalOnly,
			EventStreamId: template.EventStreamId,
			Tags:          append(make([]string, 0), template.Tags...),
			SourcePath:    template.SourcePath,
		},
		variables:  variables,
		inputs:     inputs,
		focusIndex: 0,
	}
	return m
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}
			fallthrough
		case "tab", "down", "shift+tab", "up":
			if keypress == "tab" || keypress == "down" || keypress == "enter" {
				m.focusIndex++
			} else {
				m.focusIndex--
			}
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	summary := m.baseTemplate.Summary
	description := m.baseTemplate.Description
	for i, variable := range m.variables {
		value := m.inputs[i].Value()
		if value == "" {
			continue
		}
		summary = strings.ReplaceAll(summary, variable, value)
		description = strings.ReplaceAll(description, variable, value)
	}
	m.filledInTemplate.Summary = summary
	m.filledInTemplate.Description = description
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) markdownView() string {
	md, err := glamour.Render(m.filledInTemplate.String(), "dark")
	if err != nil {
		return ""
	}
	return md
}

func (m *model) formView() string {
	var s strings.Builder

	for i := range m.inputs {
		s.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			s.WriteRune('\n')
		}
	}

	button := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render("[ Submit ]")
	}
	s.WriteString(fmt.Sprintf("\n\n%s", button))

	return s.String()
}

func (m *model) View() string {
	mdView := m.markdownView()
	formView := m.formView()
	return lipgloss.JoinVertical(lipgloss.Left, mdView, formView)
}

func Program(template *templates.Template) *templates.Template {
	lipgloss.SetColorProfile(termenv.TrueColor)
	tm, err := tea.NewProgram(initialModel(template), tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error running editor program: %v\n", err)
		os.Exit(1)
	}
	m, ok := tm.(*model)
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "received unexpected model type from program: %v\n", err)
		os.Exit(1)
	}
	return m.filledInTemplate
}
