package transposer

import "fmt"

func Transpose(chord string, t int) string {
	message := fmt.Sprintf("given cord\n %v\ngiven t:%v", chord, t)
	return message
}
