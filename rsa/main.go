package main

import (
	"fmt"
	"math"
)

// Variables for the prime numbers p and q
var p int64 = 3
var q int64 = 7

// Variables for the public key exponent 'e', a multiplier 'k', and the message 'msg'
var e int64 = 2
var k int64 = 2
var msg int64 = 12

// Function to calculate the greatest common divisor (GCD) of two numbers
func gcd(a int64, h int64) int64 {
	var temp int64
	for {
		temp = a % h
		if temp == 0 {
			return h
		}
		a = h
		h = temp
	}
}

func main() {
	// Calculate n which is the product of p and q
	n := p * q
	fmt.Println("Calculated n (p * q) = ", n)

	// Calculate Euler's Totient function (phi) for n
	phi := (p - 1) * (q - 1)
	fmt.Println("Calculated phi ((p - 1) * (q - 1)) = ", phi)

	// Find a value for e such that 1 < e < phi and gcd(e, phi) == 1

	if gcd(e, phi) != 1 {
		
		fmt.Println("e =", e, "is not coprime with phi. Incrementing e.")

	}

	for e < phi {
		if gcd(e, phi) == 1 {
			fmt.Println("Chosen e = ", e, " as it is coprime with phi")
			break
		} else {
			e++
		}
	}

	// Calculate the private key exponent 'd' using the formula: d = (1 + (k * phi)) / e
	var d int64 = (1 + (k * phi)) / e
	fmt.Println("Calculated d (private key) = ", d)

	// Print the original message data
	fmt.Println("Message data = ", msg)

	// Encrypt the message using the formula: c = (msg^e) % n
	c := int64(math.Mod(math.Pow(float64(msg), float64(e)), float64(n)))
	fmt.Println("Encrypted data = ", c)

	// Decrypt the message using the formula: m = (c^d) % n
	m := int64(math.Mod(math.Pow(float64(c), float64(d)), float64(n)))
	fmt.Println("Original Message Sent = ", m)
}
