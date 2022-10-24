package game

import (
	"fmt"
	"image"
	"reflect"
)

const (
	DIRECTION_IDLE  = 0
	DIRECTION_DOWN  = 1
	DIRECTION_LEFT  = 2
	DIRECTION_RIGHT = 3
	DIRECTION_UP    = 4
)

type IEntity interface {
	load()
	init()
	update()
	draw()

	addSprite(*Sprite) int
	setActiveSprite(*Sprite)
	addState(string, IState, SpriteSlice)
	getState(string) IState
	setStateActiveSprite(IState, int) *Sprite
	getStateActiveSprite(IState) *Sprite

	addBubble(string)
	setActiveBubble(string)

	setLevel(ILevel)
	getLevel() ILevel
	setID(string)
	getID() string
	getGridRespect() bool
	setGridRespect(bool)

	setSpeed(float32)
	setFPS(int64)

	move(int) bool

	setPixPosition(int, int)
	getPixPosition() (int, int)
	setGridPosition(int, int)
	getGridPosition() (int, int)
	getRect() image.Rectangle
	getSize() (int, int)
}

type Entity struct {
	id           string
	level        ILevel
	states       []EntityStateEntry
	sprites      []*Sprite
	bubbles      map[string]*Bubble
	activeState  IState
	activeSprite *Sprite
	activeBubble *Bubble
	x, y         int
	col, row     int
	speed        float32
	gridRespect  bool
}

type EntityStateEntry struct {
	id             string
	state          IState
	spritesIndices []int   // Sprite index in Entity.sprites slice
	activeSprite   *Sprite // Sprite index in sprites_indices
}

type AnonymousEntity struct {
	Entity
}

type PlayerEntity struct {
	Entity
}

// ----------------------------------------------------------------------------
// FACTORY
// ----------------------------------------------------------------------------

func entityFactory(kind string, id string, l ILevel) IEntity {
	var entity IEntity

	switch kind {
	case "player":
		entity = &PlayerEntity{}
	case "dog":
		entity = &AnonymousEntity{}
	default:
		panic(fmt.Sprintf("Entity kind %s is unknown", kind))
	}
	entity.setID(id)
	entity.init()
	entity.setGridRespect(true)
	return entity
}

// ----------------------------------------------------------------------------
//
// ----------------------------------------------------------------------------

func (e *Entity) load() {
	SNBH()
}

func (e *Entity) init() {
	SNBH()
}

func (e *Entity) update() {
	if e.activeState == nil && len(e.states) > 0 {
		e.activeState = e.states[0].state
		e.activeState.enter()
	}
	e.activeState.update()

	if e.activeBubble != nil {
		if e.activeBubble.mustHide() {
			e.activeBubble = nil
		}
	}
}

func (e *Entity) draw() {
	e.activeSprite.x = e.x
	e.activeSprite.y = e.y
	if e.activeSprite != nil {
		e.activeSprite.draw()
	}
	if e.activeBubble != nil {
		e.drawBubble()
	}
}

func (e *Entity) move(direction int) bool {
	new_x := e.x
	new_y := e.y
	switch direction {
	case DIRECTION_LEFT:
		new_x -= int(e.speed * gm.delta)
	case DIRECTION_RIGHT:
		new_x += int(e.speed * gm.delta)
	case DIRECTION_DOWN:
		new_y += int(e.speed * gm.delta)
	case DIRECTION_UP:
		new_y -= int(e.speed * gm.delta)
	}
	success := e.level.entityCanMoveToPix(e, new_x, new_y)
	if success {
		e.setPixPosition(new_x, new_y)
	}
	return success
}

// ----------------------------------------------------------------------------
// Pix, grid pos, sizes and misc
// ----------------------------------------------------------------------------

func (e *Entity) _updateXY(x int, y int) {
	if e.activeSprite != nil { // A virtual entity might not have any sprite
		e.activeSprite.x = x
		e.activeSprite.y = y
	}
	e.x = x
	e.y = y
}

func (e *Entity) setPixPosition(x int, y int) {
	csz := e.level.getGridCellSize()
	e.col = x / csz
	e.row = y / csz
	e._updateXY(x, y)
}

func (e *Entity) getPixPosition() (int, int) {
	return e.x, e.y
}

func (e *Entity) setGridPosition(col int, row int) {
	e.col = col
	e.row = row
	s := e.level.getGridCellSize()
	e._updateXY(col*s, row*s)
}

