package game

import (
	"testing"
)

func TestEntity(t *testing.T) {
	gm = new(Game)
	gm.Init(640, 320)

	l := levelFactory(1)
	l.load()

	t.Run("setGridPosition", func(t *testing.T) {
		entity := entityFactory("player", "Mister test", l)
		l.addEntity(entity)
		sp := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		entity.setActiveSprite(sp)

		entity.setGridPosition(3, 1)
		ce := entity.(*PlayerEntity)
		if ce.col != 3 || ce.row != 1 {
			t.Errorf("Expected 3 1 got %d %d", ce.col, ce.row)
		}
		if ce.x != 96 || ce.y != 32 {
			t.Errorf("Expected 96 32 got %d %d", ce.x, ce.y)
		}

		entity.setGridPosition(0, 0)
		if ce.col != 0 || ce.row != 0 {
			t.Errorf("Expected 0 0 got %d %d", ce.col, ce.row)
		}
		if ce.x != 0 || ce.y != 0 {
			t.Errorf("Expected 0 0 got %d %d", ce.x, ce.y)
		}

		entity.setGridPosition(19, 9)
		if ce.col != 19 || ce.row != 9 {
			t.Errorf("Expected 19 9 got %d %d", ce.col, ce.row)
		}
		if ce.x != 608 || ce.y != 288 {
			t.Errorf("Expected 608 288 got %d %d", ce.x, ce.y)
		}
	})

	t.Run("setPixPosition", func(t *testing.T) {
		entity := entityFactory("player", "Mister test", l)
		l.addEntity(entity)
		sp := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		entity.setActiveSprite(sp)

		entity.setPixPosition(96, 32)
		ce := entity.(*PlayerEntity)
		if ce.col != 3 || ce.row != 1 {
			t.Errorf("Expected 3 1 got %d %d", ce.col, ce.row)
		}
		if ce.x != 96 || ce.y != 32 {
			t.Errorf("Expected 96 32 got %d %d", ce.x, ce.y)
		}

		entity.setPixPosition(5, 5)
		if ce.col != 0 || ce.row != 0 {
			t.Errorf("Expected 0 0 got %d %d", ce.col, ce.row)
		}
		if ce.x != 5 || ce.y != 5 {
			t.Errorf("Expected 5 5 got %d %d", ce.x, ce.y)
		}

		entity.setPixPosition(610, 290)
		if ce.col != 19 || ce.row != 9 {
			t.Errorf("Expected 19 9 got %d %d", ce.col, ce.row)
		}
		if ce.x != 610 || ce.y != 290 {
			t.Errorf("Expected 610 290 got %d %d", ce.x, ce.y)
		}
	})

	t.Run("getSize", func(t *testing.T) {
		entity := entityFactory("player", "Mister test", l)
		sp := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		entity.setActiveSprite(sp)
		h, w := entity.getSize()
		if w != 32 || h != 32 {
			t.Errorf("Expected 32 32 got %d %d", w, h)
		}
		sp = loadSprite(SPRITE_STATIC, "sprite_96.png", 1)
		entity.setActiveSprite(sp)
		h, w = entity.getSize()
		if w != 96 || h != 96 {
			t.Errorf("Expected 96 96 got %d %d", w, h)
		}
	})

}
