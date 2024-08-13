package main

import (
	"fmt"
)

type Point struct {
	x, y int
}

// Calculate modular inverse
func inv(n, p int) int {
	for i := 0; i < p; i++ {
		if (n*i)%p == 1 {
			return i
		}
	}
	return -1
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

// Addition of points on the curve
func addPoints(a1, a3 Point, a, p int) Point {
	if a1.x == a3.x && a1.y != a3.y {
		return a1
	}
	if a1.x == a3.x {
		l := (3*a1.x*a1.x + a) * inv(2*a1.y, p) % p
		xr := (l*l - a1.x - a3.x) % p
		yr := (l*(a1.x-xr) - a1.y) % p
		a3 = Point{xr, yr}
		return a3
	} else {
		l := (a3.y - a1.y) * inv(a3.x-a1.x, p) % p
		xr := (l*l - a1.x - a3.x) % p
		yr := (l*(a1.x-xr) - a1.y) % p
		a3 = Point{xr, yr}
		return a3
	}
}

func genPA(na int, genr Point, a, p int) Point {
	g2 := genr
	for k := 1; k < na; k++ {
		g2 = addPoints(g2, genr, a, p)
	}
	return g2
}

func encrypt(Pm Point, na int, a, b, p int) (Point, Point) {
	var g, h int
	fmt.Print("Enter generator (i,j) coordinate i: ")
	fmt.Scan(&g)
	fmt.Print("Enter generator (i,j) coordinate j: ")
	fmt.Scan(&h)
	genr := Point{g, h}
	pa := genPA(na, genr, a, p)
	k := 41
	c1 := genPA(k, genr, a, p)
	c2 := genPA(k, pa, a, p)
	c2 = addPoints(c2, Pm, a, p)
	return c1, c2
}

func decrypt(na int, a3, a4 Point, a, p int) Point {
	h1 := genPA(na, a3, a, p)
	h1.y = -h1.y
	h2 := addPoints(a4, h1, a, p)
	return h2
}

func main() {
	var a, b, p int
	fmt.Print("Enter a: ")
	fmt.Scan(&a)
	fmt.Print("Enter b: ")
	fmt.Scan(&b)
	fmt.Print("Enter p: ")
	fmt.Scan(&p)

	var x, y int
	fmt.Print("Enter x: ")
	fmt.Scan(&x)
	fmt.Print("Enter y: ")
	fmt.Scan(&y)
	a1 := Point{x, y}

	fmt.Println("Check if point lies on the curve?")
	fmt.Println(isPoint(a1, a, b, p))

	fmt.Println("Adding points:")
	var x1, y1 int
	fmt.Print("Enter xq: ")
	fmt.Scan(&x1)
	fmt.Print("Enter yq: ")
	fmt.Scan(&y1)
	a3 := Point{x1, y1}

	fmt.Println("Sum:")
	a3 = addPoints(a1, a3, a, p)
	fmt.Println(a3.x, a3.y)

	var x2, y2, na int
	fmt.Print("Enter message(x,y) point x: ")
	fmt.Scan(&x2)
	fmt.Print("Enter message(x,y) point y: ")
	fmt.Scan(&y2)
	Pm := Point{x2, y2}
	fmt.Print("Enter na (private key): ")
	fmt.Scan(&na)
	a3, a4 := encrypt(Pm, na, a, b, p)
	fmt.Println("Encrypted message:")
	fmt.Println(a3.x, a3.y, a4.x, a4.y)

	a6 := decrypt(na, a3, a4, a, p)
	fmt.Println("Decrypted message:")
	fmt.Println(a6.x, a6.y)

	fmt.Println("All Points on the curve are:")
	allPoints(p, a, b)
}
