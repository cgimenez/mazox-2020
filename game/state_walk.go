package game

import "mazox/runtime"

type StateWalk struct {
	CompoundState
	direction int
}

type StatePlayerWalk struct {
	State
	stateWalk StateWalk
}

func (state *StatePlayerWalk) init() {
}

func (state *StateWalk) update(newDirection int) {
	if newDirection != DIRECTION_IDLE && newDirection != state.direction {
		sp := state.owner.setActiveSprite(newDirection - 1)
		sp.rewind()
		sp.play()
	}
	if newDirection == DIRECTION_IDLE {
		state.owner.getActiveSprite().stop()
	} else {
		state.owner.getEntity().move(newDirection)
	}
	state.direction = newDirection
}

func (state *StatePlayerWalk) update() {
	direction := DIRECTION_IDLE
	if gm.rtm.HasKeyPressed(runtime.KEY_RIGHT) {
		direction = DIRECTION_RIGHT
	}
	if gm.rtm.HasKeyPressed(runtime.KEY_LEFT) {
		direction = DIRECTION_LEFT
	}
	if gm.rtm.HasKeyPressed(runtime.KEY_UP) {
		direction = DIRECTION_UP
	}
	if gm.rtm.HasKeyPressed(runtime.KEY_DOWN) {
		direction = DIRECTION_DOWN
	}
	state.stateWalk.update(direction)
}
