package game

import (
	"image"
	_ "log"
	"mazox/runtime"
	"time"
)

const (
	SPRITE_STATIC = iota
	SPRITE_ANIMATED
)

var sGID int

type SpriteSlice []*Sprite

type Sprite struct {
	id        int
	kind      int
	images    []*image.RGBA
	x         int
	y         int
	width     int
	height    int
	curFrame  int
	numFrames int
	fps       int64
	timestamp time.Time
	visible   bool
	stopped   bool
}

func loadSprite(kind int, filename string, numFrames int) *Sprite {
	if numFrames == 1 {
		var sl []runtime.Image
		sl = append(sl, gm.loader.loadImage(filename))
		return newSprite(kind, sl)
	} else {
		return newSprite(kind, gm.loader.loadImageSequence(filename, numFrames))
	}
}

func newSprite(kind int, images []runtime.Image) *Sprite {
	sGID++
	s := &Sprite{id: sGID, numFrames: len(images)}
	s.images = images
	s.timestamp = time.Now()
	s.visible = true
	s.stopped = false
	s.width = s.images[0].Bounds().Max.X
	s.height = s.images[0].Bounds().Max.Y
	return s
}

func (s *Sprite) getSize() (int, int) {
	return s.width, s.height
}

func (s *Sprite) setFPS(fps int64) {
	s.fps = 1000 / fps
}

func (s *Sprite) stop() {
	s.stopped = true
}

func (s *Sprite) play() {
	s.stopped = false
}

func (s *Sprite) rewind() {
	s.curFrame = 0
}

func (s *Sprite) draw() {
	if !s.visible {
		return
	}
	gm.rtm.DrawBitmap(s.images[s.curFrame], s.x, s.y, false)
	if s.kind == SPRITE_ANIMATED && !s.stopped {
		elapsed := time.Since(s.timestamp).Milliseconds()
		if elapsed >= s.fps {
			s.curFrame++
			if s.curFrame >= s.numFrames {
				s.curFrame = 0
			}
			s.timestamp = time.Now()
		}
	}
}
