package main

import (
	"fmt"
)

type Point struct {
	x, y int
}


// Check if point lies on the curve
func isPoint(a1 Point, a, b, p int) bool {
	y := (a1.y * a1.y) % p
	x := (a1.x*a1.x*a1.x + a*a1.x + b) % p
	return x == y
}

// Print all points on the curve
func allPoints(p, a, b int) {
	for i := 0; i < p; i++ {
		for j := 0; j < p; j++ {
			a2 := Point{i, j}
			if isPoint(a2, a, b, p) {
				fmt.Println(i, j)
			}
		}
	}
}




func main() {
	var a, b, p int
	fmt.Print("Enter a: ")
	fmt.Scan(&a)
	fmt.Print("Enter b: ")
	fmt.Scan(&b)
	fmt.Print("Enter p: ")
	fmt.Scan(&p)


	fmt.Println("All Points on the curve are:")
	allPoints(p, a, b)
}
