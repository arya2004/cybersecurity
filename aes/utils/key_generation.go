package utils

import "fmt"

func KeyGeneration(plaintext uint16, key uint16) (uint16, uint16, uint16) {
	var key_0 uint16
	var key_1 uint16
	var key_2 uint16

	word_0 := uint8(key >> 8)
	word_1 := uint8(key & 0xFF)

	fmt.Println(word_0, word_1)

	return key_0, key_1, key_2
}
