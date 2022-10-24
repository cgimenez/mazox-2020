package game

import "mazox/runtime"

type Level1 struct {
	Level
}

func (l *Level1) init() {
	player := entityFactory("player", "Chris", l)
	state := stateFactory("player_walk")

	images := gm.loader.loadTiledImages("L1/student-fmale-02.png", 3, 4, 32, 32)
	player.addState("", state, SpriteSlice{
		newSprite(SPRITE_ANIMATED, images[0:3]),
		newSprite(SPRITE_ANIMATED, images[3:6]),
		newSprite(SPRITE_ANIMATED, images[6:9]),
		newSprite(SPRITE_ANIMATED, images[9:12]),
	})
	l.addEntity(player)
	player.setFPS(5)
	player.setGridPosition(1, 1)

	for i := 0; i < 1; i++ {
		entity := entityFactory("dog", "Medor", l)
		entity.addBubble("l1.e1.1")
		state = stateFactory("random_walk")

		images = gm.loader.loadTiledImages("L1/dog-01-1.png", 3, 4, 32, 32)
		entity.addState("", state, SpriteSlice{
			newSprite(SPRITE_ANIMATED, images[0:3]),
			newSprite(SPRITE_ANIMATED, images[3:6]),
			newSprite(SPRITE_ANIMATED, images[6:9]),
			newSprite(SPRITE_ANIMATED, images[9:12]),
		})
		l.addEntity(entity)
		entity.setGridPosition(10, 1)
		/*
			for {
				spc := randomInt(0, 19)
				spr := randomInt(0, 9)
				if l.entityCanMoveToCell(entity, spc, spr) {
					entity.setGridPosition(spc, spr)
					break
				}
			} */
		// entity.setSpeed(float32(randomInt(30, 200)))
		entity.setSpeed(50)
		entity.setFPS(5)
	}
}

func (l *Level1) update() {
	l._update()
	if gm.rtm.HasKeyPressed(runtime.KEY_SPACE) {
		l.entities[1].setActiveBubble("l1.e1.1")
	}
}
