package main

import (
	"mazox/game"
)

func main() {
	g := game.Game{}
	g.Init(640, 320)
	g.Start()
}
