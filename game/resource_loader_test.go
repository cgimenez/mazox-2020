package game

import (
	"testing"
)

func TestResourceLoader(t *testing.T) {
	gm = new(Game)
	gm.Init(640, 320)

	im1 := gm.loader.loadImage("sprite_32.png")
	im2 := gm.loader.loadImage("sprite_32.png")

	if im1 != im2 {
		t.Errorf("Resources should be the same")
	}

	gm.loader.clear()
	if gm.loader.resources != nil {
		t.Errorf("Resources should be released")
	}
}
