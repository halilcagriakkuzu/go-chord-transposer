package chordTransposer

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Chord is retained for v1 source compatibility.
//
// Deprecated: the type has no usable exported state. A future major version
// will replace it with an explicit parsed-chord model.
type Chord struct{}

// ChordRegex is retained for v1 source compatibility.
//
// Deprecated: it describes the legacy grammar and is not used by the parser.
const ChordRegex = `^[A-G][b\#]?(2|4|5|6|7|9|11|13|6\/9|7\-5|7\-9|7\#5|7\#9|7\+5|7\+9|b5|#5|#9|7b5|7b9|7sus2|7sus4|add2|add4|add9|aug|dim|dim7|m\/maj7|m6|m7|m7b5|m9|m11|m13|maj7|maj9|maj11|maj13|M7|M9|M11|M13|mb5|m|sus|sus2|sus4)*(\/[A-G][b\#]*)*$`

// ChordReplaceRegex is retained for v1 source compatibility.
//
// Deprecated: it describes the legacy grammar and is not used by the parser.
const ChordReplaceRegex = `([A-G][b\#]?(2|4|5|6|7|9|11|13|6\/9|7\-5|7\-9|7\#5|7\#9|7\+5|7\+9|b5|#5|#9|7b5|7b9|7sus2|7sus4|add2|add4|add9|aug|dim|dim7|m\/maj7|m6|m7|m7b5|m9|m11|m13|maj7|maj9|maj11|maj13|M7|M9|M11|M13|mb5|m|sus|sus2|sus4)*)`

// SpellingPreference controls how accidentals are rendered after a chord is
// transposed.
type SpellingPreference uint8

const (
	// SpellingAuto preserves a note's flat preference and otherwise uses sharps.
	SpellingAuto SpellingPreference = iota
	// SpellingSharps renders accidental target notes with sharps.
	SpellingSharps
	// SpellingFlats renders accidental target notes with flats.
	SpellingFlats
)

// ErrInvalidSpellingPreference is returned when Options.Spelling is unknown.
var ErrInvalidSpellingPreference = errors.New("invalid spelling preference")

// Options configures Transpose. Its zero value is valid.
type Options struct {
	// Spelling controls accidental rendering. The zero value is SpellingAuto.
	Spelling SpellingPreference
	// Formatter, when non-nil, is called once for each complete transposed chord.
	Formatter func(chord string) string
}

type note struct {
	text       string
	pitch      int
	accidental byte
}

type parsedChord struct {
	root    note
	quality string
	bass    *note
}

type chordSpan struct {
	start int
	end   int
	chord parsedChord
}

var qualities = []string{
	"m/maj7",
	"7sus2", "7sus4",
	"maj13", "maj11", "maj9", "maj7",
	"M13", "M11", "M9", "M7",
	"m13", "m11", "m9", "m7b5", "m7", "m6",
	"7-5", "7-9", "7#5", "7#9", "7+5", "7+9", "7b5", "7b9",
	"add2", "add4", "add9",
	"dim7", "dim", "aug",
	"sus2", "sus4", "sus",
	"mb5", "b5", "#5", "#9",
	"6/9", "13", "11", "9", "7", "6", "5", "4", "2",
	"m",
	"",
}

var qualitySet = func() map[string]struct{} {
	result := make(map[string]struct{}, len(qualities)-1)
	for _, quality := range qualities {
		if quality != "" {
			result[quality] = struct{}{}
		}
	}
	return result
}()

var naturalPitches = map[byte]int{
	'C': 0,
	'D': 2,
	'E': 4,
	'F': 5,
	'G': 7,
	'A': 9,
	'B': 11,
}

var sharpNames = [...]string{
	"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B",
}

var flatNames = [...]string{
	"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B",
}

// Transpose transposes chords on chord-context lines while preserving all
// other input bytes. The semitone count may be any integer.
func Transpose(song string, semitones int, options Options) (string, error) {
	if options.Spelling > SpellingFlats {
		return "", ErrInvalidSpellingPreference
	}

	return transformSong(song, semitones, options), nil
}

// TransposeChords transposes and formats chords using the legacy v1 API.
// Invalid format strings safely fall back to %v.
func TransposeChords(song string, semitones int, format string) string {
	format = validLegacyFormatOrDefault(format)
	result, _ := Transpose(song, semitones, Options{
		Formatter: func(chord string) string {
			return fmt.Sprintf(format, chord)
		},
	})
	return result
}

func transformSong(song string, semitones int, options Options) string {
	if song == "" {
		return ""
	}

	var result strings.Builder
	result.Grow(len(song))
	start := 0

	for start < len(song) {
		relativeEnd := strings.IndexByte(song[start:], '\n')
		if relativeEnd < 0 {
			result.WriteString(transformLine(song[start:], semitones, options))
			break
		}

		end := start + relativeEnd
		if end > start && song[end-1] == '\r' {
			result.WriteString(transformLine(song[start:end-1], semitones, options))
			result.WriteString("\r\n")
		} else {
			result.WriteString(transformLine(song[start:end], semitones, options))
			result.WriteByte('\n')
		}
		start = end + 1
	}

	return result.String()
}

