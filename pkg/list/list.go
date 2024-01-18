package list

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/geowa4/servicelogger/pkg/ocm"
)

type ServiceLogView struct {
	Log ocm.ServiceLog
}

func (s ServiceLogView) FilterValue() string {
	internalOrExternal := "external"
	if s.Log.InternalOnly {
		internalOrExternal = "internal"
	}
	return fmt.Sprintf(
		"%s\n%s%s\n\n%s\n%s\n%s\n%s\n%s\n%s",
		s.Title(),
		s.Description(),
		s.Log.CreatedBy,
		s.Log.Severity,
		s.Log.LogType,
		internalOrExternal,
		s.Log.ClusterId,
		s.Log.ClusterUuid)
}

func (s ServiceLogView) Title() string {
	return fmt.Sprintf("%s (%s)", s.Log.Summary, s.Log.ServiceName)
}

func (s ServiceLogView) Description() string {
	return s.Log.Desc
}

func markdown(log ocm.ServiceLog) string {
	description := log.Desc
	if description == "" {
		description = "_empty description_"
	}
	return fmt.Sprintf(
		"# [%s] %s\n\n%s\n\n_Created at %s by %s_",
		log.ServiceName,
		log.Summary,
		description,
		log.CreatedAt,
		log.CreatedBy,
	)
}

var (
	verticalPadding   = 1
	horizontalPadding = 2
	paddingStyle      = lipgloss.NewStyle().Padding(verticalPadding, horizontalPadding)
)

type model struct {
	serviceLogs        []ocm.ServiceLog
	selectedServiceLog ocm.ServiceLog
	totalCount         int

	list list.Model

	windowWidth  int
	windowHeight int
}

func initialModel(serviceLogs []ocm.ServiceLog) *model {
	items := make([]list.Item, len(serviceLogs))
	for i, sl := range serviceLogs {
		items[i] = ServiceLogView{sl}
	}
	d := list.NewDefaultDelegate()
	l := list.New(items, d, 0, 0)
	l.Title = "Service Log List"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)
	l.InfiniteScrolling = true
	return &model{
		serviceLogs:        serviceLogs,
		selectedServiceLog: serviceLogs[0],
		totalCount:         len(serviceLogs),

		list: l,
	}
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

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	item := newListModel.SelectedItem()
	if sl, ok := item.(ServiceLogView); ok {
		m.selectedServiceLog = sl.Log
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
	md := markdown(m.selectedServiceLog)
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

func Program(accessToken, refreshToken string) {
	conn, err := ocm.NewConnection(accessToken, refreshToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse input: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := ocm.NewClient(conn)
	list, err := client.ListServiceLogs("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get serviceLogs: %v", err)
		os.Exit(1)
	}
	if len(list) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "no service logs to view")
		os.Exit(0)
	}
	tm, err := tea.NewProgram(initialModel(list), tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
	if err != nil {
		return
	}

	if m, ok := tm.(*model); ok {
		if md, mdErr := glamour.Render(markdown(m.selectedServiceLog), "notty"); mdErr == nil {
			fmt.Println(md)
		}
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "received unexpected model type from program: %v\n", err)
		os.Exit(1)
		return
	}
}
