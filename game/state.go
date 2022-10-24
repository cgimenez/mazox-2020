package game

import (
	"fmt"
	"log"
)

type IState interface {
	init()
	enter()
	update()
	leave()

	setID(string)
	getID() string
	getEntity() IEntity
	setEntity(IEntity)

	setActiveSprite(int) *Sprite
	getActiveSprite() *Sprite
}

type State struct {
	id     string
	entity IEntity
}

type CompoundState struct {
	owner IState
}

func stateFactory(kind string) IState {
	var new_state IState

	switch kind {
	case "player_walk":
		state := &StatePlayerWalk{}
		state.stateWalk.owner = state
		new_state = state

	case "random_walk":
		state := &StateAutonomousWalk{}
		state.stateWalk.owner = state
		new_state = state

	case "void":
		state := &StateVoid{}
		new_state = state
	default:
		panic(fmt.Sprintf("State kind %s is unknown", kind))
	}
	new_state.init()
	return new_state
}

func (st *State) getID() string {
	return st.id
}

func (st *State) setID(id string) {
	st.id = id
}

func (st *State) getEntity() IEntity {
	return st.entity
}

func (st *State) setEntity(e IEntity) {
	st.entity = e
}

func (st *State) setActiveSprite(spriteIndex int) *Sprite {
	return st.entity.setStateActiveSprite(st, spriteIndex)
}

func (st *State) getActiveSprite() *Sprite {
	return st.entity.getStateActiveSprite(st)
}

func (st *State) init() {
	log.Println("Default state init() handler called !")
}

func (st *State) enter() {
	log.Println("Default state enter() handler called !")
}

func (st *State) update() {
	log.Println("Default state update() handler called !")
}

func (st *State) leave() {
	log.Println("Default state leave() handler called !")
}
