package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"opuscriber/pipeline"
)

var (
	lang        string
	modelSize   string
	inputDir    string
	outputDir   string
	modelsDir   string
	interactive bool
)

var rootCmd = &cobra.Command{
	Use:   "opuscriber",
	Short: "Transcribe voice notes to clean text + SRT subtitles",
	Long: `A fully local, offline CLI that transcribes .opus/.ogg/.oga voice notes
(WhatsApp/Telegram) into clean plain-text transcripts and DaVinci-Resolve-ready
SRT subtitle files. No cloud APIs, no telemetry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if interactive {
			return runTUI()
		}
		return runCLI()
	},
}

func Execute() {
	ctx := context.Background()
	if err := fang.Execute(ctx, rootCmd,
		fang.WithVersion("0.1.0"),
	); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&lang, "lang", "l", "auto", "Language code for transcription (auto for auto-detect)")
	rootCmd.Flags().StringVarP(&modelSize, "model", "m", "medium", "Whisper model size (tiny, base, small, medium, large)")
	rootCmd.Flags().StringVar(&inputDir, "input", "/audio/in", "Input audio directory")
	rootCmd.Flags().StringVar(&outputDir, "output", "/audio/out", "Output directory for .txt and .srt")
	rootCmd.Flags().StringVar(&modelsDir, "models", "/models", "Model storage directory")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Launch interactive TUI")

	rootCmd.AddCommand(tuiCmd)
}

func runCLI() error {
	modelPath := filepath.Join(modelsDir, "ggml-"+modelSize+".bin")

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("reading input dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".opus" || ext == ".ogg" || ext == ".oga" {
			files = append(files, filepath.Join(inputDir, e.Name()))
		}
	}

	if len(files) == 0 {
		fmt.Println("no audio files found")
		return nil
	}

	for _, f := range files {
		base := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		txtPath := filepath.Join(outputDir, base+".txt")
		srtPath := filepath.Join(outputDir, base+".srt")

		if _, errT := os.Stat(txtPath); errT == nil {
			if _, errS := os.Stat(srtPath); errS == nil {
				fmt.Printf("skip %s (already processed)\n", base)
				continue
			}
		}

		fmt.Printf("transcribing %s...\n", base)
		if err := pipeline.ProcessFile(f, modelPath, lang, outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "error processing %s: %v\n", base, err)
			continue
		}
		fmt.Printf("  done: %s.txt, %s.srt\n", base, base)
	}

	return nil
}

func runTUI() error {
	if !isatty.IsTerminal(os.Stdin.Fd()) || !isatty.IsTerminal(os.Stdout.Fd()) {
		return fmt.Errorf("interactive mode requires a TTY; add -it to your docker run command")
	}
	return startTUI()
}
