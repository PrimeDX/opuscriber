## 1. Demo Assets

- [x] 1.1 Transcode 30s JFK speech clip to `demo/jfk-moon-speech.opus` (16kHz mono Opus) using ffmpeg from `sample-audio/Jfk_rice_university_we_choose_to_go_to_the_moon.ogg`
- [x] 1.2 Create `demo.tape` vhs script: IDLE → file picker → select demo clip → transcribe → DONE screen → cat output .txt
- [x] 1.3 Run `vhs demo.tape` to generate `demo.gif`; verify it renders correctly and is ≤3 MB
- [x] 1.4 Add `demo` target to Makefile: `vhs demo.tape`

## 2. Architecture Documentation

- [x] 2.1 Create `docs/ARCHITECTURE.md` extracting internals from current README Architecture section: CGo, whisper.cpp build chain, pion/opus, Bubble Tea, Cobra+Fang
- [x] 2.2 Expand ARCHITECTURE.md with project structure diagram and build flow

## 3. README Rewrite

- [x] 3.1 Write new H1 + tagline: "Opuscriber — Offline Speech-to-Text for Opus Audio" with SEO-dense opening paragraph
- [x] 3.2 Add demo GIF display below quickstart
- [x] 3.3 Restructure Quickstart: `make build` → `model install medium` → `opuscriber --input ./in --output ./out`
- [x] 3.4 Add "What is Opuscriber?" section (2-3 sentences, natural keyword density)
- [x] 3.5 Add "Why Opuscriber?" comparison prose (vs. whisper.cpp raw, vs. cloud APIs, vs. GUI apps)
- [x] 3.6 Keep Features section; add "no ffmpeg" and platform callouts
- [x] 3.7 Add Model Guide table with real data from `internal/catalog/models.json` (model ID, disk, RAM, speed, best for)
- [x] 3.8 Keep Usage section (commands + flags); verify flag defaults match current `cmd/root.go`
- [x] 3.9 Rewrite Installation: build-from-source first, Homebrew (coming soon), prebuilt binaries (coming soon)
- [x] 3.10 Replace Architecture section with simplified 3-line pipeline diagram + link to `docs/ARCHITECTURE.md`
- [x] 3.11 Replace Contributing section with link to `docs/CONTRIBUTING.md` and `docs/MAINTAINING.md`
- [x] 3.12 Add Documentation section linking to docs/ARCHITECTURE, docs/CONTRIBUTING, docs/MAINTAINING
- [x] 3.13 Add macOS + Linux platform mention in opening and installation

## 4. GitHub Repo

- [x] 4.1 Create GitHub repository `primedx/opuscriber` via `gh repo create`
- [x] 4.2 Set repo description: "Offline speech-to-text and SRT subtitle generation for Opus audio — powered by Whisper. Single binary, no cloud."
- [x] 4.3 Set topics: `go`, `golang`, `whisper`, `speech-to-text`, `transcription`, `offline`, `opus`, `srt`, `cli`, `tui`, `whisper-cpp`
- [x] 4.4 Push all changes (committed files + demo.gif) to `main`

## 5. Verification

- [x] 5.1 Read final README.md against spec: check all SEO terms present, model table accurate, install section honest
- [x] 5.2 Verify `docs/ARCHITECTURE.md` exists with no placeholder content
- [x] 5.3 Verify demo.gif displays correctly in GitHub README preview
- [x] 5.4 Verify GitHub repo description and topics are set correctly