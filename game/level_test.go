package game

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

//
// blocked cells 32 squared pix grid
// - at col,row => 1,1 (32,32 <-> 63,63)
func TestLevel(t *testing.T) {
	gm = new(Game)
	gm.Init(640, 320)

	l := Level1{Level{num: 1}}
	l.load()

	t.Run("getGridCellSize", func(t *testing.T) {
		s := l.getGridCellSize()
		if s != 32 {
			t.Errorf("s expected 32 found %d", s)
		}
	})

	t.Run("Grid cell blocked", func(t *testing.T) {
		for row := 0; row < l.gridRows; row++ {
			for col := 0; col < l.gridCols; col++ {
				if row == 1 && col == 1 {
					if l.grid[row][col] != LEVEL_GRID_BLOCKED {
						t.Errorf("Expecting LEVEL_GRID_BLOCKED %d %d", row, col)
					}
				} else {
					if l.grid[row][col] != LEVEL_GRID_FREE {
						t.Errorf("Expecting LEVEL_GRID_FREE %d %d", row, col)
					}
				}
			}
		}
	})

	t.Run("translatePixToGridLoc", func(t *testing.T) {
		ex := []struct {
			x, y, c, r int
		}{
			{0, 0, 0, 0},
			{31, 31, 0, 0},
			{32, 32, 1, 1},
			{33, 33, 1, 1},
			{639, 319, 19, 9},
		}

		for _, v := range ex {
			c, r := l.translatePixToGridLoc(v.x, v.y)
			if c != v.c || r != v.r {
				t.Errorf("Expected %d %d got %d %d", v.c, v.r, c, r)
			}
		}
	})

}

func TestEntityCanMove(t *testing.T) {
	gm = new(Game)
	gm.Init(640, 320)

	var l ILevel
	var entity IEntity

	setup := func() {
		l = levelFactory(1)
		l.load()
		entity = entityFactory("player", "Mister test", l)
		sp := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		entity.setActiveSprite(sp)
	}

	teardown := func() {
		l = nil
		entity = nil
	}

	t.Run("entityCanMoveToPix GRID BLOCKED", func(t *testing.T) {
		setup()
		l.addEntity(entity)

		type rect struct{ x1, y1, x2, y2 int }
		rects := []rect{
			{1, 1, 63, 63},
		}

		inside := func(x, y int, r rect) bool {
			return x >= r.x1 && x <= r.x2 && y >= r.y1 && y <= r.y2
		}

		var reject_impossible, reject_possible []string

		sw, sh := entity.getSize()
		for x := -1; x < 640; x++ {
			for y := -1; y < 320; y++ {
				entity.setPixPosition(x, y)
				is_inside := false
				for _, r := range rects {
					if inside(x, y, r) {
						is_inside = true
						break
					}
				}
				can_move := l.entityCanMoveToPix(entity, x, y)
				if is_inside || x < 0 || y < 0 || x >= 640-sw || y >= 320-sh {
					if can_move {
						t.Errorf("Expected impossible move x = %d y = %d", x, y)
						reject_impossible = append(reject_impossible, fmt.Sprintf("Expected impossible move x = %d y = %d", x, y))
						os.Exit(1)
					}
				} else {
					if !can_move {
						reject_possible = append(reject_possible, fmt.Sprintf("Expected possible move x = %d y = %d", x, y))
						t.Errorf("Expected possible move x = %d y = %d", x, y)
						os.Exit(1)
					}
				}
			}
		}
		if len(reject_impossible) > 0 {
			t.Errorf("Impossible moves rejected %d\nThe first one was %s", len(reject_impossible), reject_impossible[0])
		}
		if len(reject_possible) > 0 {
			t.Errorf("Possible moves rejected %d\nThe first one was %s", len(reject_possible), reject_possible[0])
		}

		teardown()
	})

	t.Run("entityCanMoveToCell", func(t *testing.T) {
		setup()

		entity2 := entityFactory("player", "Mister test2", l)
		l.addEntity(entity)
		l.addEntity(entity2)
		entity.setGridPosition(5, 5)
		entity2.setGridPosition(3, 3)

		for col := 0; col < 20; col++ {
			for row := 0; row < 10; row++ {
				if (row == 1 && col == 1) || (row == 3 && col == 3) {
					if l.entityCanMoveToCell(entity, col, row) {
						t.Errorf("Expected impossible move x = %d y = %d", col, row)
					}
				} else {
					if !l.entityCanMoveToCell(entity, col, row) {
						t.Errorf("Expected possible move x = %d y = %d", col, row)
					}
				}
			}
		}

		teardown()
	})

	t.Run("add and remove entity to or from level", func(t *testing.T) {
		setup()

		l.addEntity(entity)
		ce := l.(*Level1)
		if ce.entities[0] != entity {
			t.Errorf("Entity should be in level")
		}
		//fmt.Printf("%p %p\n", l, entity.getLevel())
		if reflect.ValueOf(entity.getLevel()).Pointer() != reflect.ValueOf(l).Pointer() {
			t.Errorf("Entity should belongs to level")
		}

		teardown()
	})
}
