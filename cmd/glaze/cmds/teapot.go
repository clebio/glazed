package cmds

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type sessionState uint

const (
	defaultTime              = time.Minute
	spinnerView sessionState = iota
)

var TeapotCmd = &cobra.Command{
	Use:   "teapot [flags] file [file...]",
	Short: "Teapot: a porcelain for Glazed using Bubbletea",
	// Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(newModel(defaultTime))

		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

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
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type model struct {
	state   sessionState
	timer   timer.Model
	spinner spinner.Model
}

func newModel(timeout time.Duration) model {
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
		focusedModelStyle.Render(fmt.Sprintf("%4s", m.timer.View())),
		modelStyle.Render(m.spinner.View()))
}
