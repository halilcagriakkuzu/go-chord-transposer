# Go Chord Transposer

A small Go library for transposing chord sheets while preserving lyrics,
spacing, measure markers, and line endings.

## Requirements

- Go 1.25 or newer

## Install

```bash
go get github.com/halilcagriakkuzu/go-chord-transposer@latest
```

## Usage

The typed API is recommended for new code:

```go
package main

import (
	"fmt"

	chordTransposer "github.com/halilcagriakkuzu/go-chord-transposer"
)

func main() {
	song := "A / / / | F#m / / / | D / / / | E/G# / / /"

	result, err := chordTransposer.Transpose(song, 2, chordTransposer.Options{})
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// B / / / | G#m / / / | E / / / | F#/A# / / /
}
```

Choose the accidental style with `Options.Spelling`:

```go
result, err := chordTransposer.Transpose(
	"A",
	1,
	chordTransposer.Options{Spelling: chordTransposer.SpellingFlats},
)
// result == "Bb"
```

Use a typed formatter when the complete chord needs markup:

```go
result, err := chordTransposer.Transpose(
	"C/E",
	2,
	chordTransposer.Options{
		Formatter: func(chord string) string {
			return "<span class='chord'>" + chord + "</span>"
		},
	},
)
// result == "<span class='chord'>D/F#</span>"
```

## Legacy API

`TransposeChords` remains available for v1 compatibility:

```go
result := chordTransposer.TransposeChords("Em", 1, "<%v>")
// result == "<Fm>"
```

The legacy format accepts exactly one `%v` or `%s` placeholder. Invalid or
empty formats safely fall back to `%v`.

## Supported notation

- Natural, sharp, and flat roots, including enharmonic spellings such as `B#`,
  `Cb`, `E#`, and `Fb`
- Major, minor, suspended, added-tone, diminished, augmented, and altered
  qualities supported by the original v1 grammar
- Slash chords such as `E/G#`
- `6/9` and `m/maj7`
- Parenthesized qualities such as `D(sus4)`
- Measure separators and beat markers such as `A / / / | F#m / / /`

Double accidentals and inline lyric notation such as `[C]hello` are not
supported. Unrecognized lines are returned unchanged.

The complete target and compatibility rules are in
[`docs/behavior-contract.md`](docs/behavior-contract.md).

## Development

```bash
go test -race -cover ./...
go vet ./...
go test -run=^$ -fuzz=FuzzTransposePreservesInputAtZero -fuzztime=10s
go test -run=^$ -fuzz=FuzzTransposePreservesLineStructure -fuzztime=10s
```

CI also runs `staticcheck` and `govulncheck`.
