package main

import "fmt"

type SimplifiedAES struct {
	preRoundKey []int
	round1Key   []int
	round2Key   []int
}

var sBox = [16]int{
	0x9, 0x4, 0xA, 0xB, 0xD, 0x1, 0x8, 0x5, 0x6, 0x2, 0x0, 0x3, 0xC, 0xE, 0xF, 0x7,
}

var sBoxI = [16]int{
	0xA, 0x5, 0x9, 0xB, 0x1, 0x7, 0x8, 0xF, 0x6, 0x0, 0x2, 0x3, 0xC, 0x4, 0xD, 0xE,
}

func NewSimplifiedAES(key int) *SimplifiedAES {
	preRoundKey, round1Key, round2Key := keyExpansion(key)
	return &SimplifiedAES{
		preRoundKey: preRoundKey,
		round1Key:   round1Key,
		round2Key:   round2Key,
	}
}

func subWord(word int) int {
	return (sBox[(word>>4)] << 4) + sBox[word&0x0F]
}

func rotWord(word int) int {
	return ((word & 0x0F) << 4) + ((word & 0xF0) >> 4)
}

func keyExpansion(key int) ([]int, []int, []int) {
	Rcon1 := 0x80
	Rcon2 := 0x30

	w := make([]int, 6)
	w[0] = (key & 0xFF00) >> 8
	w[1] = key & 0x00FF
	w[2] = w[0] ^ (subWord(rotWord(w[1])) ^ Rcon1)
	w[3] = w[2] ^ w[1]
	w[4] = w[2] ^ (subWord(rotWord(w[3])) ^ Rcon2)
	w[5] = w[4] ^ w[3]

	return intToState((w[0] << 8) + w[1]), intToState((w[2] << 8) + w[3]), intToState((w[4] << 8) + w[5])
}

func gfMult(a, b int) int {
	product := 0
	a = a & 0x0F
	b = b & 0x0F

	for a != 0 && b != 0 {
		if b&1 != 0 {
			product = product ^ a
		}
		a = a << 1
		if a&(1<<4) != 0 {
			a = a ^ 0b10011
		}
		b = b >> 1
	}

	return product
}

func intToState(n int) []int {
	return []int{(n >> 12) & 0xF, (n >> 4) & 0xF, (n >> 8) & 0xF, n & 0xF}
}

func stateToInt(m []int) int {
	return (m[0] << 12) + (m[2] << 8) + (m[1] << 4) + m[3]
}

func addRoundKey(s1, s2 []int) []int {
	result := make([]int, len(s1))
	for i := range s1 {
		result[i] = s1[i] ^ s2[i]
	}
	return result
}

func subNibbles(sbox [16]int, state []int) []int {
	result := make([]int, len(state))
	for i := range state {
		result[i] = sbox[state[i]]
	}
	return result
}

func shiftRows(state []int) []int {
	return []int{state[0], state[1], state[3], state[2]}
}

func mixColumns(state []int) []int {
	return []int{
		state[0] ^ gfMult(4, state[2]),
		state[1] ^ gfMult(4, state[3]),
		state[2] ^ gfMult(4, state[0]),
		state[3] ^ gfMult(4, state[1]),
	}
}

func inverseMixColumns(state []int) []int {
	return []int{
		gfMult(9, state[0]) ^ gfMult(2, state[2]),
		gfMult(9, state[1]) ^ gfMult(2, state[3]),
		gfMult(9, state[2]) ^ gfMult(2, state[0]),
		gfMult(9, state[3]) ^ gfMult(2, state[1]),
	}
}

func (saes *SimplifiedAES) Encrypt(plaintext int) int {
	state := addRoundKey(saes.preRoundKey, intToState(plaintext))
	state = mixColumns(shiftRows(subNibbles(sBox, state)))
	state = addRoundKey(saes.round1Key, state)
	state = shiftRows(subNibbles(sBox, state))
	state = addRoundKey(saes.round2Key, state)
	return stateToInt(state)
}

func (saes *SimplifiedAES) Decrypt(ciphertext int) int {
	state := addRoundKey(saes.round2Key, intToState(ciphertext))
	state = subNibbles(sBoxI, shiftRows(state))
	state = inverseMixColumns(addRoundKey(saes.round1Key, state))
	state = subNibbles(sBoxI, shiftRows(state))
	state = addRoundKey(saes.preRoundKey, state)
	return stateToInt(state)
}

func main() {
	key := 0b0100101011110101
	plaintext := 0b1101011100101000
	

	saes := NewSimplifiedAES(key)

	enc := saes.Encrypt(plaintext)
	fmt.Printf("Encrypted: %016b\n", enc)

	dec := saes.Decrypt(enc)
	fmt.Printf("Decrypted: %016b\n", dec)
}
