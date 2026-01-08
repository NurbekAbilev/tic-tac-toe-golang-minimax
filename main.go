package main

import (
	"bufio"
	"fmt"
	"os"
)

// 00 -> empty
// 01 -> X
// 10 -> Y
// 10th bit represents whose turn it is where bit = 1 means X turn and bit = 0 means "O" turn
// [00,00][00,00][00,00]... [0]-> bestmove

const (
	empty = 0 // 00
	Xf    = 1 // 01
	Yf    = 2 // 10
)

const MOVECOUNT = 59049 // 3^10
var s [1048576]byte
var vis [1048576]bool

func calc(b int, t bool) byte {
	if vis[b] {
		return s[b]
	}

	check := checkWin(b)
	if check != NotFinished {
		vis[b] = true
		s[b] = check
		// fmt.Println("board won cond:")
		// printBoard(b)
		return check
	}

	var res byte = Draw // draw

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			v := getSquare(b, y, x)
			if v != 0 {
				continue
			}

			sv := Xf
			if t == true {
				sv = Yf
			}

			idx := ((3 * y) + x) * 2
			newb := b | (sv << idx)

			gameVal := calc(newb, !t)

			if t == false {
				res = max(res, gameVal)
			} else {
				res = min(res, gameVal)
			}
		}
	}

	vis[b] = true
	s[b] = res

	return res
}

func max(a, b byte) byte {
	if a > b {
		return a
	}
	return b
}

func min(a, b byte) byte {
	if a < b {
		return a
	}
	return b
}

func getSquare(b, y, x int) byte {
	idx := ((3 * y) + x) * 2
	fb := (b >> idx) & 1
	sb := (b >> (idx + 1)) & 1
	if sb == 1 {
		return Xf
	} else if fb == 1 {
		return Yf
	}

	return 0
}

const (
	Ywon        = 0
	Draw        = 1
	Xwon        = 2
	NotFinished = 3
)

func checkWin(b int) byte {
	// 1. Check Rows
	for y := 0; y < 3; y++ {
		s1 := getSquare(b, y, 0)
		if s1 != 0 && s1 == getSquare(b, y, 1) && s1 == getSquare(b, y, 2) {
			if s1 == Xf {
				return Xwon
			}
			return Ywon
		}
	}
	// 2. Check Columns
	for x := 0; x < 3; x++ {
		s1 := getSquare(b, 0, x)
		if s1 != 0 && s1 == getSquare(b, 1, x) && s1 == getSquare(b, 2, x) {
			if s1 == Xf {
				return Xwon
			}
			return Ywon
		}
	}
	// 3. Check Diagonal (Top-Left to Bottom-Right)
	sDiag1 := getSquare(b, 0, 0)
	if sDiag1 != 0 && sDiag1 == getSquare(b, 1, 1) && sDiag1 == getSquare(b, 2, 2) {
		if sDiag1 == Xf {
			return Xwon
		}
		return Ywon
	}
	// 4. Check Anti-Diagonal (Top-Right to Bottom-Left)
	sDiag2 := getSquare(b, 0, 2)
	if sDiag2 != 0 && sDiag2 == getSquare(b, 1, 1) && sDiag2 == getSquare(b, 2, 0) {
		if sDiag2 == Xf {
			return Xwon
		}
		return Ywon
	}

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if getSquare(b, y, x) == 0 {
				return NotFinished
			}
		}
	}

	// 5. No winner found
	return Draw
}

func printBoard(b int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			v := getSquare(b, i, j)
			if v == Xf {
				fmt.Print("x")
			} else if v == Yf {
				fmt.Print("y")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func readFromInput(reader *bufio.Reader) int {
	// var rd bufio.Reader
	b := 0
	for y := 0; y < 3; y++ {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return 0
		}

		for x := 0; x < 3; x++ {
			val := empty

			if input[x] == 'x' {
				val = Xf
			} else if input[x] == 'y' {
				val = Yf
			}

			fb := val & 1
			sb := (val >> 1) & 1

			// fmt.Println("fb =", fb, " sb = ", sb)

			idx := ((3 * y) + x) * 2
			fbb := (fb << idx)
			sbb := (sb << (idx + 1))

			// fmt.Printf("%010b %010b\n", fbb, sbb)
			b |= fbb
			b |= sbb
		}
	}
	// 10,10,10,10,10,10,01,01,01
	// 11001100110011001

	// fmt.Println("Whose turn is it[x/y]:")
	input, err := reader.ReadString('\n')
	if err != nil {
		// fmt.Println("Error reading input:", err)
		return 0
	}
	turn := 0
	if input[0] == 'x' {
		turn = 1
	}
	// fmt.Printf("Turn binary: %022b\n", (turn << (9 * 2)))
	b |= (turn << (9 * 2))

	return b
}

func printBinary(a int) {
	for y := 0; y < 3; y++ {
		fmt.Print("[")
		for x := 0; x < 3; x++ {
			idx := ((3 * y) + x) * 2
			mask := 0b0011

			valByte := (a & (mask << idx)) >> idx

			fmt.Printf("%02b", valByte)
			// fmt.Print((a >> (idx + 1)) & 1)
			// fmt.Print((a >> idx) & 1)
			if x < 2 {
				fmt.Print(",")
			}
		}
		fmt.Print("]\n")
	}
	fmt.Printf("Turn bit is: 0%d\n", (a>>(9*2))&1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter board:\n")
	board := readFromInput(reader)

	fmt.Printf("Binary board: \n")
	printBinary(board)
	printBoard(board)

	c := calc(board, false)

	fmt.Println(c)
}
