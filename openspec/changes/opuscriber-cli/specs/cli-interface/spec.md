## Purpose

Provide two operation modes for the CLI: a non-interactive pipeline-safe mode (default) and an interactive Bubble Tea TUI (--interactive flag or tui subcommand) for browsing, selecting, and monitoring file transcription.

## ADDED Requirements

### Requirement: Non-interactive CLI mode
The system SHALL by default run in non-interactive mode: process all matching audio files in the input directory, write outputs to the output directory, and exit. Clean stdout/stderr, safe in Docker without -it.

#### Scenario: Batch processing
- **WHEN** opuscriber runs without --interactive
- **THEN** it SHALL process all matching files and exit with code 0 on success

#### Scenario: No audio files
- **WHEN** the input directory contains no matching audio files
- **THEN** the system SHALL print "no audio files found" and exit with code 0

### Requirement: Interactive TUI mode
The system SHALL provide a Bubble Tea TUI when --interactive or -i flag is passed, or the tui subcommand is used. The TUI SHALL include a file picker, processing queue with progress bars, and results summary.

#### Scenario: TTY required
- **WHEN** --interactive is passed but stdin or stdout is not a TTY
- **THEN** the system SHALL print an error message instructing the user to add -it to docker run and exit with non-zero code

#### Scenario: File selection
- **WHEN** the TUI starts
- **THEN** the user SHALL be able to browse and select audio files from the input directory

#### Scenario: Progress display
- **WHEN** files are being transcribed in the TUI
- **THEN** the system SHALL display per-file progress bars and completion indicators

### Requirement: Fang-wrapped cobra
The command layer SHALL use cobra wrapped with charmbracelet/fang for styled help output, automatic --version, shell completions, and man page generation. Fang SHALL auto-detect terminal color profile and downsample gracefully when piped.

#### Scenario: Help output
- **WHEN** opuscriber --help is invoked
- **THEN** the system SHALL print styled help text with all available flags

#### Scenario: Version flag
- **WHEN** opuscriber --version is invoked
- **THEN** the system SHALL print the version string

### Requirement: CLI flags
The CLI SHALL support the following flags: --lang (default "pt"), --model (default "medium"), --input (default "/audio/in"), --output (default "/audio/out"), --models (default "/models"), --interactive.

#### Scenario: Flag defaults
- **WHEN** no flags are provided
- **THEN** the system SHALL use documented default values for all flags
