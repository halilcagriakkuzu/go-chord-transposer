# Go Chord Transposer
Chord transposer for Golang.
It will detect chords inside multiline string(Song guitar chords) then replace the chords given string, if you want, it will transpose chords too.

# Install
`$ go get github.com/halilcagriakkuzu/go-chord-transposer@upgrade`

## Simple Example
```
package main

import (
	"fmt"

	chordTransposer "github.com/halilcagriakkuzu/go-chord-transposer"
)

func main() {
	fmt.Println(chordTransposer.TransposeChords("Em", 0, "%v"))
}
```

## TransposeChords(song string, transposeValue int, format string)
### Param1 : song
```
// You can pass whole song in multi line string, it will detect the chord lines automatically
song := `
[Verse 1]
 
D
She's got a smile that it seems to me
  C
Reminds me of childhood memories
       G
Where everything
                                D
Was as fresh as a bright blue sky
D
Now and then when I see her face
    C
She takes me away to that special place
         G
And if I stared too long
                              D
I'd probably break down and cry

[Chorus]
 
A               C              D
Whoa Oh, Sweet child o' mine
A               C              D
Whoa, Oh, Oh, Oh Sweet love o' mine
`
chordTransposer.TransposeChords(song, 0, "%v")
```

### Param2 : transposeValue
```
// You can give integer value for transpose, between -11,+11
chordTransposer.TransposeChords("Em", 1, "%v")
// output : Fm
chordTransposer.TransposeChords("Em", -2, "%v")
// output : Dm
```

### Param3 : formatString
```
// You can use %v for chord string and put anything you want. If you don't want to format, just use "%v"
chordTransposer.TransposeChords("Em", 0, "%v")
// output : Em
chordTransposer.TransposeChords("Em", 0, "<span class='chord'>%v</span>")
// output : <span class='chord'>Em</span>
```
