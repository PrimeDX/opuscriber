// Package pipeline provides the transcription pipeline: ffmpeg decode,
// whisper-cli transcription, and text reflow.
//
// SRT output is produced directly by whisper-cli and passed through
// without modification — whisper.cpp generates correct SRT with
// HH:MM:SS,mmm timestamps compatible with DaVinci Resolve.
package pipeline
