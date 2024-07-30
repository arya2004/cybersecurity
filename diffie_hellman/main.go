package main

import (
	"fmt"
	"math/big"
)

// Power function to return value of a ^ b mod P
func power(a, b, P int64) int64 {
	// Convert inputs to big.Int
	A := big.NewInt(a)
	B := big.NewInt(b)
	Pmod := big.NewInt(P)
	
	// Perform modular exponentiation
	result := new(big.Int).Exp(A, B, Pmod)
	
	return result.Int64()
}

// Encrypt function using XOR
func encryptDecrypt(message string, key int64) string {
	var result string
	for _, char := range message {
		result += string(char ^ rune(key))
	}
	return result
}

func main() {
	var P, G, x, a, y, b, ka, kb int64

	// Both the persons will be agreed upon the
	// public keys G and P
	P = 23 // A prime number P is taken
	fmt.Println("The value of P:", P)

	G = 9 // A primitive root for P, G is taken
	fmt.Println("The value of G:", G)

	// Alice will choose the private key a
	a = 4 // a is the chosen private key
	fmt.Println("The private key a for Alice:", a)

	x = power(G, a, P) // gets the generated key
	fmt.Println("The public key x for Alice:", x)

	// Bob will choose the private key b
	b = 3 // b is the chosen private key
	fmt.Println("The private key b for Bob:", b)

	y = power(G, b, P) // gets the generated key
	fmt.Println("The public key y for Bob:", y)

	// Generating the secret key after the exchange of keys
	ka = power(y, a, P) // Secret key for Alice
	kb = power(x, b, P) // Secret key for Bob
	fmt.Println("Secret key for the Alice is:", ka)
	fmt.Println("Secret key for the Bob is:", kb)

	// Alice encrypts a message
	message := "Hello Bob!"
	encryptedMessage := encryptDecrypt(message, ka)
	fmt.Println("Encrypted Message:", encryptedMessage)

	// Bob decrypts the message
	decryptedMessage := encryptDecrypt(encryptedMessage, kb)
	fmt.Println("Decrypted Message:", decryptedMessage)
}
