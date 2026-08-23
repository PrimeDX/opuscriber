## Why

WhatsApp and Telegram both deliver voice notes as OGG-container/Opus-codec audio, but converting these into usable text transcripts or subtitle files requires multiple manual steps. This change builds a self-contained, offline Docker CLI that takes raw .opus/.ogg/.oga voice notes and produces two outputs per file: a clean paragraph-reflowed plain-text transcript and a DaVinci-Resolve-ready SRT subtitle file — all running locally with no cloud API calls.

## What Changes

- New Go CLI binary (`opuscriber`) with two modes: non-interactive CLI (default, pipeline-safe) and Bubble Tea TUI (`--interactive`)
- Multi-stage Dockerfile that builds whisper-cli from source (cross-platform for amd64 + arm64), then packages the Go binary + ffmpeg + whisper-cli into a slim Ubuntu runtime image
- Pipeline per audio file: ffmpeg decodes to 16kHz mono WAV → whisper-cli transcribes → Go binary reflows raw segment-by-segment text into clean paragraphs
- SRT output is produced directly by whisper-cli (correct for DaVinci Resolve import); no post-processing needed
- Full TUI: file picker (bubbles), progress bars per file, spinner, styled output with lipgloss (purple accent theme)
- Idempotent: skips files that already have both .txt and .srt outputs
- No telemetry, no API calls, fully local/offline

## Capabilities

### New Capabilities
- `audio-transcription`: Transcribe .opus/.ogg/.oga voice notes into text using a local whisper.cpp model. ffmpeg decodes to 16kHz mono WAV, whisper-cli runs transcription, Go binary reflows raw segment text into clean paragraphs. Configurable language and model size.
- `subtitle-generation`: Generate DaVinci-Resolve-ready .srt subtitle files from transcription output. Uses whisper-cli's built-in SRT output directly.
- `cli-interface`: Two operation modes: a non-interactive CLI (default, clean stdout/stderr, safe in pipelines) and an interactive Bubble Tea TUI (--interactive flag or `tui` subcommand).
- `docker-packaging`: Multi-architecture Docker image (linux/amd64, linux/arm64) packaging the Go binary, ffmpeg, whisper-cli, and model download script. Model persisted via mounted volume.

### Modified Capabilities

None — this is a new project with no existing specs.

## Impact

- New repository at /Users/adrian/src/primedx/opuscriber
- Dependencies: Go 1.27, cobra, fang, bubbletea, bubbles, lipgloss, huh, go-isatty, ffmpeg, whisper.cpp (source-built)
- Docker: multi-stage build with buildx for cross-platform support
- No external services, no API keys, no telemetry
