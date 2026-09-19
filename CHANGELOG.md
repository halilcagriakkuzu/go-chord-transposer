# Changelog

## Unreleased

### Added

- A typed `Transpose` API with accidental-style and formatter options.
- Support for measure separators, standalone beat markers, parenthesized
  qualities, and all single-accidental note spellings.
- Regression, property, fuzz, example, and benchmark tests.
- Automated test, static-analysis, and vulnerability-check workflows.

### Changed

- Transposition now supports any integer semitone value.
- Formatting applies to a complete chord, including its slash bass.
- Input whitespace, line endings, and final-newline state are preserved.
- The minimum supported Go version is now 1.25.

### Fixed

- Removed process termination on long input lines.
- Fixed silent corruption of `B#`, `Cb`, `E#`, and `Fb`.
- Fixed partial formatting of slash and `6/9` chords.
- Fixed rejection of chord lines containing `|` and standalone `/` markers.

### Deprecated

- `Chord`, `ChordRegex`, and `ChordReplaceRegex`; they remain available during
  v1 for source compatibility.
