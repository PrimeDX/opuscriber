# Opuscriber — Offline Speech-to-Text for Opus Audio

> Local audio transcription and SRT subtitle generation powered by whisper.cpp.
> Single binary. Cross-platform (macOS, Linux). No cloud. No telemetry.

Opuscriber transcribes Opus-encoded audio files (.opus, .ogg, .oga) into clean
plain-text transcripts and DaVinci-Resolve-ready SRT subtitles. Runs entirely on
your machine — download a Whisper model once, then transcribe forever offline.

Built for WhatsApp and Telegram voice notes, works with any Opus audio.

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="demo.gif">
  <img alt="Opuscriber TUI in action: selecting files, transcribing, showing results" src="demo.gif">
</picture>

## Quickstart

```bash
# Requires Go 1.27+, cmake, C compiler, ffmpeg, and whisper-cli:
brew install ffmpeg
# Download whisper-cli for your platform from github.com/ggml-org/whisper.cpp
# Place it at /usr/local/bin/whisper-cli

# Clone and build
git clone https://github.com/primedx/opuscriber.git
cd opuscriber
make build

# Download a Whisper model (e.g. ggml-medium.bin from Hugging Face)
# Place it in ./models/ggml-medium.bin

# Transcribe all audio files
./opuscriber --input ./in --output ./out --models ./models
```

That's it. No Docker, no Python, no cloud setup.

## What Is Opuscriber?

Opuscriber turns spoken-word Opus audio into two files: a clean plain-text
transcript and a timestamped SRT subtitle file importable directly into DaVinci
Resolve, VLC, or any editor that supports SRT.

It's purpose-built for **offline use** — once you download a Whisper model, you
never need the internet again. Perfect for journalists processing field recordings,
legal transcription, podcasters generating subtitles, or anyone who doesn't want
their voice notes uploaded to a cloud API.

## Why Opuscriber?

**Instead of using whisper.cpp directly** — you need ffmpeg to convert Opus to
WAV, manual shell scripts to generate SRT, and you re-process files you already
transcribed. Opuscriber handles all of this in a single command: ffmpeg decode,
built-in SRT generation, and idempotent resume.

**Instead of a cloud speech-to-text API** — you pay per minute, you need internet,
and your audio leaves your machine. Opuscriber is fully local. Download a model
once and transcribe forever with no network. Your data stays on your machine.

**Instead of a desktop GUI app** — you get a polished interface but lose
automation. Opuscriber works as both a terminal UI for interactive use *and* a
headless CLI for scripts, cron jobs, and pipelines.

## Features

- **Single binary**: No Docker, no Python, no shell wrappers — one static binary
- **External whisper-cli**: Uses whisper.cpp for inference — works with Metal (macOS) and CUDA (Linux) accelerated builds
- **Two modes**: Non-interactive CLI (pipeline-safe) and interactive Bubble Tea TUI
- **Idempotent**: Skips files that already have both `.txt` and `.srt` outputs
- **Text reflow**: Joins Whisper segment lines into clean paragraphs
- **SRT subtitles**: Standard format, imports into DaVinci Resolve, VLC, Premiere
- **Cross-platform**: macOS (Apple Silicon + Intel), Linux

## Which Model Should I Use?

Models range from tiny (fast, okay accuracy) to large (slow, near-perfect).

| Model | Disk | RAM (approx) | Speed | Best For |
|---|---|---|---|---|
| `tiny` | 74 MB | ~300 MB | Fastest | Quick tests, low-resource devices |
| `base` | 141 MB | ~500 MB | Fast | Short clips, casual use |
| `small` | 465 MB | ~1 GB | Moderate | Good balance for most users |
| `medium` ★ | 1.4 GB | ~2 GB | Slower | Excellent accuracy for most use |
| `large` | 2.9 GB | ~4 GB | Slowest | Highest accuracy, long recordings |

★ = recommended default. Start with `medium`. Switch to `large` for near-perfect
transcription, or `small` for a speed boost on older machines.

Download models from [Hugging Face](https://huggingface.co/ggml-org) and place
them in your models directory.

## Usage

### Commands

```
opuscriber [flags]    Transcribe audio files
opuscriber tui        Launch interactive TUI
opuscriber help       Show help text
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-l, --lang` | `auto` | Language code (ISO 639-1), or `auto` for auto-detect |
| `-m, --model` | `medium` | Model size: tiny, base, small, medium, large |
| `--input` | `/audio/in` | Input audio directory |
| `--output` | `/audio/out` | Output directory for `.txt` and `.srt` |
| `--models` | `/models` | Model storage directory |
| `-i, --interactive` | | Launch interactive TUI |

### Interactive TUI

```bash
opuscriber --interactive
# or
opuscriber tui
```

The TUI lets you browse and select audio files, track transcription progress
with a progress bar and live per-file checkmarks, then view results.

## Prerequisites

- **ffmpeg** — installed and on `PATH` (used for decoding audio to 16kHz WAV)
- **whisper-cli** — at `/usr/local/bin/whisper-cli` (the [whisper.cpp](https://github.com/ggml-org/whisper.cpp) command-line tool)
- **Whisper model** — a `.bin` model file (e.g. `ggml-medium.bin`) in your models directory

## Installation

### Build from source

```bash
git clone https://github.com/primedx/opuscriber.git
cd opuscriber
make build
```

**Prerequisites**: Go 1.27+, cmake, and a C compiler (Clang on macOS, GCC on Linux).

The build clones whisper.cpp at a pinned version and compiles the static library
via CMake. First build takes ~2 minutes.

### Homebrew (coming soon)

```bash
brew install primedx/tap/opuscriber
```

### Prebuilt binaries (coming soon)

Platform-specific binaries will be attached to [GitHub Releases](https://github.com/primedx/opuscriber/releases).

## Pipeline

```
.opus/.ogg/.oga
     │
     ▼
  ffmpeg decode ──→ 16kHz mono WAV
     │
     ▼
  whisper-cli ──→ timestamped text segments
     │
     ├──→ TXT: reflowed paragraphs
     └──→ SRT: subtitle file
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the full build chain,
GPU backends, and project structure.

## Environment

| Variable | Default | Description |
|---|---|---|
| `OPUSCRIBER_MODELS` | `/models` | Override model storage directory |

## License

MIT