func transformLine(line string, semitones int, options Options) string {
	spans, ok := findChordSpans(line)
	if !ok || len(spans) == 0 {
		return line
	}

	var result strings.Builder
	result.Grow(len(line))
	last := 0
	for _, span := range spans {
		result.WriteString(line[last:span.start])
		chord := renderChord(span.chord, semitones, options.Spelling)
		if options.Formatter != nil {
			chord = options.Formatter(chord)
		}
		result.WriteString(chord)
		last = span.end
	}
	result.WriteString(line[last:])
	return result.String()
}

func findChordSpans(line string) ([]chordSpan, bool) {
	var spans []chordSpan
	for index := 0; index < len(line); {
		r, size := utf8.DecodeRuneInString(line[index:])
		if unicode.IsSpace(r) {
			index += size
			continue
		}
		if line[index] == '|' {
			index++
			continue
		}

		start := index
		for index < len(line) {
			r, size = utf8.DecodeRuneInString(line[index:])
			if unicode.IsSpace(r) || line[index] == '|' {
				break
			}
			index += size
		}

		token := line[start:index]
		if token == "/" {
			continue
		}

		chord, ok := parseChord(token)
		if !ok {
			return nil, false
		}
		spans = append(spans, chordSpan{start: start, end: index, chord: chord})
	}

	return spans, true
}

func parseChord(token string) (parsedChord, bool) {
	root, consumed, ok := parseNotePrefix(token)
	if !ok {
		return parsedChord{}, false
	}
	remainder := token[consumed:]

	if strings.HasPrefix(remainder, "(") {
		closing := strings.IndexByte(remainder, ')')
		if closing <= 1 {
			return parsedChord{}, false
		}
		quality := remainder[1:closing]
		if _, supported := qualitySet[quality]; !supported {
			return parsedChord{}, false
		}
		bass, ok := parseOptionalBass(remainder[closing+1:])
		if !ok {
			return parsedChord{}, false
		}
		return parsedChord{root: root, quality: remainder[:closing+1], bass: bass}, true
	}

	for _, quality := range qualities {
		if !strings.HasPrefix(remainder, quality) {
			continue
		}
		bass, ok := parseOptionalBass(remainder[len(quality):])
		if ok {
			return parsedChord{root: root, quality: quality, bass: bass}, true
		}
	}

	return parsedChord{}, false
}

func parseOptionalBass(input string) (*note, bool) {
	if input == "" {
		return nil, true
	}
	if len(input) < 2 || input[0] != '/' {
		return nil, false
	}
	bass, consumed, ok := parseNotePrefix(input[1:])
	if !ok || consumed != len(input)-1 {
		return nil, false
	}
	return &bass, true
}

func parseNotePrefix(input string) (note, int, bool) {
	if input == "" {
		return note{}, 0, false
	}
	pitch, ok := naturalPitches[input[0]]
	if !ok {
		return note{}, 0, false
	}

	consumed := 1
	var accidental byte
	if len(input) > 1 && (input[1] == '#' || input[1] == 'b') {
		accidental = input[1]
		consumed++
		if accidental == '#' {
			pitch++
		} else {
			pitch--
		}
	}

	return note{
		text:       input[:consumed],
		pitch:      normalizePitch(pitch),
		accidental: accidental,
	}, consumed, true
}

func renderChord(chord parsedChord, semitones int, spelling SpellingPreference) string {
	var result strings.Builder
	result.WriteString(renderNote(chord.root, semitones, spelling))
	result.WriteString(chord.quality)
	if chord.bass != nil {
		result.WriteByte('/')
		result.WriteString(renderNote(*chord.bass, semitones, spelling))
	}
	return result.String()
}

func renderNote(source note, semitones int, spelling SpellingPreference) string {
	shift := normalizePitch(semitones)
	if shift == 0 {
		return source.text
	}

	target := normalizePitch(source.pitch + shift)
	if spelling == SpellingFlats || (spelling == SpellingAuto && source.accidental == 'b') {
		return flatNames[target]
	}
	return sharpNames[target]
}

func normalizePitch(value int) int {
	return ((value % 12) + 12) % 12
}

func validLegacyFormatOrDefault(format string) string {
	placeholders := 0
	for index := 0; index < len(format); index++ {
		if format[index] != '%' {
			continue
		}
		if index+1 >= len(format) {
			return "%v"
		}
		index++
		switch format[index] {
		case '%':
		case 's', 'v':
			placeholders++
		default:
			return "%v"
		}
	}
	if placeholders != 1 {
		return "%v"
	}
	return format
}
