package teaspoon

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"os"
)

func Program(m tea.Model) (tea.Model, error) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	return tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
}
