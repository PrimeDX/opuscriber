package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ProcessFile transcribes a single audio file through the pipeline:
// ffmpeg decode → whisper-cli → txt reflow.
func ProcessFile(audioPath, modelPath, lang, outputDir string) error {
	base := filepath.Base(audioPath)
	ext := filepath.Ext(audioPath)
	baseName := base[:len(base)-len(ext)]

	// Create temp WAV for ffmpeg output
	tmpWAV, err := os.CreateTemp("", "opuscriber-*.wav")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpWAV.Close()
	defer os.Remove(tmpWAV.Name())

	// Step 1: ffmpeg decode to 16kHz mono WAV
	if err := decodeToWAV(audioPath, tmpWAV.Name()); err != nil {
		return fmt.Errorf("ffmpeg decode: %w", err)
	}

	// Step 2: whisper-cli transcription
	outBase := filepath.Join(outputDir, baseName)
	if err := runWhisper(tmpWAV.Name(), modelPath, lang, outBase); err != nil {
		return fmt.Errorf("whisper transcription: %w", err)
	}

	// Step 3: reflow the raw txt output
	rawTxtPath := outBase + ".txt"
	if err := reflowTxtFile(rawTxtPath); err != nil {
		return fmt.Errorf("txt reflow: %w", err)
	}

	return nil
}

func decodeToWAV(input, output string) error {
	cmd := exec.Command("ffmpeg",
		"-i", input,
		"-ar", "16000",
		"-ac", "1",
		"-f", "wav",
		"-y",
		output,
	)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}
	return nil
}

func runWhisper(wavPath, modelPath, lang, outBase string) error {
	cmd := exec.Command("whisper-cli",
		"-m", modelPath,
		"-l", lang,
		"-f", wavPath,
		"-of", outBase,
		"-osrt",
		"-otxt",
		"-np",
	)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("whisper-cli failed: %w", err)
	}
	return nil
}

func reflowTxtFile(txtPath string) error {
	data, err := os.ReadFile(txtPath)
	if err != nil {
		return err
	}

	reflowed := ReflowTxt(string(data))
	if reflowed == "" {
		reflowed = "(empty recording)"
	}

	if err := os.WriteFile(txtPath, []byte(reflowed), 0644); err != nil {
		return err
	}
	return nil
}
