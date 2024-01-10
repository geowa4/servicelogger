package list

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"os"
	"time"
)

type ServiceLogResponse struct {
	Kind        string        `json:"kind"`
	Page        int           `json:"page"`
	Size        int           `json:"size"`
	Total       int           `json:"total"`
	ServiceLogs []*ServiceLog `json:"items"`
}

type ServiceLog struct {
	ClusterId     string    `json:"cluster_id"`
	ClusterUuid   string    `json:"cluster_uuid"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	Desc          string    `json:"description"`
	EventStreamId string    `json:"event_stream_id"`
	Href          string    `json:"href"`
	Id            string    `json:"id"`
	InternalOnly  bool      `json:"internal_only"`
	Kind          string    `json:"kind"`
	LogType       string    `json:"log_type"`
	ServiceName   string    `json:"service_name"`
	Severity      string    `json:"severity"`
	Summary       string    `json:"summary"`
	Timestamp     time.Time `json:"timestamp"`
	Username      string    `json:"username"`
}

func (s *ServiceLog) FilterValue() string {
	internalOrExternal := "external"
	if s.InternalOnly {
		internalOrExternal = "internal"
	}
	return fmt.Sprintf(
		"%s\n%s%s\n\n%s\n%s\n%s\n%s\n%s\n%s",
		s.Title(),
		s.Description(),
		s.CreatedBy,
		s.Severity,
		s.LogType,
		internalOrExternal,
		s.ClusterId,
		s.ClusterUuid)
}

func (s *ServiceLog) Title() string {
	return fmt.Sprintf("%s (%s)", s.Summary, s.ServiceName)
}

func (s *ServiceLog) Description() string {
	return s.Desc
}

func (s *ServiceLog) Markdown() string {
	description := s.Desc
	if description == "" {
		description = "_empty description_"
	}
	return fmt.Sprintf(
		"# [%s] %s\n\n%s\n\n_Created at %s by %s_",
		s.ServiceName,
		s.Summary,
		description,
		s.CreatedAt,
		s.CreatedBy,
	)
}

var (
	verticalPadding   = 1
	horizontalPadding = 2
	paddingStyle      = lipgloss.NewStyle().Padding(verticalPadding, horizontalPadding)
)

type model struct {
	serviceLogs        []*ServiceLog
	selectedServiceLog *ServiceLog
	totalCount         int

	list list.Model

	windowWidth  int
	windowHeight int
}

func initialModel(slResponse *ServiceLogResponse) *model {
	items := make([]list.Item, len(slResponse.ServiceLogs))
	for i, sl := range slResponse.ServiceLogs {
		items[i] = sl
	}
	d := list.NewDefaultDelegate()
	l := list.New(items, d, 0, 0)
	l.Title = "Service Logs"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)
	return &model{
		serviceLogs:        slResponse.ServiceLogs,
		selectedServiceLog: slResponse.ServiceLogs[0],
		totalCount:         slResponse.Total,

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
	if sl, ok := item.(*ServiceLog); ok {
		m.selectedServiceLog = sl
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
		paddingStyle.Width(m.getPaneWidth()).Render(m.selectedServiceLog.Markdown()),
	)
}

func Program(slResponseBytes []byte) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	slResponse := ServiceLogResponse{}
	err := json.Unmarshal(slResponseBytes, &slResponse)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not parse input: %v", err)
		os.Exit(1)
	}
	tm, err := tea.NewProgram(initialModel(&slResponse), tea.WithOutput(os.Stderr), tea.WithAltScreen()).Run()
	if err != nil {
		return
	}

	if m, ok := tm.(*model); ok {
		if md, mdErr := glamour.Render(m.selectedServiceLog.Markdown(), "dark"); mdErr == nil {
			fmt.Println(md)
		}
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "received unexpected model type from program: %v\n", err)
		os.Exit(1)
		return
	}

}
