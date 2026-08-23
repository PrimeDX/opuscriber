package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"opuscriber/pipeline"
)

var (
	accent     = lipgloss.Color("#7C3AED")
	text       = lipgloss.Color("#C0CAF5")
	textDim    = lipgloss.Color("#565F89")
	textGreen  = lipgloss.Color("#9ECE6A")
	textRed    = lipgloss.Color("#F7768E")
	textYellow = lipgloss.Color("#E0AF68")

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(accent).MarginLeft(2).MarginTop(1).MarginBottom(1)
	helpStyle   = lipgloss.NewStyle().Foreground(textDim).MarginLeft(2).MarginTop(1).MarginBottom(1)
	statusStyle = lipgloss.NewStyle().Foreground(text).MarginLeft(2)
	bannerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accent).Padding(0, 2).MarginLeft(2).MarginRight(2).MarginTop(1).Width(60)
	progressSty = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).MarginBottom(1)
	btnStyle    = lipgloss.NewStyle().Foreground(accent).Bold(true).MarginLeft(2)
	checkStyle  = lipgloss.NewStyle().Foreground(textGreen).MarginLeft(2)
	doneStyle   = lipgloss.NewStyle().Foreground(textYellow).Bold(true).MarginLeft(2).MarginTop(1)
	divider     = lipgloss.NewStyle().Foreground(textDim).Width(64).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(textDim).MarginLeft(2).Render("")
)

type state int

const (
	idle state = iota
	browsing
	transcribing
	done
)

type fileItem struct {
	name string
	path string
}

type jobResult struct {
	name string
	err  error
}

type uiModel struct {
	state       state
	fp          filepicker.Model
	spin        spinner.Model
	pbar        progress.Model
	queue       []fileItem
	results     []jobResult
	doneMsg     string
	fileIdx     int
	progressVal float64

	inputDir  string
	outputDir string
	modelsDir string
	model     string
	lang      string
}

