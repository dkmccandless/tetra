package main

import "fmt"

func main() {
	switch m := PlayGame(); m {
	case X, O:
		fmt.Println(m, "wins!")
	default:
		fmt.Println("It's a draw. Good game!")
	}
}

// PlayGame conducts a game and reports the winner, or Empty in the case of a draw.
func PlayGame() Mark {
	g := NewGrid()
	g.Print()

	m := X
	for {
		var c Cell
		fmt.Printf("%v's move: ", m)
		if n, err := fmt.Scan(&c.i, &c.j, &c.k); n != 3 || err != nil {
			fmt.Println(err)
			continue
		}
		if c.i < 0 || c.i > 3 || c.j < 0 || c.j > 3 || c.k < 0 || c.k > 3 {
			fmt.Println("Invalid cell")
			continue
		}

		// Make the move if it's legal
		if ok := g.Move(c, m); !ok {
			fmt.Println("Invalid move")
			continue
		}
		g.Print()

		// Check for game-ending conditions
		switch {
		case g.isWin(c):
			return m
		case g.isFull():
			return Empty
		default:
			m = m.Opp()
		}
	}
}

// Mark is an X or an O in a grid cell.
// The zero value represents an unmarked cell.
type Mark uint8

const (
	Empty Mark = iota
	X
	O
)

// Opp returns m's opponent.
func (m Mark) Opp() Mark {
	if m == X {
		return O
	}
	return X
}

// String returns a string representation of m.
func (m Mark) String() string { return []string{" ", "X", "O"}[m] }

// Cell is the coordinates of a cell in a Grid.
type Cell struct{ i, j, k int }

// Grid is a 4x4x4 cube of cells into which Marks can be placed.
// The zero value represents an empty grid.
type Grid [][][]Mark

// NewGrid returns an empty Grid ready to play.
func NewGrid() Grid {
	g := make([][][]Mark, 4)
	for i := range g {
		g[i] = make([][]Mark, 4)
		for j := range g[i] {
			g[i][j] = make([]Mark, 4)
		}
	}
	return g
}

// Move writes m into g at c if it is legal to do so and reports whether the move is legal.
// A move is legal if c is originally Empty.
func (g Grid) Move(c Cell, m Mark) bool {
	ok := g[c.i][c.j][c.k] == Empty
	if ok {
		g[c.i][c.j][c.k] = m
	}
	return ok
}

// isWin reports whether g contains a straight line of four identical Marks through c.
func (g Grid) isWin(c Cell) bool {
	i, j, k := c.i, c.j, c.k
	// Each cell is part of three lines, one parallel to each axis.
	if sameMark(g[0][j][k], g[1][j][k], g[2][j][k], g[3][j][k]) ||
		sameMark(g[i][0][k], g[i][1][k], g[i][2][k], g[i][3][k]) ||
		sameMark(g[i][j]...) {
		return true
	}

	// Each cell is part of three planes. Within each plane, diagonal lines pass through the plane's four corners and four interior cells.
	// If cell (i, j) is on a plane diagonal, the other three cells on the same diagonal are (i^1, j^1), (i^2, j^2), and (i^3, j^3).
	var onij, onik, onjk bool
	if onDiagonal(i, j) {
		onij = true
		if sameMark(g[i][j][k], g[i^1][j^1][k], g[i^2][j^2][k], g[i^3][j^3][k]) {
			return true
		}
	}
	if onDiagonal(i, k) {
		onik = true
		if sameMark(g[i][j][k], g[i^1][j][k^1], g[i^2][j][k^2], g[i^3][j][k^3]) {
			return true
		}
	}
	if onDiagonal(j, k) {
		onjk = true
		if sameMark(g[i][j][k], g[i][j^1][k^1], g[i][j^2][k^2], g[i][j^3][k^3]) {
			return true
		}
	}

	// The cells on all three plane diagonals are the eight cells at the vertices of the grid and the eight cells in its interior.
	// These cells are also on a fourth diagonal passing through the grid's volume.
	// If cell (i, j, k) is on a volume diagonal, the other three cells on the same diagonal are (i^1, j^1, k^1), (i^2, j^2, k^2), and (i^3, j^3, k^3).
	if onij == onik && onik == onjk {
		return sameMark(g[i][j][k], g[i^1][j^1][k^1], g[i^2][j^2][k^2], g[i^3][j^3][k^3])
	}

	return false
}

// sameMark reports whether all Marks in s are the same.
func sameMark(s ...Mark) bool {
	for i := range s[:len(s)-1] {
		if s[i] != s[i+1] {
			return false
		}
	}
	return true
}

// onDiagonal reports whether the cell represented by plane coordinates i and j lies on a diagonal of the plane.
func onDiagonal(i, j int) bool {
	// The cells in a 4x4 plane that lie on the plane's diagonals are those for which both or neither of i and j are 0 or 3:
	//   0 1 2 3
	// 0 * . . *
	// 1 . * * .
	// 2 . * * .
	// 3 * . . *
	return (i == 0 || i == 3) == (j == 0 || j == 3)
}

// isFull reports whether all cells in g have been filled.
func (g Grid) isFull() bool {
	for i := range g {
		for j := range g[i] {
			for k := range g[i][j] {
				if g[i][j][k] == Empty {
					return false
				}
			}
		}
	}
	return true
}

// Print prints g.
func (g Grid) Print() {
	for i := range g {
		for j := range g[i] {
			for k := range g[i][j] {
				fmt.Printf(" %v ", g[i][j][k])
				if k != len(g[i][j])-1 {
					fmt.Printf("|")
				}
			}
			fmt.Println()
			if j != len(g[i])-1 {
				fmt.Println("---+---+---+---")
			}
		}
		fmt.Println()
	}
}
