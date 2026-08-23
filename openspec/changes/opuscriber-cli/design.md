## Context

See proposal.md - Why and What Changes. This is a greenfield Go project wrapping whisper.cpp and ffmpeg in a Dockerized CLI with two modes: non-interactive batch and Bubble Tea TUI.

## Goals / Non-Goals

**Goals:**
- CLI processes a directory of .opus/.ogg/.oga files in one command
- Pipeline: ffmpeg → whisper-cli → txt reflow → .txt + .srt output
- Multi-arch Docker image (linux/amd64 + linux/arm64) from source-built whisper.cpp
- TUI: file picker, per-file progress bars, results summary
- Idempotent: skip files already having both outputs
- Pure-function ReflowTxt with table-driven tests

**Non-Goals:**
- Not a server/daemon — single-shot CLI only
- No cloud API or telemetry
- No GPU acceleration (CPU-only whisper.cpp via Docker)
- No speaker diarization or streaming

## Decisions

| Decision | Choice | Rationale | Alternatives |
|----------|--------|-----------|-------------|
| Language | Go | Cross-compiles trivially (CGO_ENABLED=0). os/exec for subprocess. Bubble Tea for TUI. | Python → larger image, no mature TUI. Rust → steeper curve. |
| Whisper runtime | whisper-cli subprocess | Battle-tested. No CGo linking. Voice notes are short, subprocess overhead negligible. | CGo bindings → coupling to specific whisper.cpp version. |
| TUI framework | Bubble Tea | Standard Go TUI with filepicker, progress bars, lipgloss styling. | None — Bubble Tea is the standard. |
| Cobra wrapper | Fang | Auto-styled help, --version, shell completions, man pages. Color-aware for pipe safety. | Plain cobra → manual styling. |
| Docker base | Ubuntu 22.04 | Prebuilt whisper.cpp image only ships amd64 + SIGILL on ARM64. Source build fixes both. | Debian slim, Alpine (needs glibc/gomp). |
| Model storage | Volume mount | ~1.5GB models. Download once, reuse across runs. | Bundling in image → image grows 1.5GB. |
| SRT handling | Pass-through | whisper-cli -osrt produces correct DaVinci Resolve SRT. | Hand-written generator → duplicating tested code. |

## Risks / Trade-offs

- **Build time**: First Docker build compiles whisper.cpp from source (~5 min). Docker layer caching mitigates.
- **whisper-cli flag drift**: If upstream changes -osrt/-otxt flags. Build from git HEAD tracks current.
- **No GPU**: CPU-only. Fine for voice notes (2-60s). Medium model ~30s real-time per audio minute.
- **Temp disk**: Each file decoded to WAV (~10MB/min). Immediately deleted after transcription.
