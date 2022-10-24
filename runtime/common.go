package runtime

import (
	"image"
	"image/draw"
	"sync"
	"time"
)

type IGame interface {
	Update()
	Draw()
}

type Image = *image.RGBA

const (
	KEY_UP    = 38
	KEY_DOWN  = 40
	KEY_LEFT  = 37
	KEY_RIGHT = 39
	KEY_SPACE = 32
	KEY_ESC   = 27
	KEY_D     = 68
)

type _runtime struct {
	width, height   int
	keyStates       map[int]bool
	rgbaImageBuffer Image
	frame_count     int64
	frameCountStart time.Time
	startedAt       time.Time
	mux             sync.Mutex
}

func (r *Runtime) _init(width int, height int) {
	r.width = width
	r.height = height
	r.keyStates = make(map[int]bool)
	r.rgbaImageBuffer = image.NewRGBA(image.Rect(0, 0, width, height))
	r.frameCountStart = time.Now()
	r.startedAt = time.Now()
}

func (r *Runtime) HasKeyPressed(k int) bool {
	r.mux.Lock()
	res := r.keyStates[k]
	r.mux.Unlock()
	return res
}

func (r *Runtime) SetKeyState(k int, state bool) {
	r.mux.Lock()
	r.keyStates[k] = state
	r.mux.Unlock()
}

func (r *Runtime) DrawBitmap(img Image, x int, y int, dop bool) {
	var op draw.Op
	if dop {
		op = draw.Src
	} else {
		op = draw.Over
	}
	dp := image.Point{x, y}
	r.mux.Lock()
	draw.Draw(r.rgbaImageBuffer, image.Rectangle{dp, dp.Add(img.Bounds().Size())}, img, image.Point{0, 0}, op)
	r.mux.Unlock()
}
