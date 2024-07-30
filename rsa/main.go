// package main

// import "fmt"

// func gcd(a int64, h int64) int64{
// 	var temp int64
// 	for {
// 		temp =  a % h
// 		if(temp == 0){
// 			return h
// 		}
// 		a = h
// 		h = temp

// 	}
// }

// func main(){

// 	var p float64 = 3
// 	var q float64 = 11

// 	n := p * q

// 	var e float64 = 2

// 	phi := (p - 1) * (q - 1)

// 	for e < phi {
// 		if(gcd(int64(e), int64(phi)) == 1){
// 			break
// 		}else{
// 			e++

// 		}
// 	}

// 	var k float64 = 2

// 	var d float64 = (1 + (k * phi)) / e

// 	var msg float64 = 12

// 	fmt.Println("Message data = ", msg)

// 	var c float64 = float64(msg) ^ e % n

// 	fmt.Println("Encrypted data = ", c)

// 	var m float64 = c ^ d % n

// 	fmt.Println("Original Message Sent = ", m)

// 	fmt.Print("sdgf\n")
// }

package main

import (
	"fmt"
	"math"
)

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
	var p int64 = 3
	var q int64 = 7

	n := p * q

	var e int64 = 2

	phi := (p - 1) * (q - 1)

	for e < phi {
		if gcd(e, phi) == 1 {
			break
		} else {
			e++
		}
	}

	var k int64 = 2

	var d int64 = (1 + (k * phi)) / e

	var msg int64 = 12

	fmt.Println("Message data = ", msg)

	// Encrypt the message
	c := int64(math.Mod(math.Pow(float64(msg), float64(e)), float64(n)))

	fmt.Println("Encrypted data = ", c)

	// Decrypt the message
	m := int64(math.Mod(math.Pow(float64(c), float64(d)), float64(n)))

	fmt.Println("Original Message Sent = ", m)

}
