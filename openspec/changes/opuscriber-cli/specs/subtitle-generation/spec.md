## Purpose

Generate standard SRT (SubRip) subtitle files from whisper-cli transcription output, compatible with DaVinci Resolve's subtitle import.

## ADDED Requirements

### Requirement: SRT output format
The system SHALL produce SRT files with standard formatting: 1-indexed sequence numbers, HH:MM:SS,mmm timestamps with comma as millisecond separator, and blank line between entries.

#### Scenario: Correct SRT structure
- **WHEN** a file is transcribed
- **THEN** the .srt output SHALL be a valid SRT file importable by DaVinci Resolve via File > Import > Subtitles

### Requirement: SRT from whisper-cli
The system SHALL use whisper-cli's --output-srt flag to generate SRT files directly. No post-processing SHALL modify the SRT content.

#### Scenario: Direct pipeline
- **WHEN** whisper-cli runs with -osrt
- **THEN** the system SHALL copy the resulting .srt file to the output directory without modification

### Requirement: File naming
SRT output files SHALL use the same base name as the input audio file, with .srt extension, in the output directory.

#### Scenario: Output path
- **WHEN** processing /audio/in/recording.opus
- **THEN** the SRT SHALL be written to /audio/out/recording.srt
