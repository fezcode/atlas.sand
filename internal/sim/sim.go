package sim

import (
	"math/rand"
)

type ParticleType int

const (
	Empty ParticleType = iota
	Wall
	Sand
	Water
	Fire
	Salt
)

type Cell struct {
	Type    ParticleType
	Color   string
	Visited bool // To avoid double-updating in a single frame
}

type Simulation struct {
	Width  int
	Height int
	Grid   [][]Cell
}

func NewSimulation(width, height int) *Simulation {
	grid := make([][]Cell, height)
	for i := range grid {
		grid[i] = make([]Cell, width)
	}
	return &Simulation{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

func (s *Simulation) SetCell(x, y int, t ParticleType) {
	if x >= 0 && x < s.Width && y >= 0 && y < s.Height {
		s.Grid[y][x] = Cell{Type: t}
	}
}

func (s *Simulation) Clear() {
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.Grid[y][x] = Cell{Type: Empty}
		}
	}
}

func (s *Simulation) Update() {
	// Reset visited status
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.Grid[y][x].Visited = false
		}
	}

	// Update from bottom to top to handle gravity correctly
	for y := s.Height - 1; y >= 0; y-- {
		// Randomize horizontal order to avoid bias
		order := rand.Perm(s.Width)
		for _, x := range order {
			s.updateCell(x, y)
		}
	}
}

func (s *Simulation) updateCell(x, y int) {
	cell := s.Grid[y][x]
	if cell.Type == Empty || cell.Type == Wall || cell.Visited {
		return
	}

	switch cell.Type {
	case Sand:
		s.updateSand(x, y)
	case Water:
		s.updateWater(x, y)
	case Salt:
		s.updateSand(x, y) // Salt behaves like sand for now
	case Fire:
		s.updateFire(x, y)
	}
}

func (s *Simulation) updateSand(x, y int) {
	if y+1 >= s.Height {
		return
	}

	// 1. Move straight down
	if s.Grid[y+1][x].Type == Empty || s.Grid[y+1][x].Type == Water {
		s.swap(x, y, x, y+1)
		return
	}

	// 2. Move diagonally down
	dir := 1
	if rand.Float32() < 0.5 {
		dir = -1
	}

	if s.canMove(x+dir, y+1) {
		s.swap(x, y, x+dir, y+1)
	} else if s.canMove(x-dir, y+1) {
		s.swap(x, y, x-dir, y+1)
	}
}

func (s *Simulation) updateWater(x, y int) {
	if y+1 < s.Height && s.Grid[y+1][x].Type == Empty {
		s.swap(x, y, x, y+1)
		return
	}

	// Try move sideways
	dir := 1
	if rand.Float32() < 0.5 {
		dir = -1
	}

	if s.canMove(x+dir, y+1) {
		s.swap(x, y, x+dir, y+1)
	} else if s.canMove(x-dir, y+1) {
		s.swap(x, y, x-dir, y+1)
	} else if s.canMove(x+dir, y) {
		s.swap(x, y, x+dir, y)
	} else if s.canMove(x-dir, y) {
		s.swap(x, y, x-dir, y)
	}
}

func (s *Simulation) updateFire(x, y int) {
	// Fire moves up or disappears
	if rand.Float32() < 0.1 {
		s.Grid[y][x] = Cell{Type: Empty}
		return
	}

	newY := y - 1
	if newY < 0 {
		s.Grid[y][x] = Cell{Type: Empty}
		return
	}

	dir := rand.Intn(3) - 1 // -1, 0, 1
	newX := x + dir

	if newX >= 0 && newX < s.Width && s.Grid[newY][newX].Type == Empty {
		s.swap(x, y, newX, newY)
	} else if rand.Float32() < 0.2 {
		s.Grid[y][x] = Cell{Type: Empty}
	}
}

func (s *Simulation) canMove(x, y int) bool {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
		return false
	}
	return s.Grid[y][x].Type == Empty
}

func (s *Simulation) swap(x1, y1, x2, y2 int) {
	s.Grid[y1][x1], s.Grid[y2][x2] = s.Grid[y2][x2], s.Grid[y1][x1]
	s.Grid[y2][x2].Visited = true
}