func (m uiModel) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, tea.SetWindowTitle("Opuscriber"))
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Let filepicker process messages in browsing state
	var fpCmd tea.Cmd
	if m.state == browsing {
		m.fp, fpCmd = m.fp.Update(msg)
		cmds = append(cmds, fpCmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.fp.Height = msg.Height - 12

	case tea.KeyMsg:
		switch m.state {
		case idle:
			switch msg.String() {
			case "enter", " ":
				m.state = browsing
				m.fp.CurrentDirectory = m.inputDir
				return m, tea.Batch(m.spin.Tick, m.fp.Init())
			case "q", "ctrl+c":
				return m, tea.Quit
			}

		case browsing:
			// After filepicker processes enter, check if a file was selected
			if msg.String() == "enter" && m.fp.Path != "" {
				path := m.fp.Path
				name := filepath.Base(path)
				ext := strings.ToLower(filepath.Ext(name))
				if ext == ".opus" || ext == ".ogg" || ext == ".oga" {
					already := false
					for _, item := range m.queue {
						if item.path == path {
							already = true
							break
						}
					}
					if !already {
						m.queue = append(m.queue, fileItem{name: name, path: path})
					}
				}
				m.fp.Path = ""
			}
			if msg.String() == "tab" && len(m.queue) > 0 {
				m.state = transcribing
				m.fileIdx = 0
				m.progressVal = 0
				return m, m.startJob()
			}
			if msg.String() == "esc" {
				m.state = idle
			}
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case transcribing:
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case done:
			switch msg.String() {
			case "enter", "r":
				m.state = idle
				m.queue = nil
				m.results = nil
				m.fileIdx = 0
				m.progressVal = 0
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

	case progressMsg:
		if m.state != transcribing {
			break
		}
		m.progressVal = float64(msg)
		if m.progressVal >= 1.0 {
			return m, m.finishJob()
		}
		return m, tea.Batch(
			func() tea.Msg { time.Sleep(100 * time.Millisecond); return progressMsg(m.progressVal + 0.05) },
			m.spin.Tick,
		)

	case jobDoneMsg:
		if m.state != transcribing {
			break
		}
		m.results = append(m.results, jobResult{name: msg.name, err: msg.err})
		m.fileIdx++
		if m.fileIdx >= len(m.queue) {
			success, fail := 0, 0
			for _, r := range m.results {
				if r.err != nil {
					fail++
				} else {
					success++
				}
			}
			if fail > 0 {
				m.doneMsg = fmt.Sprintf("%d transcribed, %d failed", success, fail)
			} else {
				m.doneMsg = fmt.Sprintf("%d files transcribed successfully!", success)
			}
			m.state = done
			return m, tea.Batch(m.spin.Tick)
		}
		m.progressVal = 0
		return m, m.startJob()
	}

	// Component updates
	var cmd tea.Cmd
	m.spin, cmd = m.spin.Update(msg)
	cmds = append(cmds, cmd)
	pm, pcmd := m.pbar.Update(msg)
	if p, ok := pm.(progress.Model); ok {
		m.pbar = p
	}
	cmds = append(cmds, pcmd)

	return m, tea.Batch(cmds...)
}

type progressMsg float64
type jobDoneMsg struct {
	name string
	err  error
}

func (m *uiModel) startJob() tea.Cmd {
	item := m.queue[m.fileIdx]
	return func() tea.Msg {
		modelPath := filepath.Join(m.modelsDir, "ggml-"+m.model+".bin")
		err := pipeline.ProcessFile(item.path, modelPath, m.lang, m.outputDir)
		return jobDoneMsg{name: item.name, err: err}
	}
}

func (m *uiModel) finishJob() tea.Cmd {
	return func() tea.Msg {
		return jobDoneMsg{name: m.queue[m.fileIdx].name, err: nil}
	}
}

func (m uiModel) View() string {
	var b strings.Builder
	b.WriteString("\n")

	switch m.state {
	case idle:
		b.WriteString(titleStyle.Render("Opuscriber"))
		b.WriteString("\n")
		b.WriteString(bannerStyle.Render("Transcribe voice notes → clean text + SRT subtitles\nFully local · Offline · No API calls\nFormats: .opus .ogg .oga (WhatsApp/Telegram)"))
		b.WriteString("\n\n")
		b.WriteString(statusStyle.Render("Press Enter to select audio files"))
		b.WriteString("\n")
		b.WriteString(btnStyle.Render("  [ Enter ] Browse for files"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("q / Ctrl+C to quit"))

	case browsing:
		b.WriteString(titleStyle.Render("Select Audio Files"))
		b.WriteString("\n")
		b.WriteString(m.fp.View())
		b.WriteString("\n")
		b.WriteString(divider)
		b.WriteString("\n")
		b.WriteString(helpStyle.Render(fmt.Sprintf("Selected: %d files  |  Tab → transcribe  |  q quit", len(m.queue))))

	case transcribing:
		b.WriteString(titleStyle.Render("Transcribing"))
		b.WriteString("\n")
		b.WriteString(divider)
		b.WriteString("\n")
		cur := "starting..."
		if m.fileIdx < len(m.queue) {
			cur = m.queue[m.fileIdx].name
		}
		b.WriteString(lipgloss.NewStyle().MarginLeft(2).Render(fmt.Sprintf("%s %s", m.spin.View(), cur)))
		b.WriteString("\n")
		b.WriteString(progressSty.Render(m.pbar.ViewAs(m.progressVal)))
		b.WriteString("\n")
		if len(m.results) > 0 {
			b.WriteString(divider)
			b.WriteString("\n")
			for _, r := range m.results {
				if r.err != nil {
					b.WriteString(lipgloss.NewStyle().Foreground(textRed).MarginLeft(2).Render("✗ " + r.name))
				} else {
					b.WriteString(checkStyle.Render("✓ " + r.name))
				}
				b.WriteString("\n")
			}
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Processing...  q / Ctrl+C to quit"))

	case done:
		b.WriteString(titleStyle.Render("Complete!"))
		b.WriteString("\n")
		b.WriteString(divider)
		b.WriteString("\n")
		for _, r := range m.results {
			if r.err != nil {
				b.WriteString(lipgloss.NewStyle().Foreground(textRed).MarginLeft(2).Render("✗ " + r.name + " — " + r.err.Error()))
			} else {
				b.WriteString(checkStyle.Render("✓ " + r.name))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(divider)
		b.WriteString("\n")
		b.WriteString(doneStyle.Render(m.doneMsg))
		b.WriteString("\n\n")
		b.WriteString(btnStyle.Render("  [ Enter / r ] Start over  •  [ q ] Quit"))
	}

	return b.String()
}

func Run(inputDir, outputDir, modelsDir, model, lang string) error {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(accent)
	s.Spinner = spinner.Dot

	p := progress.New(progress.WithSolidFill("#7C3AED"), progress.WithWidth(50))

	fp := filepicker.New()
	fp.AllowedTypes = []string{".opus", ".ogg", ".oga"}
	fp.Height = 15
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.AutoHeight = false
	fp.DirAllowed = false

	m := uiModel{
		state:     idle,
		fp:        fp,
		spin:      s,
		pbar:      p,
		inputDir:  inputDir,
		outputDir: outputDir,
		modelsDir: modelsDir,
		model:     model,
		lang:      lang,
	}

	prog := tea.NewProgram(m, tea.WithAltScreen())
	_, err := prog.Run()
	return err
}
