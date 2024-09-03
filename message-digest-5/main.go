package main

import (
	"fmt"
	"math"
)

const (
	S11 = 7
	S12 = 12
	S13 = 17
	S14 = 22
	S21 = 5
	S22 = 9
	S23 = 14
	S24 = 20
	S31 = 4
	S32 = 11
	S33 = 16
	S34 = 23
	S41 = 6
	S42 = 10
	S43 = 15
	S44 = 21
)

var T [64]uint32

func init() {
	for i := 0; i < 64; i++ {
		T[i] = uint32(math.Abs(math.Sin(float64(i+1))) * math.Pow(2, 32))
	}
}

func F(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

func G(x, y, z uint32) uint32 {
	return (x & z) | (y & ^z)
}

func H(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func I(x, y, z uint32) uint32 {
	return y ^ (x | ^z)
}

func rotateLeft(x, n uint32) uint32 {
	return (x << n) | (x >> (32 - n))
}

func encode(input []byte) []uint32 {
	length := ((len(input) + 8) / 64 * 64) + 64
	output := make([]uint32, length/4)

	for i := 0; i < len(input); i++ {
		output[i>>2] |= uint32(input[i]) << ((i % 4) * 8)
	}

	output[len(input)>>2] |= 0x80 << ((len(input) % 4) * 8)

	bitLength := uint64(len(input)) * 8
	output[len(output)-2] = uint32(bitLength)
	output[len(output)-1] = uint32(bitLength >> 32)

	return output
}

func decode(input []uint32) []byte {
	output := make([]byte, len(input)*4)

	for i := 0; i < len(input); i++ {
		output[i*4] = byte(input[i] & 0xff)
		output[i*4+1] = byte((input[i] >> 8) & 0xff)
		output[i*4+2] = byte((input[i] >> 16) & 0xff)
		output[i*4+3] = byte((input[i] >> 24) & 0xff)
	}

	return output
}

func md5Transform(input []byte) []byte {
	x := encode(input)

	a := uint32(0x67452301)
	b := uint32(0xefcdab89)
	c := uint32(0x98badcfe)
	d := uint32(0x10325476)

	for i := 0; i < len(x); i += 16 {
		aa, bb, cc, dd := a, b, c, d

		// Round 1
		a = b + rotateLeft(a+F(b, c, d)+x[i+0]+T[0], S11)
		d = a + rotateLeft(d+F(a, b, c)+x[i+1]+T[1], S12)
		c = d + rotateLeft(c+F(d, a, b)+x[i+2]+T[2], S13)
		b = c + rotateLeft(b+F(c, d, a)+x[i+3]+T[3], S14)
		a = b + rotateLeft(a+F(b, c, d)+x[i+4]+T[4], S11)
		d = a + rotateLeft(d+F(a, b, c)+x[i+5]+T[5], S12)
		c = d + rotateLeft(c+F(d, a, b)+x[i+6]+T[6], S13)
		b = c + rotateLeft(b+F(c, d, a)+x[i+7]+T[7], S14)
		a = b + rotateLeft(a+F(b, c, d)+x[i+8]+T[8], S11)
		d = a + rotateLeft(d+F(a, b, c)+x[i+9]+T[9], S12)
		c = d + rotateLeft(c+F(d, a, b)+x[i+10]+T[10], S13)
		b = c + rotateLeft(b+F(c, d, a)+x[i+11]+T[11], S14)
		a = b + rotateLeft(a+F(b, c, d)+x[i+12]+T[12], S11)
		d = a + rotateLeft(d+F(a, b, c)+x[i+13]+T[13], S12)
		c = d + rotateLeft(c+F(d, a, b)+x[i+14]+T[14], S13)
		b = c + rotateLeft(b+F(c, d, a)+x[i+15]+T[15], S14)

		// Round 2
		a = b + rotateLeft(a+G(b, c, d)+x[i+1]+T[16], S21)
		d = a + rotateLeft(d+G(a, b, c)+x[i+6]+T[17], S22)
		c = d + rotateLeft(c+G(d, a, b)+x[i+11]+T[18], S23)
		b = c + rotateLeft(b+G(c, d, a)+x[i+0]+T[19], S24)
		a = b + rotateLeft(a+G(b, c, d)+x[i+5]+T[20], S21)
		d = a + rotateLeft(d+G(a, b, c)+x[i+10]+T[21], S22)
		c = d + rotateLeft(c+G(d, a, b)+x[i+15]+T[22], S23)
		b = c + rotateLeft(b+G(c, d, a)+x[i+4]+T[23], S24)
		a = b + rotateLeft(a+G(b, c, d)+x[i+9]+T[24], S21)
		d = a + rotateLeft(d+G(a, b, c)+x[i+14]+T[25], S22)
		c = d + rotateLeft(c+G(d, a, b)+x[i+3]+T[26], S23)
		b = c + rotateLeft(b+G(c, d, a)+x[i+8]+T[27], S24)
		a = b + rotateLeft(a+G(b, c, d)+x[i+13]+T[28], S21)
		d = a + rotateLeft(d+G(a, b, c)+x[i+2]+T[29], S22)
		c = d + rotateLeft(c+G(d, a, b)+x[i+7]+T[30], S23)
		b = c + rotateLeft(b+G(c, d, a)+x[i+12]+T[31], S24)

		// Round 3
		a = b + rotateLeft(a+H(b, c, d)+x[i+5]+T[32], S31)
		d = a + rotateLeft(d+H(a, b, c)+x[i+8]+T[33], S32)
		c = d + rotateLeft(c+H(d, a, b)+x[i+11]+T[34], S33)
		b = c + rotateLeft(b+H(c, d, a)+x[i+14]+T[35], S34)
		a = b + rotateLeft(a+H(b, c, d)+x[i+1]+T[36], S31)
		d = a + rotateLeft(d+H(a, b, c)+x[i+4]+T[37], S32)
		c = d + rotateLeft(c+H(d, a, b)+x[i+7]+T[38], S33)
		b = c + rotateLeft(b+H(c, d, a)+x[i+10]+T[39], S34)
		a = b + rotateLeft(a+H(b, c, d)+x[i+13]+T[40], S31)
		d = a + rotateLeft(d+H(a, b, c)+x[i+0]+T[41], S32)
		c = d + rotateLeft(c+H(d, a, b)+x[i+3]+T[42], S33)
		b = c + rotateLeft(b+H(c, d, a)+x[i+6]+T[43], S34)
		a = b + rotateLeft(a+H(b, c, d)+x[i+9]+T[44], S31)
		d = a + rotateLeft(d+H(a, b, c)+x[i+12]+T[45], S32)
		c = d + rotateLeft(c+H(d, a, b)+x[i+15]+T[46], S33)
		b = c + rotateLeft(b+H(c, d, a)+x[i+2]+T[47], S34)

		 // Round 4
		a = b + rotateLeft(a+I(b, c, d)+x[i+0]+T[48], S41)
		d = a + rotateLeft(d+I(a, b, c)+x[i+7]+T[49], S42)
		c = d + rotateLeft(c+I(d, a, b)+x[i+14]+T[50], S43)
		b = c + rotateLeft(b+I(c, d, a)+x[i+5]+T[51], S44)
		a = b + rotateLeft(a+I(b, c, d)+x[i+12]+T[52], S41)
		d = a + rotateLeft(d+I(a, b, c)+x[i+3]+T[53], S42)
		c = d + rotateLeft(c+I(d, a, b)+x[i+10]+T[54], S43)
		b = c + rotateLeft(b+I(c, d, a)+x[i+1]+T[55], S44)
		a = b + rotateLeft(a+I(b, c, d)+x[i+8]+T[56], S41)
		d = a + rotateLeft(d+I(a, b, c)+x[i+15]+T[57], S42)
		c = d + rotateLeft(c+I(d, a, b)+x[i+6]+T[58], S43)
		b = c + rotateLeft(b+I(c, d, a)+x[i+13]+T[59], S44)
		a = b + rotateLeft(a+I(b, c, d)+x[i+4]+T[60], S41)
		d = a + rotateLeft(d+I(a, b, c)+x[i+11]+T[61], S42)
		c = d + rotateLeft(c+I(d, a, b)+x[i+2]+T[62], S43)
		b = c + rotateLeft(b+I(c, d, a)+x[i+9]+T[63], S44)

		a += aa
		b += bb
		c += cc
		d += dd
	}

	hash := []uint32{a, b, c, d}

	return decode(hash)
}

func calculateMD5(input string) string {
	md5Bytes := md5Transform([]byte(input))
	hexString := ""

	for i := 0; i < 16; i++ {
		hex := fmt.Sprintf("%02x", md5Bytes[i])
		hexString += hex
	}

	return hexString
}

func main() {
	input := "Hello, World!"
	md5Hash := calculateMD5(input)
	fmt.Println("MD5 Hash:", md5Hash)
}
