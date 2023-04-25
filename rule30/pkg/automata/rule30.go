package automata

import "math/rand"

type Grid struct {
	rows   int
	cols   int
	seed   int64
	wrap   bool
	invert bool
	head   *row
	tail   *row
}

type row struct {
	next *row
	data []bool
}

func newRowFromSeed(n int, seed int64) *row {
	r := &row{}
	r.data = make([]bool, n)
	for i := 0; i < n; i++ {
		r.data[i] = seed&1 == 1
		seed >>= 1
	}
	return r
}

func rule30(r *row, wrap, invert bool) *row {
	n := len(r.data)
	next := &row{}
	next.data = make([]bool, n)
	for i := 0; i < n; i++ {
		var left, center, right bool
		if i == 0 {
			if wrap {
				left = r.data[n-1]
			} else {
				left = rand.Float32() > 0.5
			}
		} else {
			left = r.data[i-1]
		}
		if i == n-1 {
			if wrap {
				right = r.data[0]
			} else {
				right = rand.Float32() < 0.5
			}
		} else {
			right = r.data[i+1]
		}
		center = r.data[i]
		next.data[i] = left != (center || right)
		if invert {
			next.data[i] = !next.data[i]
		}
	}
	return next
}

func NewGrid(rows, cols int, seed int64) *Grid {
	g := &Grid{}
	g.rows = rows
	g.cols = cols
	g.seed = seed
	g.init()
	return g
}

func (g *Grid) init() {
	g.head = newRowFromSeed(g.cols, g.seed)
	g.tail = g.head
	for i := 0; i < g.rows; i++ {
		g.tail.next = rule30(g.tail, g.wrap, g.invert)
		g.tail = g.tail.next
	}
}

func (g *Grid) Update() {
	g.head = g.head.next
	newRow := rule30(g.tail, g.wrap, g.invert)
	g.tail.next = newRow
	g.tail = newRow
}

func (g *Grid) ToString() string {
	var s string
	for r := g.head; r != nil; r = r.next {
		for _, c := range r.data {
			if c {
				s += "█"
			} else {
				s += " "
			}
		}
		s += "\n"
	}
	return s
}

func (g *Grid) ToArray() [][]bool {
	var a [][]bool
	for r := g.head; r != nil; r = r.next {
		a = append(a, r.data)
	}
	return a
}

func (g *Grid) ToggleWrap() {
	g.wrap = !g.wrap
}

func (g *Grid) ToggleInvert() {
	g.invert = !g.invert
}

func (g *Grid) IncrementSeed() {
	g.seed++
	g.init()
}

func (g *Grid) DecrementSeed() {
	g.seed--
	g.init()
}
