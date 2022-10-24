package game

import (
	"math/rand"
	"time"
)

type Vector2 struct {
	x, y int
}

func randomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func surroundingCellsCoords(ec int, er int, maxCols int, maxRows int) []Vector2 {
	var cells []Vector2
	cells = append(cells, []Vector2{
		{ec - 1, er - 1}, {ec, er - 1}, {ec + 1, er - 1},
		{ec - 1, er}, {ec, er}, {ec + 1, er},
		{ec - 1, er + 1}, {ec, er + 1}, {ec + 1, er + 1},
	}...)
	for i, cell := range cells {
		if cell.x == 0 && cell.y == 0 && cell.x > maxCols && cell.y > maxRows {
			cells[i] = cells[len(cells)-1]
			cells = cells[:len(cells)-1] // slice order is lost
		}
	}
	return cells
}
