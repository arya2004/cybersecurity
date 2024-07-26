package main

import (
	"fmt"
)

var substitutionMap = map[uint8]uint8{
	0b0000: 0b1001,
	0b0001: 0b0100,
	0b0010: 0b1010,
	0b0011: 0b1011,
	0b0100: 0b1101,
	0b0101: 0b0001,
	0b0110: 0b1000,
	0b0111: 0b0101,
	0b1000: 0b0110,
	0b1001: 0b0010,
	0b1010: 0b0000,
	0b1011: 0b0011,
	0b1100: 0b1100,
	0b1101: 0b1110,
	0b1110: 0b1111,
	0b1111: 0b0111,
}

func SwapNibble(byteVal uint8) uint8 {
	return (byteVal << 4) | (byteVal >> 4)
}

func SubNibble(byteVal uint8) uint8 {
	highNibble := byteVal >> 4
	lowNibble := byteVal & 0x0F
	subHigh := substitutionMap[highNibble]
	subLow := substitutionMap[lowNibble]
	return (subHigh << 4) | subLow
}

func KeyGeneration(key uint16) (uint16, uint16, uint16) {
	word_0 := uint8(key >> 8)
	word_1 := uint8(key & 0xFF)
	fmt.Printf("Initial words: w0 = %08b, w1 = %08b\n", word_0, word_1)

	word_2 := word_0 ^ 0b10000000 ^ SubNibble(SwapNibble(word_1))
	word_3 := word_2 ^ word_1
	fmt.Printf("Intermediate words after first round: w2 = %08b, w3 = %08b\n", word_2, word_3)

	word_4 := word_2 ^ 0b00110000 ^ SubNibble(SwapNibble(word_3))
	word_5 := word_4 ^ word_3
	fmt.Printf("Intermediate words after second round: w4 = %08b, w5 = %08b\n", word_4, word_5)

	key_0 := (uint16(word_0) << 8) | uint16(word_1)
	key_1 := (uint16(word_2) << 8) | uint16(word_3)
	key_2 := (uint16(word_4) << 8) | uint16(word_5)

	return key_0, key_1, key_2
}

func AES(plaintext uint16, key uint16) (cypherText uint16){
	k0, k1, k2 := KeyGeneration(key)
	fmt.Printf("Key0 = %04b %04b %04b %04b\n", k0>>12, (k0>>8)&0xF, (k0>>4)&0xF, k0&0xF)
	fmt.Printf("Key1 = %04b %04b %04b %04b\n", k1>>12, (k1>>8)&0xF, (k1>>4)&0xF, k1&0xF)
	fmt.Printf("Key2 = %04b %04b %04b %04b\n", k2>>12, (k2>>8)&0xF, (k2>>4)&0xF, k2&0xF)
	
	roundKey := plaintext ^ k0
	fmt.Print(roundKey)
	return roundKey

}


func main() {
	AES(0b1101011100101000, 0b0100101011110101)
}
