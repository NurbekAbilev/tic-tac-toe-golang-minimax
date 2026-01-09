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

// (2*(9+1)^10)
// amount of boxes = 9
// extra twobit for to indicate move (00 for X and 01 for Y)
const MOVECOUNT = 1048576

var s [MOVECOUNT]byte
var vis [MOVECOUNT]bool

var globalRecCount int = 0

func calc(b int) byte {
	// fmt.Println("state:")
	// printBinary(b)

	// if globalRecCount > 30 {
	// 	log.Fatal("finish")
	// }

	globalRecCount++
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

	t := (b >> (9 * 2)) & 1

	var res byte = Draw // draw
	if t == 0 {
		res = 0
	} else {
		res = 255
	}

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			v := getSquare(b, y, x)
			if v != 0 {
				continue
			}

			sv := Xf
			if t == 1 {
				sv = Yf
			}

			idx := ((3 * y) + x) * 2
			newb := b | (sv << idx)
			newb ^= 1 << (9 * 2)

			gameVal := calc(newb)

			if t == 0 {
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

type Move struct {
	Y     int
	X     int
	Value byte
}

func getAvailableMoves(b int) []Move {
	moves := make([]Move, 0)
	t := (b >> (9 * 2)) & 1
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			v := getSquare(b, y, x)
			if v != 0 {
				continue
			}

			sv := Xf
			if t == 1 {
				sv = Yf
			}

			idx := ((3 * y) + x) * 2
			newb := b | (sv << idx)
			newb ^= 1 << (9 * 2)

			moves = append(moves, Move{
				y, x, s[newb],
			})
		}
	}

	return moves
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

	mask := 0b0011
	valByte := (b & (mask << idx)) >> idx
	return byte(valByte)
}

const (
	Ywon        = 1
	Draw        = 2
	Xwon        = 3
	NotFinished = 4
)

func checkWin(b int) byte {
	// 1. Check Rows
	for y := 0; y < 3; y++ {
		s1 := getSquare(b, y, 0)
		if s1 != empty && s1 == getSquare(b, y, 1) && s1 == getSquare(b, y, 2) {
			if s1 == Xf {
				return Xwon
			}
			return Ywon
		}
	}
	// 2. Check Columns
	for x := 0; x < 3; x++ {
		s1 := getSquare(b, 0, x)
		if s1 != empty && s1 == getSquare(b, 1, x) && s1 == getSquare(b, 2, x) {
			if s1 == Xf {
				return Xwon
			}
			return Ywon
		}
	}
	// 3. Check Diagonal (Top-Left to Bottom-Right)
	sDiag1 := getSquare(b, 0, 0)
	if sDiag1 != empty && sDiag1 == getSquare(b, 1, 1) && sDiag1 == getSquare(b, 2, 2) {
		if sDiag1 == Xf {
			return Xwon
		}
		return Ywon
	}

	// 4. Check Anti-Diagonal (Top-Right to Bottom-Left)
	sDiag2 := getSquare(b, 0, 2)
	if sDiag2 != empty && sDiag2 == getSquare(b, 1, 1) && sDiag2 == getSquare(b, 2, 0) {
		if sDiag2 == Xf {
			return Xwon
		}
		return Ywon
	}

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if getSquare(b, y, x) == empty {
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
	fmt.Println((b >> (9 * 2)) & 1)

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

	// fmt.Println("Whose turn is it[x/y]:")
	input, err := reader.ReadString('\n')
	if err != nil {
		// fmt.Println("Error reading input:", err)
		return 0
	}
	turn := 1
	if input[0] == 'x' {
		turn = 0
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

	check := checkWin(board)
	fmt.Println("check:", check)

	c := calc(0)
	_ = c

	fmt.Println("s[board]", s[board])

	moves := getAvailableMoves(board)
	fmt.Printf("%+v\n", moves)
	for _, mv := range moves {
		fmt.Printf("y=%d, x=%d, value=%d\n", mv.Y, mv.X, mv.Value)
	}

	for i := 0; i < MOVECOUNT; i++ {
		if s[i] == Xwon && vis[i] {
			printBoard(i)
			fmt.Println()
		}
	}
	fmt.Println(c)

	// limit := 0
	// for i := 0; i < MOVECOUNT; i++ {
	// 	check := checkWin(i)
	// 	if check != NotFinished {
	// 		printBoard(i)
	// 		limit++
	// 	}

	// 	if limit > 100 {
	// 		break
	// 	}
	// }
}
