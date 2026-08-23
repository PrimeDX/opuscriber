# Opuscriber Architecture

## Pipeline

```
.opus/.ogg/.oga
     │
     ▼
  OGG/Opus decode (pure Go, pion/opus) ──→ 16kHz mono float32 PCM
     │
     ▼
  whisper.cpp (CGo, linked into binary) ──→ timestamped text segments
     │
     ├──→ TXT: reflow segments → clean paragraphs
     └──→ SRT: segment timestamps → subtitle file
```

The entire pipeline runs in one binary. No Docker, no subprocesses, no shell glue.

## Technology Stack

- **Go** drives everything: OGG parsing, Opus decode, whisper inference, text output
- **whisper.cpp** is compiled as a static library (`libwhisper.a`) and linked via CGo
- **pion/opus** provides pure-Go Opus decoding (no C deps for audio)
- **Bubble Tea** powers the TUI with file picker, progress bars, lipgloss styling
- **Cobra + Fang** for the CLI layer (styled help, --version, completions)

## Project Structure

```
opuscriber/
├── cmd/                  # CLI commands (Cobra)
│   ├── root.go           # Flag registration, subcommand wiring, Execute()
│   ├── tui.go            # TUI subcommand registration
│   ├── model.go          # Model subcommands (list, install, installed, remove)
│   ├── config.go         # Config subcommands (show, set)
│   └── transcribe.go     # runCLI() + runTUI() for transcription
├── internal/             # Private packages
│   ├── catalog/          # Embedded model catalog (models.json + Go accessors)
│   ├── config/           # User config file management (~/.config/opuscriber)
│   └── model/            # Model download, scan, sha256 integrity
├── pipeline/             # Transcription pipeline
│   ├── transcribe.go     # Entry point: ProcessFile()
│   ├── decode.go         # OGG container parsing + pion/opus decode
│   ├── whisper.go        # CGo whisper.cpp wrapper (TranscribeAudio)
│   ├── txt.go            # Segment-to-TXT paragraph reflow
│   ├── srt.go            # SRT subtitle generation from segments
│   ├── ..._test.go       # Tests
├── tui/                  # Bubble Tea TUI
│   └── app.go            # TUI application (states: idle, browsing, transcribing, done)
├── demo/                 # Demo assets
│   ├── jfk-moon-speech.opus  # Public domain demo audio (PD, US Government)
│   └── demo.tape         # vhs terminal recording script
├── sample-audio/         # Original audio source files for provenance
├── docs/
│   ├── CONTRIBUTING.md   # Contributor guide (build setup, tests)
│   ├── MAINTAINING.md    # Release process, whisper.cpp version bumps
│   └── ARCHITECTURE.md   # This file
├── Dockerfile            # Optional: containerized build
├── Makefile
└── go.mod / go.sum
```

## Build Chain

The build requires three stages:

1. **Clone whisper.cpp** at a pinned version tag (defined in `Makefile` as `WHISPER_VERSION`). whisper.cpp is **not** a git submodule — it lives in a `.gitignore`d `whisper.cpp/` directory.

2. **Build libwhisper.a** via CMake:
   ```bash
   cmake -S whisper.cpp -B whisper.cpp/build \
     -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release
   cmake --build whisper.cpp/build -j$(nproc)
   ```
   This produces `libwhisper.a`, `libggml.a`, and backend `.a` files (`libggml-metal.a`, `libggml-cpu.a`, `libggml-blas.a`, `libggml-base.a`, `libcommon.a`). These are copied to the project root.

3. **Compile the Go binary** with CGo:
   ```bash
   CGO_CFLAGS="-I$(pwd)" CGO_LDFLAGS="-L$(pwd)" CGO_ENABLED=1 go build -o opuscriber .
   ```
   The `-I` flag finds `whisper.h` and all `ggml*.h` headers. The `-L` flag finds the static libraries.

### CGo Wrapper

`pipeline/whisper.go` wraps the C whisper.cpp API via CGo. The key function `TranscribeAudio` takes 16kHz mono float32 PCM samples and returns timestamped text segments. Go strings are bridged to C via `C.CString` and freed with `C.free` after use.

### GPU Acceleration

GPU backend is auto-detected at whisper.cpp init time:
- **macOS**: Metal (`libggml-metal.a`)
- **Linux**: CUDA (NVIDIA) or CPU BLAS fallback
- **All**: CPU fallback works everywhere

whisper.cpp logs the selected backend at startup:
```
whisper_init_with_params_no_state: use gpu    = 1
ggml_metal_device_init: GPU name: MTL0 (Apple M2 Max)
```

## TUI State Machine

The TUI uses a simple four-state model:

```
IDLE ──Enter──► BROWSING ──Tab──► TRANSCRIBING ──► DONE
  ▲                                        │         │
  └────────────────────────────────────────┴──Enter──┘
```

- **IDLE**: Welcome banner, "Press Enter to select audio files"
- **BROWSING**: File picker (`.opus`, `.ogg`, `.oga` only). Select files with Tab, confirm with Enter
- **TRANSCRIBING**: Spinner + progress bar per file, live checkmarks for completed files
- **DONE**: Summary of results (✓ passed / ✗ failed), options to start over or quit

The file picker uses `bubbles/filepicker` with height 15 and auto-navigation. Files that already have both `.txt` and `.srt` outputs are skipped silently (idempotent).

## Model Catalog

Models are stored as ggml binary files (`.bin`) downloaded from HuggingFace's `ggerganov/whisper.cpp` repository. The catalog at `internal/catalog/models.json` contains:

- `id`: CLI identifier (e.g., `tiny`, `medium`, `large-v3-turbo`)
- `filename`: The file on HuggingFace (`ggml-<model>.bin`)
- `size_bytes`: Exact file size for disk display
- `sha256`: Integrity hash verified on every download
- `download_url`: HuggingFace resolve URL
- `recommended`: Whether this is the default choice

Download verification is mandatory — sha256 mismatch aborts the install. Models are stored in `~/.config/opuscriber/models/` by default, overridable with `--models` or `OPUSCRIBER_MODELS`.

## Configuration

Stored in `~/.config/opuscriber/config.json`:

```json
{
  "version": 1,
  "default_model": "medium"
}
```

Path overridable with `OPUSCRIBER_CONFIG`.