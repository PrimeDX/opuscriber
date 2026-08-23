package cmd

import (
	"github.com/spf13/cobra"
	"opuscriber/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startTUI()
	},
}

func startTUI() error {
	return tui.Run(inputDir, outputDir, modelsDir, modelSize, lang)
}
