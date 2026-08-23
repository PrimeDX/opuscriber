## 1. Project Scaffold

- [x] 1.1 Git init, go mod init, .gitignore, initial commit

## 2. CLI Layer

- [x] 2.1 main.go entrypoint delegating to cobra root command
- [x] 2.2 cmd/root.go: cobra root wrapped with fang, all flags (--lang, --model, --input, --output, --models, --interactive)
- [x] 2.3 cmd/tui.go: --interactive flag handler with TTY check; fails with message if no TTY

## 3. Pipeline

- [x] 3.1 pipeline/transcribe.go: ProcessFile() — ffmpeg decode to 16kHz WAV, whisper-cli invocation with -osrt -otxt, temp file cleanup
- [x] 3.2 pipeline/txt.go: ReflowTxt() pure function — join segment lines, preserve paragraph breaks, trim whitespace
- [x] 3.3 pipeline/txt_test.go: table-driven tests for ReflowTxt (empty, single line, multi-line join, paragraph breaks, trim)
- [x] 3.4 Batch directory scan: iterate audio dir, skip files with both .txt+.srt, call ProcessFile per file

## 4. TUI

- [x] 4.1 tui/app.go: Bubble Tea model with state machine (idle → browsing → transcribing → done)
- [x] 4.2 File picker integration (bubbles filepicker), processing queue with progress bars, results summary
- [x] 4.3 Lipgloss styling: accent color #7C3AED, rounded borders, consistent spacing

## 5. Docker & Build

- [x] 5.1 Dockerfile: multi-stage — whisper.cpp source build, Go build, Ubuntu runtime with ffmpeg
- [x] 5.2 Makefile: build, run, run-interactive, test, docker-build, model-download targets
- [x] 5.3 README.md: what it does, quickstart, CLI and TUI usage, flags, multi-arch note, SRT note

## 6. Verification

- [x] 6.1 go build compiles, go test passes
- [x] 6.2 docker build succeeds
- [x] 6.3 Model download + end-to-end transcription with real audio file
- [x] 6.4 Idempotency: re-run skips already-processed files
- [x] 6.5 Interactive TTY check: docker run without -it prints error
