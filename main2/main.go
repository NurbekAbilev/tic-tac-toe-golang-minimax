package main

import "fmt"

const (
	empty = 0 // 00
	Xf    = 1 // 01
	Yf    = 2 // 10
)

func getBoard(b int) [3][3]byte {
	arr := [3][3]byte{}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			idx := ((3 * y) + x) * 2
			fb := (b >> idx) & 1
			sb := (b >> (idx + 1)) & 1

			if sb == 1 {
				arr[y][x] = Xf
			} else if fb == 1 {
				arr[y][x] = Yf
			} else {
				arr[y][x] = empty
			}
		}
	}
	return arr
}

func main() {
	a := 52

	fmt.Printf("%b\n", a)

	for i := 0; i < 15; i++ {
		fmt.Printf("%d 'nth bit is %d\n", i, a>>i&1)
	}

}
