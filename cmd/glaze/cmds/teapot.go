package cmds

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/wesen/glazed/internal/teapot"
)

type sessionState uint

const (
	defaultTime              = 10 * time.Second
	spinnerView sessionState = iota
)

var TeapotCmd = &cobra.Command{
	Use:   "teapot [flags] file [file...]",
	Short: "Teapot: a porcelain for Glazed using Bubbletea",
	// Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(teapot.NewModel(defaultTime))

		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