func (e *Entity) getGridPosition() (int, int) {
	return e.col, e.row
}

func (e *Entity) getRect() image.Rectangle {
	return image.Rect(e.x, e.y, e.x+e.level.getGridCellSize(), e.y+e.level.getGridCellSize())
}

func (e *Entity) getSize() (int, int) {
	return e.activeSprite.getSize()
}

func (e *Entity) getGridRespect() bool {
	return e.gridRespect
}

func (e *Entity) setGridRespect(b bool) {
	e.gridRespect = b
}

// ----------------------------------------------------------------------------
// Simple getters and setters
// ----------------------------------------------------------------------------

func (e *Entity) setID(id string) {
	e.id = id
}

func (e *Entity) getID() string {
	return e.id
}

func (e *Entity) setLevel(l ILevel) {
	e.level = l
}

func (e *Entity) getLevel() ILevel {
	return e.level
}

func (e *Entity) setSpeed(speed float32) {
	e.speed = speed
}

// ----------------------------------------------------------------------------
// Bubbles handling
// ----------------------------------------------------------------------------

func (e *Entity) addBubble(k string) {
	if e.bubbles == nil { // uninitialized map
		e.bubbles = make(map[string]*Bubble, 5)
	}
	e.bubbles[k] = newBubble(k)
}

func (e *Entity) setActiveBubble(k string) {
	e.activeBubble = e.bubbles[k]
	e.activeBubble.display()
}

func (e *Entity) drawBubble() {
	gm.rtm.DrawBitmap(e.activeBubble.image, e.x+32, e.y, true)
}

// ----------------------------------------------------------------------------
// Sprites handling
// ----------------------------------------------------------------------------

func (e *Entity) setFPS(fps int64) {
	for sprite := range e.sprites {
		e.sprites[sprite].setFPS(fps)
	}
}

//
// set the entity active sprite from a state sprite
//
func (e *Entity) setStateActiveSprite(state IState, spriteIndex int) *Sprite {
	var result *Sprite

	found_state_entry := e.findStateEntry(state)
	result = e.sprites[found_state_entry.spritesIndices[spriteIndex]]
	e.activeSprite = result
	found_state_entry.activeSprite = result
	return result
}

func (e *Entity) getStateActiveSprite(state IState) *Sprite {
	return e.findStateEntry(state).activeSprite
}

func (e *Entity) addSprite(sprite *Sprite) int {
	e.activeSprite = sprite
	e.sprites = append(e.sprites, sprite)
	return len(e.sprites) - 1
}

func (e *Entity) setActiveSprite(sprite *Sprite) {
	e.activeSprite = sprite
}

// ----------------------------------------------------------------------------
// States handling
// ----------------------------------------------------------------------------

//
// Add a state to the entity
// first state sprite becomes (if any) the entity active sprite
//
func (e *Entity) addState(id string, state IState, sprites SpriteSlice) {
	var indices []int
	for i := range sprites {
		indices = append(indices, e.addSprite(sprites[i]))
	}
	state_entry := EntityStateEntry{id: id, state: state, spritesIndices: indices}
	if len(sprites) > 0 {
		state_entry.activeSprite = sprites[0]
	}
	e.states = append(e.states, state_entry)
	state.setEntity(e)
}

//
// Find the entry state corresponding with state
// Panic if not found
//
func (e *Entity) findStateEntry(state IState) *EntityStateEntry {
	var found_state_entry *EntityStateEntry

	for i := range e.states {
		if reflect.ValueOf(e.states[i].state).Pointer() == reflect.ValueOf(state).Pointer() {
			found_state_entry = &e.states[i]
			break
		}
	}
	if found_state_entry == nil {
		panic(fmt.Sprintf("State %v not found", state))
	}
	return found_state_entry
}

//
// find a state by its ID
//
func (e *Entity) getState(id string) IState {
	var found_state IState

	for i := range e.states {
		if e.states[i].id == id {
			found_state = e.states[i].state
			break
		}
	}
	if found_state == nil {
		panic(fmt.Sprintf("State with id %s not found", id))
	}
	return found_state
}

// ----------------------------------------------------------------------------
// To be moved somewhere else
// ----------------------------------------------------------------------------

func (e *PlayerEntity) load() {
}

func (e *PlayerEntity) init() {
	e.speed = 200
}

func (e *AnonymousEntity) load() {
}

func (e *AnonymousEntity) init() {
	e.speed = 200
}
