## Purpose

The README.md is the public face of Opuscriber on GitHub. It must serve search engine discoverability, first-time user comprehension, and quickstart success — without burdening readers with internal architecture or contributing/maintaining instructions that belong in dedicated docs.

## ADDED Requirements

### Requirement: SEO-optimized heading and opening
The README H1 and first two paragraphs MUST contain primary search terms: speech-to-text, offline audio transcription, Whisper, SRT subtitles, Opus audio.

#### Scenario: Search engine indexing
- **WHEN** a search engine crawls the repo page
- **THEN** the H1, repo description, and first `<p>` contain "speech-to-text", "audio transcription", "Whisper", "offline", and "SRT"

### Requirement: Quickstart-first layout
The Quickstart section MUST appear before any explanatory prose, showing the three commands needed to install, download a model, and transcribe in under 60 seconds.

#### Scenario: First-time visitor
- **WHEN** a user scrolls past the demo GIF
- **THEN** they see a 3-line quickstart: `make build`, `model install`, `--input/--output`

### Requirement: Honest pre-release install section
The Installation section MUST lead with the working build-from-source path, then state Homebrew and prebuilt binaries as "coming soon" (not yet published).

#### Scenario: User tries to install
- **WHEN** a user reads the Installation section
- **THEN** `make build` is the first option shown, with prerequisites listed clearly
- **AND** Homebrew is shown second with a "coming soon" qualifier
- **AND** prebuilt binaries are mentioned third with a link to future GitHub Releases

### Requirement: Model guide table
The README MUST include a model comparison table with real data from `internal/catalog/models.json`: model ID, disk size, approximate RAM, relative speed, and recommended use case.

#### Scenario: User chooses a model
- **WHEN** a user reads the Model Guide section
- **THEN** they see all six models listed with accurate disk sizes in MB/GB
- **AND** `medium` is marked as recommended
- **AND** the table uses model IDs exactly as they appear in the catalog: `tiny`, `base`, `small`, `medium`, `large-v3`, `large-v3-turbo`

### Requirement: Comparison section
The README MUST include a "Why Opuscriber?" section comparing Opuscriber to: raw whisper.cpp, cloud APIs (OpenAI), and local GUI apps — in prose, not a table.

#### Scenario: User weighing alternatives
- **WHEN** a user reads "Why Opuscriber?"
- **THEN** they understand the tradeoffs vs. whisper.cpp CLI, vs. cloud APIs, vs. GUI alternatives

### Requirement: Architecture moved out
Architecture internals (CGo, libwhisper.a, build chain, pion/opus) MUST live in `docs/ARCHITECTURE.md`, not README. The README keeps only a simplified 3-step pipeline diagram.

#### Scenario: End-user reads README
- **WHEN** a user scans the README
- **THEN** they see a simple pipeline diagram: "Opus file → decode → Whisper → TXT + SRT"
- **AND** no CGo, libwhisper.a, or build-chain internals appear

### Requirement: Contributing and maintaining links
The README MUST link to `docs/CONTRIBUTING.md` and `docs/MAINTAINING.md` without reproducing their content.

#### Scenario: Contributor looks for build instructions
- **WHEN** a user clicks the contributing link
- **THEN** they land on `docs/CONTRIBUTING.md` with full setup, project structure, and test instructions

### Requirement: Platform support
The README MUST state macOS + Linux support in the opening text and installation section. Metal (macOS), CUDA (Linux), and CPU fallback MUST be mentioned.

#### Scenario: Linux user checks compatibility
- **WHEN** a Linux user reads the opening paragraph
- **THEN** they see "macOS and Linux" mentioned
- **AND** GPU acceleration details are in the Features section