package search

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
