package game

type Zone struct {
	num         int
	levels      []ILevel
	activeLevel ILevel
}

func (z *Zone) load() {
	z.levels = nil

	switch z.num {
	case 1:
		level := &Level1{Level{num: 1, zone: z}}
		level.load()
		z.levels = append(z.levels, level)
	}
}

func (z *Zone) init() {
	for level := range z.levels {
		z.levels[level].init()
	}
	z.activeLevel = z.levels[0]
}

func (z *Zone) update() {
	for level := range z.levels {
		z.levels[level].update()
	}
}

func (z *Zone) draw() {
	z.activeLevel.draw()
}
