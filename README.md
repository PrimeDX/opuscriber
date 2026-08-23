# Opuscriber

Transcribe WhatsApp/Telegram voice notes (.opus/.ogg/.oga) into clean plain-text
transcripts and DaVinci-Resolve-ready SRT subtitle files — fully offline, no cloud
APIs, no telemetry.

```
docker run --rm -v $(pwd)/models:/models ghcr.io/ggml-org/whisper.cpp:main \
  sh -c "download-ggml-model.sh medium /models"

docker run --rm \
  -v $(pwd)/in:/audio/in \
  -v $(pwd)/out:/audio/out \
  -v $(pwd)/models:/models \
  opuscriber --lang pt --model medium
```

## Features

- **Two modes**: Non-interactive CLI (default, pipeline-safe) and interactive Bubble
  Tea TUI (`--interactive` / `-i`)
- **Idempotent**: Skips files that already have both .txt and .srt outputs
- **Text reflow**: Joins raw whisper segment lines into clean paragraphs
- **SRT for DaVinci Resolve**: whisper-cli produces standard SRT files that import
  directly via File > Import > Subtitles
- **Multi-arch**: Native linux/amd64 and linux/arm64 images
- **Portuguese default**: `--lang pt` — override with any whisper-supported language

## Quickstart

### 1. Build the image (or pull from registry)

```bash
make docker-build
```

### 2. Download a model

```bash
make model-download MODEL=medium
```

### 3. Transcribe audio

Place .opus/.ogg/.oga files in `./in/`:

```bash
make docker-run ARGS="--lang pt --model medium"
```

Output: `./out/<name>.txt` (clean text) + `./out/<name>.srt` (subtitles).

### Interactive TUI

```bash
make docker-run-interactive
```

Requires `-it` — opuscriber will print an error and exit if TTY is missing.

## Usage

```
opuscriber [flags]

Flags:
  -l, --lang string       Language code (default "pt")
  -m, --model string      Model size: tiny, base, small, medium, large (default "medium")
      --input string      Input audio directory (default "/audio/in")
      --output string     Output directory (default "/audio/out")
      --models string     Model storage directory (default "/models")
  -i, --interactive       Launch interactive TUI
  -h, --help              Help output (styled)
      --version           Version info

Subcommands:
  tui                     Launch the interactive TUI directly
```

### Supported formats

| Extension | Source |
|-----------|--------|
| `.opus`   | WhatsApp, Telegram |
| `.ogg`    | WhatsApp, Telegram |
| `.oga`    | WhatsApp, Telegram (alternative extension) |

All three use the same OGG-container/Opus-codec format.

## Architecture

```
.opus/.ogg/.oga
     │
     ▼
  ffmpeg ──→ 16kHz mono WAV
     │
     ▼
  whisper-cli ──→ raw .txt (segment per line) + .srt (timestamps)
     │
     ▼
  Go reflow ──→ clean paragraphs → <name>.txt
                                   <name>.srt (passthrough)
```

- **Go** orchestrates via `os/exec` — no CGo, no bindings
- **whisper.cpp** is compiled from source in the Docker image (native per arch)
- **Bubble Tea** powers the TUI with file picker, progress bars, lipgloss styling
- **Cobra + Fang** for the CLI layer (styled help, --version, completions)

## Container volumes

| Mount | Purpose |
|-------|---------|
| `/audio/in` | Input .opus/.ogg/.oga files |
| `/audio/out` | Output .txt and .srt files |
| `/models` | Whisper ggml model files (persistent across runs) |

## Multi-arch builds

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t opuscriber .
```

whisper.cpp is compiled from source for each target architecture, producing native
binaries on both platforms.

## License

MIT
