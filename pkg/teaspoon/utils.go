package teaspoon

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"os"
)

func Program(m tea.Model) (tea.Model, error) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	return tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
}

func renderMarkdown(md string, rendererOptions ...glamour.TermRendererOption) string {
	renderer, err := glamour.NewTermRenderer(rendererOptions...)
	if err != nil {
		return md
	}
	renderedMd, err := renderer.Render(md)
	if err != nil {
		return md
	}
	return renderedMd
}

func RenderMarkdown(md string) string {
	return renderMarkdown(
		md,
		glamour.WithStandardStyle("notty"),
	)
}

func RenderMarkdownWithWordWrap(md string, wordWrap int) string {
	return renderMarkdown(
		md,
		glamour.WithStandardStyle("notty"),
		glamour.WithWordWrap(wordWrap),
	)
}
