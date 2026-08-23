## Purpose

Launch the public GitHub repository with correct metadata (description, topics), SEO-friendly README, a terminal demo recording, and sample audio for the demo — so the project is discoverable and credible from day one.

## ADDED Requirements

### Requirement: GitHub repo description
The repository description (set in GitHub settings, displayed under the repo name and in search) MUST be: "Offline speech-to-text and SRT subtitle generation for Opus audio — powered by Whisper. Single binary, no cloud."

#### Scenario: GitHub search results
- **WHEN** a user searches GitHub for "whisper speech to text opus"
- **THEN** the Opuscriber repo appears with the description containing "speech-to-text", "Offline", "Whisper", "Opus audio"

### Requirement: GitHub topics
The repository MUST have topics set: `go`, `golang`, `whisper`, `speech-to-text`, `transcription`, `offline`, `opus`, `srt`, `cli`, `tui`, `whisper-cpp`.

#### Scenario: Topic-based discovery
- **WHEN** a user browses github.com/topics/whisper or github.com/topics/speech-to-text
- **THEN** Opuscriber is listed among the repos

### Requirement: Terminal demo GIF
The README MUST display a terminal recording (`demo.gif`) showing the TUI transcribing the JFK moon speech sample audio.

#### Scenario: Visitor sees the demo
- **WHEN** a visitor opens the README
- **THEN** the demo GIF is visible near the top, below the quickstart
- **AND** the GIF shows: IDLE screen → file picker → transcribing progress → DONE with checkmarks

### Requirement: Demo audio file
The repository MUST include a 30-second Opus clip (`demo/jfk-moon-speech.opus`) transcribed from the public domain JFK moon speech source file at `sample-audio/`. The original Wikimedia `.ogg` source remains in `sample-audio/` for provenance.

#### Scenario: Demo is self-contained
- **WHEN** a maintainer runs `make demo`
- **THEN** the demo GIF is generated using only files in the repo (no network, no external deps beyond `vhs`)

### Requirement: vhs tape script
The repository MUST include `demo.tape` — a vhs (charmbracelet/vhs) script that reproduces the terminal demo recording.

#### Scenario: Demo is reproducible
- **WHEN** a maintainer installs vhs and runs `vhs demo.tape`
- **THEN** an identical `demo.gif` is produced

### Requirement: make demo target
The Makefile MUST include a `demo` target that runs `vhs demo.tape` to regenerate the demo GIF.

#### Scenario: Regenerating the demo
- **WHEN** a maintainer runs `make demo`
- **THEN** `demo.gif` is rebuilt from `demo.tape`