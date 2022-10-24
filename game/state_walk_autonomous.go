package game

type StateAutonomousWalk struct {
	State
	stateWalk                StateWalk
	timer                    Timer
	dest_col, dest_row       int
	next_cell_x, next_cell_y int
	next_cell_reached        bool
	moving                   bool
	direction                int
	new_direction            int
}

func (state *StateAutonomousWalk) init() {
	state.timer = Timer{}
}

func (state *StateAutonomousWalk) enter() {
	directions := []int{DIRECTION_DOWN, DIRECTION_UP, DIRECTION_LEFT, DIRECTION_RIGHT}
	state.new_direction = directions[randomInt(0, len(directions)-1)]
	state.moving = true
	state.timer.start(3000)
}

func (state *StateAutonomousWalk) update() {
	if state.moving {
		if state.new_direction != DIRECTION_IDLE && state.new_direction != state.direction {
			sp := state.setActiveSprite(state.new_direction - 1)
			sp.rewind()
			sp.play()
		}
		if state.new_direction == DIRECTION_IDLE {
			state.getActiveSprite().stop()
		} else {
			if !state.getEntity().move(state.new_direction) {
				state.enter()
			}
		}
		state.direction = state.new_direction
		state.stateWalk.update(state.direction)
	}
}
