package main

import (
	"reflect"
	"testing"
)

// constructGrid returns a Grid with X and O Marks occupying the Cells in xs and os, if these are disjoint.
// It also reports whether any cell would be marked more than once: if so, it returns an empty Grid instead.
func constructGrid(t *testing.T, xs, os []Cell) (Grid, bool) {
	g := NewGrid()
	for _, c := range xs {
		if !g.Move(c, X) {
			return NewGrid(), false
		}
	}
	for _, c := range os {
		if !g.Move(c, O) {
			return NewGrid(), false
		}
	}
	return g, true
}

func TestMove(t *testing.T) {
	for _, test := range []struct {
		xs, os []Cell
		c      Cell
		ok     bool
	}{
		{[]Cell{}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}}, []Cell{}, Cell{0, 0, 0}, false},
		{[]Cell{{3, 3, 0}, {3, 3, 1}, {3, 3, 2}}, []Cell{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}}, Cell{3, 3, 3}, true},
		{[]Cell{{3, 3, 0}, {3, 3, 1}, {3, 3, 2}}, []Cell{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}}, Cell{3, 3, 3}, false},
	} {
		xg, xok := constructGrid(t, append(test.xs, test.c), test.os) // X plays at c
		og, ook := constructGrid(t, test.xs, append(test.os, test.c)) // O plays at c
		if xok != ook || xok != test.ok {
			t.Errorf("TestMove(%v, %v, %v): got valid cell to move for X %v, for O %v; want %v", test.xs, test.os, test.c, xok, ook, test.ok)
		}
		grids := []Grid{xg, og}
		for i, m := range []Mark{X, O} {
			g, ok := constructGrid(t, test.xs, test.os)
			if !ok {
				t.Fatal(test.xs, test.os)
			}
			if ok = g.Move(test.c, m); ok != test.ok {
				t.Errorf("TestMove(%v, %v, %v, %v): got valid move %v, want %v", test.xs, test.os, test.c, m, ok, test.ok)
			}
			var want Grid
			if test.ok {
				want = grids[i]
			} else {
				// Move does not change g if the move is invalid
				want, _ = constructGrid(t, test.xs, test.os)
			}
			if !reflect.DeepEqual(g, want) {
				t.Errorf("TestMove(%v, %v, %v, %v): got %v, want %v", test.xs, test.os, test.c, m, g, want) // TODO: more meaningful error message
			}
		}
	}
}

func TestIsWin(t *testing.T) {
	for _, test := range []struct {
		xs, os []Cell
		c      Cell
		want   bool
	}{
		{[]Cell{{0, 0, 0}}, []Cell{}, Cell{0, 0, 0}, false},
		{[]Cell{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}}, []Cell{}, Cell{0, 0, 0}, false},
		{[]Cell{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}, {0, 0, 3}}, []Cell{}, Cell{0, 0, 0}, false},
		{[]Cell{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}}, []Cell{{0, 0, 3}}, Cell{0, 0, 0}, false},
		{[]Cell{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}}, []Cell{{0, 0, 3}}, Cell{0, 0, 3}, false},
		{[]Cell{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}, {0, 0, 3}, {3, 3, 3}}, []Cell{}, Cell{3, 3, 3}, false},
		{[]Cell{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}, {0, 0, 3}}, []Cell{{3, 3, 3}}, Cell{3, 3, 3}, false},

		{[]Cell{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}, {0, 0, 3}}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}, {0, 1, 0}, {0, 2, 0}, {0, 3, 0}}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}, {1, 0, 0}, {2, 0, 0}, {3, 0, 0}}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}, {0, 1, 1}, {0, 2, 2}, {0, 3, 3}}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}, {1, 0, 1}, {2, 0, 2}, {3, 0, 3}}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}, {1, 1, 0}, {2, 2, 0}, {3, 3, 0}}, []Cell{}, Cell{0, 0, 0}, true},
		{[]Cell{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}, {3, 3, 3}}, []Cell{}, Cell{0, 0, 0}, true},

		{[]Cell{{0, 0, 1}, {0, 0, 2}, {0, 0, 3}, {0, 0, 0}}, []Cell{}, Cell{0, 0, 1}, true},
		{[]Cell{{0, 0, 1}, {0, 1, 1}, {0, 2, 1}, {0, 3, 1}}, []Cell{}, Cell{0, 0, 1}, true},
		{[]Cell{{0, 0, 1}, {1, 0, 1}, {2, 0, 1}, {3, 0, 1}}, []Cell{}, Cell{0, 0, 1}, true},
		{[]Cell{{0, 0, 1}, {0, 1, 2}, {0, 2, 3}, {0, 3, 0}}, []Cell{}, Cell{0, 0, 1}, false},
		{[]Cell{{0, 0, 1}, {1, 0, 2}, {2, 0, 3}, {3, 0, 0}}, []Cell{}, Cell{0, 0, 1}, false},
		{[]Cell{{0, 0, 1}, {1, 1, 1}, {2, 2, 1}, {3, 3, 1}}, []Cell{}, Cell{0, 0, 1}, true},

		{[]Cell{{0, 1, 1}, {0, 1, 2}, {0, 1, 3}, {0, 1, 0}}, []Cell{}, Cell{0, 1, 1}, true},
		{[]Cell{{0, 1, 1}, {0, 2, 1}, {0, 3, 1}, {0, 0, 1}}, []Cell{}, Cell{0, 1, 1}, true},
		{[]Cell{{0, 1, 1}, {1, 1, 1}, {2, 1, 1}, {3, 1, 1}}, []Cell{}, Cell{0, 1, 1}, true},
		{[]Cell{{0, 1, 1}, {0, 2, 2}, {0, 3, 3}, {0, 0, 0}}, []Cell{}, Cell{0, 1, 1}, true},
		{[]Cell{{0, 1, 1}, {1, 1, 2}, {2, 1, 3}, {3, 1, 0}}, []Cell{}, Cell{0, 1, 1}, false},
		{[]Cell{{0, 1, 1}, {1, 2, 1}, {2, 3, 1}, {3, 0, 1}}, []Cell{}, Cell{0, 1, 1}, false},

		{[]Cell{{1, 1, 1}, {1, 1, 2}, {1, 1, 3}, {1, 1, 0}}, []Cell{}, Cell{1, 1, 1}, true},
		{[]Cell{{1, 1, 1}, {1, 2, 1}, {1, 3, 1}, {1, 0, 1}}, []Cell{}, Cell{1, 1, 1}, true},
		{[]Cell{{1, 1, 1}, {2, 1, 1}, {3, 1, 1}, {0, 1, 1}}, []Cell{}, Cell{1, 1, 1}, true},
		{[]Cell{{1, 1, 1}, {1, 2, 2}, {1, 3, 3}, {1, 0, 0}}, []Cell{}, Cell{1, 1, 1}, true},
		{[]Cell{{1, 1, 1}, {2, 1, 2}, {3, 1, 3}, {0, 1, 0}}, []Cell{}, Cell{1, 1, 1}, true},
		{[]Cell{{1, 1, 1}, {2, 2, 1}, {3, 3, 1}, {0, 0, 1}}, []Cell{}, Cell{1, 1, 1}, true},
		{[]Cell{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}, {0, 0, 0}}, []Cell{}, Cell{1, 1, 1}, true},
	} {
		g, ok := constructGrid(t, test.xs, test.os)
		if !ok {
			t.Fatal(test.xs, test.os)
		}
		if got := g.isWin(test.c); got != test.want {
			t.Errorf("isWin(%v, %v %v): got %v, want %v", test.xs, test.os, test.c, got, test.want)
		}
	}
}
