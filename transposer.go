package transposer

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const ChordRegex = `^[A-G][b\#]?(2|4|5|6|7|9|11|13|6\/9|7\-5|7\-9|7\#5|7\#9|7\+5|7\+9|b5|#5|#9|7b5|7b9|7sus2|7sus4|add2|add4|add9|aug|dim|dim7|m\/maj7|m6|m7|m7b5|m9|m11|m13|M7|M9|M11|M13|mb5|m|sus|sus2|sus4)*(\/[A-G][b\#]*)*$`
const ChordReplaceRegex = `([A-G][b\#]?(2|4|5|6|7|9|11|13|6\/9|7\-5|7\-9|7\#5|7\#9|7\+5|7\+9|b5|#5|#9|7b5|7b9|7sus2|7sus4|add2|add4|add9|aug|dim|dim7|m\/maj7|m6|m7|m7b5|m9|m11|m13|maj7|maj9|maj11|maj13|M7|M9|M11|M13|mb5|m|sus|sus2|sus4)*)`

func Transpose(song string, t int) string {
	scanner := bufio.NewScanner(strings.NewReader(song))
	index := 0
	for scanner.Scan() {
		fmt.Println(fmt.Sprintf("[%v] %v - IsChordLine:%v", index, scanner.Text(), isChordLine(scanner.Text())))
		index++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ""
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
