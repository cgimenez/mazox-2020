package game

import (
	_ "image/png" // Warn - without this import decoder is not loaded and decoding panics
	"log"
	"mazox/runtime"
	"time"
)

type Game struct {
	delta         float32
	deltaTimeLast time.Time
	width         int
	height        int
	debug         bool
	activeZone    *Zone
	loader        *ResourceLoader
	rtm           *runtime.Runtime
	locale        string
	font          Font
}

var gm *Game

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func SNBH() {
	panic("Should not be here")
}

func (g *Game) Init(width int, height int) {
	gm = g
	g.locale = "en"
	g.width = width
	g.height = height
	g.debug = true
	g.loader = &ResourceLoader{}
	g.rtm = &runtime.Runtime{}
}

func (g *Game) Start() {
	log.Println("Init game")

	g.rtm.Init(g.width, g.height)

	g.deltaTimeLast = time.Now()

	g.font = newFont("minecraftia-regular", 14)
	g.font.load()

	g.activeZone = &Zone{num: 1}
	g.activeZone.load()
	g.activeZone.init()

	g.rtm.Run(g)
}

func (g *Game) Update() {
	if time.Since(g.deltaTimeLast).Milliseconds() > 10 {
		g.delta = float32(time.Since(g.deltaTimeLast).Milliseconds()) / 1000.0
		g.activeZone.update()
		if g.rtm.HasKeyPressed(runtime.KEY_D) {
			g.debug = !g.debug
		}
		g.deltaTimeLast = time.Now()
	}
}

func (g *Game) Draw() {
	g.activeZone.draw()
}
