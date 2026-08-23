## Purpose

Transcribe Opus-encoded voice notes into timestamped text segments using an offline, locally-running whisper.cpp model with configurable language and model size.

## ADDED Requirements

### Requirement: Supported input formats
The system SHALL accept audio files with extensions .opus, .ogg, and .oga. All three use the same OGG-container/Opus-codec format.

#### Scenario: Accept .opus files
- **WHEN** the input directory contains a file ending in .opus
- **THEN** the system SHALL process it as a valid audio input

#### Scenario: Accept .ogg and .oga files
- **WHEN** the input directory contains files ending in .ogg or .oga
- **THEN** the system SHALL process them as valid audio inputs

### Requirement: FFmpeg decoding
The system SHALL decode input audio to 16 kHz mono WAV format before passing it to the transcription engine.

#### Scenario: Successful decode
- **WHEN** whisper-cli receives the decoded WAV file
- **THEN** it SHALL be 16-bit, 16 kHz, single-channel PCM

### Requirement: Configurable language
The system SHALL accept a language code (--lang flag) passed to whisper-cli for transcription. Default language SHALL be "pt" (Portuguese).

#### Scenario: Default language
- **WHEN** no --lang flag is provided
- **THEN** the system SHALL use "pt" as the language

#### Scenario: Override language
- **WHEN** --lang en is provided
- **THEN** the system SHALL pass "en" to whisper-cli

### Requirement: Configurable model
The system SHALL accept a model name (--model flag) selecting which whisper.cpp ggml model to use. Default SHALL be "medium", resolving to /models/ggml-medium.bin.

#### Scenario: Default model
- **WHEN** no --model flag is provided
- **THEN** the system SHALL use /models/ggml-medium.bin

### Requirement: Idempotent processing
The system SHALL skip files that already have both .txt and .srt outputs in the output directory.

#### Scenario: Already processed
- **WHEN** a file's .txt and .srt outputs both exist in the output directory
- **THEN** the system SHALL skip that file with an informational message

#### Scenario: Partial output
- **WHEN** only one of .txt or .srt exists for a file
- **THEN** the system SHALL re-process that file
