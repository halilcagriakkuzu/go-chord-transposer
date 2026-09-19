package chordTransposer

import (
	"strings"
	"testing"
)

func FuzzTransposePreservesInputAtZero(f *testing.F) {
	f.Add("")
	f.Add("C/E")
	f.Add("A / / / | F#m / / / | D / / / | E/G# / / /")
	f.Add("[Verse]\r\nLyrics\r\nC G\r\n")
	f.Add("\xff\xfeC")

	f.Fuzz(func(t *testing.T, input string) {
		actual, err := Transpose(input, 0, Options{})
		if err != nil {
			t.Fatalf("Transpose() error = %v", err)
		}
		if actual != input {
			t.Fatalf("zero transposition changed input: got %q, want %q", actual, input)
		}
	})
}

func FuzzTransposePreservesLineStructure(f *testing.F) {
	f.Add("C\nD\r\nE", 1)
	f.Add("A / / / | F#m", -25)
	f.Add("lyrics", 12)

	f.Fuzz(func(t *testing.T, input string, semitones int) {
		actual, err := Transpose(input, semitones, Options{})
		if err != nil {
			t.Fatalf("Transpose() error = %v", err)
		}
		if strings.Count(actual, "\n") != strings.Count(input, "\n") {
			t.Fatalf("newline count changed: got %q from %q", actual, input)
		}
		if strings.Count(actual, "\r") != strings.Count(input, "\r") {
			t.Fatalf("carriage-return count changed: got %q from %q", actual, input)
		}
	})
}
