package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func main() {
	airplane := '\u2708'
	byteAirplane := make([]byte, utf8.RuneLen(airplane))
	utf8.EncodeRune(byteAirplane, airplane)

	buffer := new(bytes.Buffer)
	buffer.Write(byteAirplane)
	buffer.Write(byteAirplane[:2])

	output, size, err := buffer.ReadRune()

	output, size, err = buffer.ReadRune()

	err = buffer.UnreadRune()

	unread := buffer.Bytes()

	buffer.Reset()

	buffer.Write(unread)

	buffer.WriteByte(136)

	output, size, err = buffer.ReadRune()

	fmt.Println(airplane, buffer, output, size, err)
}
