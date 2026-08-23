package pipeline

import (
	"testing"
)

func TestReflowTxt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
		{
			name:  "whitespace only",
			input: "  \n  \n",
			want:  "",
		},
		{
			name:  "single line",
			input: "Hello world",
			want:  "Hello world",
		},
		{
			name:  "multiple lines joined",
			input: "First segment\nsecond part\nthird part",
			want:  "First segment second part third part",
		},
		{
			name:  "paragraph break preserved",
			input: "First paragraph first line\nFirst paragraph second line\n\nSecond paragraph only",
			want:  "First paragraph first line First paragraph second line\n\nSecond paragraph only",
		},
		{
			name:  "multiple paragraph breaks",
			input: "Para one\n\nPara two\n\nPara three",
			want:  "Para one\n\nPara two\n\nPara three",
		},
		{
			name:  "leading and trailing whitespace per line",
			input: "  Hello world  \n  Second segment  ",
			want:  "Hello world Second segment",
		},
		{
			name:  "single line with trailing newline",
			input: "Hello world\n",
			want:  "Hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReflowTxt(tt.input)
			if got != tt.want {
				t.Errorf("ReflowTxt() = %q, want %q", got, tt.want)
			}
		})
	}
}
