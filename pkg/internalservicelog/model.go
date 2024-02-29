package internalservicelog

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	inputForm   *huh.Form
	confirmForm *huh.Form

	slSummary    string
	slBody       string
	confirmation bool
}

func initialModel() *model {
	m := &model{}
	m.createNewForms()
	return m
}

func (m *model) createNewForms() {
	m.inputForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Value(&m.slSummary).Title("Brief Summary").Validate(func(s string) error {
				if len(s) == 0 {
					return errors.New("summary cannot be empty")
				}
				return nil
			}),
			huh.NewText().Value(&m.slBody).Title("Body").CharLimit(1000),
		).Title("Internal Service Log (Markdown Supported)"),
	)
	m.confirmForm = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Value(&m.confirmation).
				Title("Send this service log?").
				Negative("Edit").
				Affirmative("Send"),
		),
	)
}

func (m *model) Markdown() string {
	return fmt.Sprintf("# %s\n\n%s\n", m.slSummary, m.slBody)
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.inputForm.Init(), m.confirmForm.Init())
}

type FailedConfirmationMsg struct{}

func FailedConfirmation() tea.Msg {
	return FailedConfirmationMsg{}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmdBatch := make([]tea.Cmd, 0, 1)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case FailedConfirmationMsg:
		m.createNewForms()
		cmdBatch = append(cmdBatch, m.inputForm.Init(), m.confirmForm.Init())
	}
	if m.inputForm.State != huh.StateCompleted {
		inputForm, cmd := m.inputForm.Update(msg)
		cmdBatch = append(cmdBatch, cmd)
		if f, ok := inputForm.(*huh.Form); ok {
			m.inputForm = f
		}
	} else if m.confirmForm.State != huh.StateCompleted {
		confirmForm, cmd := m.confirmForm.Update(msg)
		cmdBatch = append(cmdBatch, cmd)
		if f, ok := confirmForm.(*huh.Form); ok {
			m.confirmForm = f
		}
		if m.confirmForm.State == huh.StateCompleted {
			if !m.confirmation {
				cmdBatch = append(cmdBatch, FailedConfirmation)
			} else {
				cmdBatch = append(cmdBatch, tea.Quit)
			}
		}
	}
	return m, tea.Batch(cmdBatch...)
}

func (m *model) View() string {
	if m.inputForm.State == huh.StateCompleted {
		md := m.Markdown()
		renderedMd, err := glamour.Render(md, "notty")
		if err != nil {
			renderedMd = md
		}
		return lipgloss.JoinVertical(
			lipgloss.Left,
			renderedMd,
			m.confirmForm.View(),
		)
	}
	return m.inputForm.View()
}
