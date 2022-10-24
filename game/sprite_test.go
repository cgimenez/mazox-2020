package game

import (
	"reflect"
	"testing"
	"time"
)

func TestSprites(t *testing.T) {
	gm = new(Game)
	gm.Init(640, 320)

	var l ILevel
	var entity IEntity

	setup := func() {
		l = levelFactory(1)
		l.load()
		entity = entityFactory("player", "Mister test", l)
	}

	teardown := func() {
		l = nil
		entity = nil
	}

	t.Run("sprite loading", func(t *testing.T) {
		sp := loadSprite(SPRITE_STATIC, "sprite_96.png", 1)
		if len(sp.images) != 1 {
			t.Errorf("Incorrect number of sprite images")
		}

		w, h := sp.getSize()
		if w != 96 || h != 96 {
			t.Errorf("Wrong image size")
		}

		sp.draw()
		if sp.curFrame != 0 {
			t.Errorf("curFrame expected 0 found %d", sp.curFrame)
		}
	})

	t.Run("sprite misc", func(t *testing.T) {
		sp := loadSprite(SPRITE_ANIMATED, "sprite_96.png", 2)
		if len(sp.images) != 2 {
			t.Errorf("Incorrect number of sprite images")
		}
		sp.setFPS(5)

		diff := (1000) * time.Millisecond
		sp.timestamp = time.Now().Add(-diff)
		sp.draw()
		if sp.curFrame != 1 {
			t.Errorf("num_frames expected 0 found %d", sp.curFrame)
		}
	})

	t.Run("sprite states", func(t *testing.T) {
		setup()

		state := stateFactory("void")
		sp0 := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		entity.addSprite(sp0)

		ent := entity.(*PlayerEntity)

		if ent.activeSprite != sp0 {
			t.Errorf("Wrong active sprite")
		}

		ssp1 := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		ssp2 := loadSprite(SPRITE_STATIC, "sprite_32.png", 1)
		entity.addState("S", state, SpriteSlice{ssp1, ssp2})
		state_entry := ent.states[0]

		if state_entry.activeSprite != ssp1 {
			t.Errorf("Got wrong state active sprite")
		}

		if len(ent.sprites) != 3 {
			t.Errorf("Expecting 3 sprites got %d", len(ent.sprites))
		}

		if len(ent.states) != 1 {
			t.Errorf("Got wrong state number")
		}

		if state_entry.id != "S" {
			t.Errorf("Got wrong state id")
		}

		if state_entry.state != state {
			t.Errorf("Got wrong state")
		}

		if !reflect.DeepEqual(state_entry.spritesIndices, []int{1, 2}) {
			t.Errorf("Expecting [1,2] got %v", state_entry.spritesIndices)
		}

		ssp := state.setActiveSprite(0)
		if ssp != ssp1 {
			t.Errorf("Expected state active sprite to be %p got %p", ssp1, ssp)
		}

		if ent.activeSprite != ssp1 {
			t.Errorf("Expected entity active sprite to be %p got %p", ssp1, ssp)
		}

		if state.getActiveSprite() != ssp1 {
			t.Errorf("Expected entity active sprite to be %p got %p", ssp1, state.getActiveSprite())
		}

		entity.draw()

		teardown()
	})

}
