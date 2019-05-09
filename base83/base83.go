package base83

import (
	"fmt"
	"strings"
)

const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

func Encode(value, length int, pos int, buffer *[]byte) int {
	divisor := 1
	for i := 0; i < length-1; i++ {
		divisor *= 83
	}

	for i := 0; i < length; i++ {
		digit := (value / divisor) % 83
		divisor /= 83
		(*buffer)[pos] = characters[digit]
		pos++
	}

	return pos
}

func Decode(str string) (value int, err error) {
	for _, r := range str {
		idx := strings.IndexRune(characters, r)
		if idx == -1 {
			return 0, fmt.Errorf("base83: invalid string")
		}
		value = value*83 + idx
	}
	return value, nil
}
