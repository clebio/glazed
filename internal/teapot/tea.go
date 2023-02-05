package teapot

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func init() {
	// TeapotCmd.Flags().SortFlags = false
	// cli.AddOutputFlags(TeapotCmd, cli.NewOutputFlagsDefaults())
	// cli.AddTemplateFlags(TeapotCmd, cli.NewTemplateFlagsDefaults())
	// cli.AddFieldsFilterFlags(TeapotCmd, &cli.FieldsFilterFlagsDefaults{
	// 	Fields:      "path,Title,SectionType,Slug,Commands,Flags,Topics,IsTopLevel,ShowPerDefault",
	// 	Filter:      "",
	// 	SortColumns: false,
	// })
}

var (
	modelStyle = lipgloss.NewStyle().
		Width(30).
		Height(10).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.HiddenBorder())
)

type model struct {
	timer   timer.Model
	spinner spinner.Model
}

const (
	defaultTime = 10 * time.Second
)

// NewModel returns a new Bubbletea model with default values.
func NewModel(timeout time.Duration) model {
	m := model{}
	m.timer = timer.New(defaultTime)
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.MiniDot
	return m
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.timer.Init(), m.spinner.Tick)
}

// Views return a string based on data in the model. That string which will be
// rendered to the terminal.
func (m model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		modelStyle.
			BorderStyle(lipgloss.HiddenBorder()).
			Render(fmt.Sprintf("%4s", m.timer.View())),
		modelStyle.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")).
			Render(m.spinner.View()),
		modelStyle.
			BorderStyle(lipgloss.HiddenBorder()).
			Render("stub"))
}
