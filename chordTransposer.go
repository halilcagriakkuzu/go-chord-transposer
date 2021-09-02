package chordTransposer

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Chord struct {
	name  string
	kind  string
	value int
}

const ChordRegex = `^[A-G][b\#]?(2|4|5|6|7|9|11|13|6\/9|7\-5|7\-9|7\#5|7\#9|7\+5|7\+9|b5|#5|#9|7b5|7b9|7sus2|7sus4|add2|add4|add9|aug|dim|dim7|m\/maj7|m6|m7|m7b5|m9|m11|m13|maj7|maj9|maj11|maj13|M7|M9|M11|M13|mb5|m|sus|sus2|sus4)*(\/[A-G][b\#]*)*$`
const ChordReplaceRegex = `([A-G][b\#]?(2|4|5|6|7|9|11|13|6\/9|7\-5|7\-9|7\#5|7\#9|7\+5|7\+9|b5|#5|#9|7b5|7b9|7sus2|7sus4|add2|add4|add9|aug|dim|dim7|m\/maj7|m6|m7|m7b5|m9|m11|m13|maj7|maj9|maj11|maj13|M7|M9|M11|M13|mb5|m|sus|sus2|sus4)*)`

func getChords() []Chord {
	return []Chord{
		{name: "Ab", value: 0, kind: "F"},
		{name: "G#", value: 0, kind: "S"},
		{name: "A", value: 1, kind: "N"},
		{name: "Bb", value: 2, kind: "F"},
		{name: "A#", value: 2, kind: "S"},
		{name: "B", value: 3, kind: "N"},
		{name: "C", value: 4, kind: "N"},
		{name: "Db", value: 5, kind: "F"},
		{name: "C#", value: 5, kind: "S"},
		{name: "D", value: 6, kind: "N"},
		{name: "Eb", value: 7, kind: "F"},
		{name: "D#", value: 7, kind: "S"},
		{name: "E", value: 8, kind: "N"},
		{name: "F", value: 9, kind: "N"},
		{name: "Gb", value: 10, kind: "F"},
		{name: "F#", value: 10, kind: "S"},
		{name: "G", value: 11, kind: "N"},
	}
}

func TransposeChords(song string, t int, format string) string {
	scanner := bufio.NewScanner(strings.NewReader(song))
	transposedSong := ""
	line := ""
	for scanner.Scan() {
		line = scanner.Text()
		if isChordLine(line) {
			line = formatChordsInLine(transposeChordLine(line, t), format)
		}
		transposedSong += line + "\n"
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return transposedSong
}

func isChordLine(line string) bool {
	// Get every words in line
	tokens := strings.Fields(line)

	// If there is no word then not chord line
	if len(tokens) == 0 {
		return false
	}

	r, _ := regexp.Compile(ChordRegex)

	// Look for the words one by one , if any of words not match the chord then not chord line
	for _, token := range tokens {
		if !r.MatchString(token) {
			return false
		}
	}

	// Else its chord line
	return true
}

func transposeChordLine(line string, t int) string {
	r, _ := regexp.Compile(ChordReplaceRegex)
	newLine := r.ReplaceAllFunc([]byte(line), func(matched []byte) []byte {
		chordStr := r.FindStringSubmatch(string(matched))[1]
		chordRoot, chordExtra := getChordRoot(chordStr)
		chordRootValue := getChordValueByRoot(chordRoot)
		chordStr = getChordNameByValue(chordRootValue + t)

		return []byte(chordStr + chordExtra)
	})
	return string(newLine)
}

func getChordValueByRoot(root string) int {
	result := -1
	for _, chord := range getChords() {
		if chord.name == root {
			result = chord.value
		}
	}
	return result
}

func getChordNameByValue(value int) string {
	result := ""

	if value > 11 {
		value -= 12
	} else if value < 0 {
		value += 12
	}

	for _, chord := range getChords() {
		if chord.value == value {
			result = chord.name
		}
	}
	return result
}

func getChordRoot(chord string) (string, string) {
	if len(chord) > 1 && (chord[1:2] == "b" || chord[1:2] == "#") {
		return chord[0:2], chord[2:len(chord)]
	}
	return chord[0:1], chord[1:len(chord)]
}

func formatChordsInLine(line string, format string) string {
	r, _ := regexp.Compile(ChordReplaceRegex)
	newLine := r.ReplaceAllFunc([]byte(line), func(matched []byte) []byte {
		chord := r.FindStringSubmatch(string(matched))[1]
		return []byte(fmt.Sprintf(format, chord))
	})
	return string(newLine)
}
