# Target Behavior Contract

Status: implemented by the unreleased refactor

Date: 2026-09-19

This document defines the behavior provided by the refactored library. It is not
a description of every behavior in the v1.0.1 release. Each rule is backed by an
automated test in the unreleased version.

## 1. Input preservation

1. The result must preserve every byte that is not part of a recognized chord.
2. Spaces, tabs, measure separators, section labels, and other non-chord text
   must remain unchanged.
3. Line endings must be preserved. LF must remain LF, CRLF must remain CRLF,
   and a final newline must not be added or removed.
4. An unrecognized token must remain unchanged and must not be partially
   interpreted as a chord.
5. The library must not log, terminate the process, or panic because of user
   input. In particular, long lines must not be constrained by
   `bufio.Scanner`'s default token limit.

## 2. Chord-line classification

The library transposes chords only on a chord-context line. This avoids treating
ordinary lyrics as chords.

A line is a chord-context line when, after whitespace is ignored, all tokens are
one of the following:

- a recognized chord;
- a measure separator (`|`);
- a standalone beat/hold marker (`/`).

Section labels such as `[Intro]` and lyric lines are not chord-context lines and
must remain unchanged. Inline chord notation embedded in lyrics, such as
`[C]hello`, is outside the scope of this refactor.

The slash has two distinct meanings:

- without surrounding whitespace in `E/G#`, it introduces a bass note;
- as a standalone token in `A / / /`, it is a beat/hold marker and is preserved.

## 3. Supported chord grammar

A recognized chord consists of:

1. a root letter from `A` through `G`;
2. an optional single accidental, `#` or `b`;
3. an optional supported quality;
4. an optional slash followed by a bass note using the same note grammar.

All single-accidental spellings, including `B#`, `Cb`, `E#`, and `Fb`, must map
to the correct pitch class. Double accidentals such as `C##` and `Ebb` are not
supported and must remain unchanged.

The initial quality set remains compatible with v1:

```text
2 4 5 6 7 9 11 13 6/9
7-5 7-9 7#5 7#9 7+5 7+9
b5 #5 #9 7b5 7b9
7sus2 7sus4 add2 add4 add9
aug dim dim7
m/maj7 m6 m7 m7b5 m9 m11 m13
maj7 maj9 maj11 maj13
M7 M9 M11 M13
mb5 m sus sus2 sus4
```

A quality may occur at most once. Arbitrary concatenations such as `Cmaj7m9`
are invalid and must remain unchanged.

Parenthesized quality notation is also supported when its contents are a
supported quality. For example, `D(sus4)` and `D(#5)` are recognized. The
parentheses are preserved in the output.

## 4. Transposition

1. The semitone argument may be any integer. It is normalized modulo 12.
2. A zero-semitone transposition must preserve the original chord spelling.
3. Roots and slash bass notes are transposed independently.
4. Qualities and their parenthesized/non-parenthesized style are preserved
   byte-for-byte.
5. In the default automatic spelling mode, a note originally written with `b`
   prefers flats after transposition. Notes written with `#` or without an
   accidental prefer sharps. Natural target notes never receive an accidental.
6. The options API may explicitly request sharp or flat output. The legacy
   wrapper uses automatic spelling mode.

Examples:

```text
transpose("Ab", 0)       = "Ab"
transpose("Cb", 1)       = "C"
transpose("C/E", 2)      = "D/F#"
transpose("D(sus4)", 2)  = "E(sus4)"
transpose("G", 25)       = "G#"
```

## 5. Formatting

Formatting applies once to the complete chord token. It must not format the
root, quality, or slash bass separately.

```text
format("C/E", "<%v>")   = "<C/E>"
format("C6/9", "<%v>")  = "<C6/9>"
```

The new API exposes a typed formatting option rather than an unrestricted
`fmt.Sprintf` format string. The v1 `TransposeChords` wrapper remains available
and accepts one `%v` or `%s` placeholder. An empty or invalid legacy format uses
`%v` instead of emitting `fmt` diagnostics or terminating the process.

## 6. API compatibility and errors

1. `TransposeChords(song string, semitones int, format string) string` remains
   available throughout v1 as a compatibility wrapper.
2. A new additive API may return an error for invalid options, but ordinary
   unrecognized text is not an error and is preserved.
3. Existing exported identifiers are not removed during v1. They may be
   deprecated and removed only in a new major version.
4. A package-name change is a major-version change and is outside the v1
   refactor.

## 7. Required acceptance scenarios

The automated test suite encodes at least these scenarios:

| ID | Scenario | Required result |
| --- | --- | --- |
| BC-01 | Input without a final newline | No newline is added |
| BC-02 | CRLF input | CRLF is preserved |
| BC-03 | Input line longer than 64 KiB | No exit, panic, or truncation |
| BC-04 | `Ab` transposed by zero | `Ab` |
| BC-05 | `G` transposed by `25` | `G#` |
| BC-06 | `G` transposed by `-23` | `G#` |
| BC-07 | `B# Cb E# Fb` | Every note maps to its correct pitch class |
| BC-08 | `C/E` formatted with `<%v>` | `<C/E>` |
| BC-09 | `C6/9` formatted with `<%v>` | `<C6/9>` |
| BC-10 | `Cmaj7m9` | Left unchanged |
| BC-11 | `C## Ebb` | Left unchanged |
| BC-12 | `D(sus4) Dsus4` transposed by two | `E(sus4) Esus4` |
| BC-13 | `A / / / \| F#m / / / \| D / / / \| E/G# / / /` transposed by two | `B / / / \| G#m / / / \| E / / / \| F#/A# / / /` |
| BC-14 | `[Intro]` and lyric text | Left unchanged |
| BC-15 | Invalid legacy format | Safe `%v` fallback |
