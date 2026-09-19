package chordTransposer

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestBehaviorContract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		semitones int
		format    string
		expected  string
	}{
		{
			name:     "BC-01 no final newline",
			input:    "C",
			format:   "%v",
			expected: "C",
		},
		{
			name:      "BC-02 CRLF",
			input:     "C\r\nD\r\n",
			semitones: 2,
			format:    "%v",
			expected:  "D\r\nE\r\n",
		},
		{
			name:     "BC-04 zero preserves spelling",
			input:    "Ab",
			format:   "%v",
			expected: "Ab",
		},
		{
			name:      "BC-05 large positive shift",
			input:     "G",
			semitones: 25,
			format:    "%v",
			expected:  "G#",
		},
		{
			name:      "BC-06 large negative shift",
			input:     "G",
			semitones: -23,
			format:    "%v",
			expected:  "G#",
		},
		{
			name:      "BC-07 enharmonic roots",
			input:     "B# Cb E# Fb",
			semitones: 1,
			format:    "%v",
			expected:  "C# C F# F",
		},
		{
			name:     "BC-08 whole slash chord formatting",
			input:    "C/E",
			format:   "<%v>",
			expected: "<C/E>",
		},
		{
			name:     "BC-09 whole six-nine formatting",
			input:    "C6/9",
			format:   "<%v>",
			expected: "<C6/9>",
		},
		{
			name:      "BC-10 repeated qualities rejected",
			input:     "Cmaj7m9",
			semitones: 2,
			format:    "<%v>",
			expected:  "Cmaj7m9",
		},
		{
			name:      "BC-11 double accidentals rejected",
			input:     "C## Ebb",
			semitones: 2,
			format:    "<%v>",
			expected:  "C## Ebb",
		},
		{
			name:      "BC-12 parenthesized quality",
			input:     "D(sus4) Dsus4",
			semitones: 2,
			format:    "%v",
			expected:  "E(sus4) Esus4",
		},
		{
			name:      "BC-13 issue one rhythm notation",
			input:     "A / / / | F#m / / / | D / / / | E/G# / / /",
			semitones: 2,
			format:    "%v",
			expected:  "B / / / | G#m / / / | E / / / | F#/A# / / /",
		},
		{
			name:      "BC-14 section and lyrics unchanged",
			input:     "[Intro]\nA lyric line\n[C]hello",
			semitones: 2,
			format:    "<%v>",
			expected:  "[Intro]\nA lyric line\n[C]hello",
		},
		{
			name:     "BC-15 invalid legacy format",
			input:    "C",
			format:   "%d",
			expected: "C",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := TransposeChords(test.input, test.semitones, test.format)
			if actual != test.expected {
				t.Fatalf("TransposeChords() = %q, want %q", actual, test.expected)
			}
		})
	}
}

func TestLongLineDoesNotTruncateOrExit(t *testing.T) {
	t.Parallel()

	input := strings.Repeat(" ", 70*1024) + "C"
	actual := TransposeChords(input, 2, "%v")
	expected := strings.Repeat(" ", 70*1024) + "D"
	if actual != expected {
		t.Fatalf("long-line result length = %d, want %d", len(actual), len(expected))
	}
}

func TestIssue1PartiallyNotWorking(t *testing.T) {
	t.Parallel()

	input := "[Intro]\nA / / / | F#m / / / | D / / / | E/G# / / /\nD(sus4)             Dsus4"
	expected := "[Intro]\nB / / / | G#m / / / | E / / / | F#/A# / / /\nE(sus4)             Esus4"
	actual := TransposeChords(input, 2, "%v")
	if actual != expected {
		t.Fatalf("TransposeChords() = %q, want %q", actual, expected)
	}
}

func TestTransposeOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  Options
		expected string
	}{
		{name: "automatic spelling", options: Options{}, expected: "A#"},
		{name: "sharp spelling", options: Options{Spelling: SpellingSharps}, expected: "A#"},
		{name: "flat spelling", options: Options{Spelling: SpellingFlats}, expected: "Bb"},
		{
			name: "typed formatter",
			options: Options{Formatter: func(chord string) string {
				return "[" + chord + "]"
			}},
			expected: "[A#]",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual, err := Transpose("A", 1, test.options)
			if err != nil {
				t.Fatalf("Transpose() error = %v", err)
			}
			if actual != test.expected {
				t.Fatalf("Transpose() = %q, want %q", actual, test.expected)
			}
		})
	}
}

func TestTransposeRejectsInvalidSpelling(t *testing.T) {
	t.Parallel()

	_, err := Transpose("C", 1, Options{Spelling: SpellingPreference(255)})
	if !errors.Is(err, ErrInvalidSpellingPreference) {
		t.Fatalf("Transpose() error = %v, want %v", err, ErrInvalidSpellingPreference)
	}
}

func TestSupportedQualities(t *testing.T) {
	t.Parallel()

	for _, quality := range qualities {
		if quality == "" {
			continue
		}
		quality := quality
		t.Run(quality, func(t *testing.T) {
			t.Parallel()
			for _, token := range []string{"C" + quality, "C(" + quality + ")"} {
				if _, ok := parseChord(token); !ok {
					t.Errorf("parseChord(%q) was rejected", token)
				}
			}
		})
	}
}

func TestChordContextPreservesStructure(t *testing.T) {
	t.Parallel()

	input := "\t|C|  G/B\t/ |\n"
	expected := "\t|D|  A/C#\t/ |\n"
	actual := TransposeChords(input, 2, "%v")
	if actual != expected {
		t.Fatalf("TransposeChords() = %q, want %q", actual, expected)
	}
}

func TestLegacyFormats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format   string
		expected string
	}{
		{format: "[%v]", expected: "[C]"},
		{format: "[%s]", expected: "[C]"},
		{format: "%% %v", expected: "% C"},
		{format: "", expected: "C"},
		{format: "%v %v", expected: "C"},
		{format: "%1000000s", expected: "C"},
	}

	for _, test := range tests {
		if actual := TransposeChords("C", 0, test.format); actual != test.expected {
			t.Errorf("TransposeChords(format %q) = %q, want %q", test.format, actual, test.expected)
		}
	}
}

func TestEmptyAndStructuralLines(t *testing.T) {
	t.Parallel()

	input := "\n| / ||\n"
	if actual := TransposeChords(input, 4, "<%v>"); actual != input {
		t.Fatalf("TransposeChords() = %q, want %q", actual, input)
	}
}

func ExampleTranspose() {
	result, err := Transpose("A | C#m/G#", 1, Options{Spelling: SpellingFlats})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output: Bb | Dm/A
}

func ExampleTransposeChords() {
	result := TransposeChords("C/E", 2, "<span class='chord'>%v</span>")
	fmt.Println(result)
	// Output: <span class='chord'>D/F#</span>
}

func BenchmarkTranspose(b *testing.B) {
	song := strings.Repeat("A / / / | F#m / / / | D / / / | E/G# / / /\nlyrics stay unchanged\n", 500)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Transpose(song, 2, Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
