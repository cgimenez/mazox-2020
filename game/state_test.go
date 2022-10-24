package game

import (
	"testing"
)

func TestStates(t *testing.T) {
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

	t.Run("states finding", func(t *testing.T) {
		setup()

		state := stateFactory("void")
		ent := entity.(*PlayerEntity)
		entity.addState("S", state, SpriteSlice{})

		st := entity.getState("S")
		if st != state {
			t.Errorf("Expected state to be %p got %p", state, st)
		}

		state_entry := ent.findStateEntry(state)
		if state_entry.state != state {
			t.Errorf("Expected state to be %p got %p", state, state_entry.state)
		}

		teardown()
	})
}
