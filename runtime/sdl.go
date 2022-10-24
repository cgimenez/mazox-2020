// +build runtime_sdl

package runtime

import (
	"io"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type Runtime struct {
	_runtime
	window            *sdl.Window
	renderer          *sdl.Renderer
	rgba_image_buffer Image
	texture           *sdl.Texture
	rect              sdl.Rect
}

var keyConv = map[sdl.Keycode]int{
	sdl.K_ESCAPE: KEY_ESC,
	sdl.K_UP:     KEY_UP,
	sdl.K_DOWN:   KEY_DOWN,
	sdl.K_LEFT:   KEY_LEFT,
	sdl.K_RIGHT:  KEY_RIGHT,
	sdl.K_SPACE:  KEY_SPACE,
	sdl.K_d:      KEY_D,
}

func (r *Runtime) ReadFile(filename string) io.ReadCloser {
	log.Printf("Loading %s", filename)
	file, err := os.Open("www/" + filename)
	if err != nil {
		panic(err)
	}
	return file
}

func (r *Runtime) LoadAsset(filename string) io.ReadCloser {
	return r.ReadFile("/assets/" + filename)
}

func (r *Runtime) Init(width int, height int) {
	var err error

	r._init(width, height)

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	r.window, err = sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	r.renderer, err = sdl.CreateRenderer(r.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	r.texture, err = r.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, int32(r.width), int32(r.height))
	if err != nil {
		panic(err)
	}
	r.rect = sdl.Rect{X: 0, Y: 0, W: int32(r.width), H: int32(r.height)}
}

func (r *Runtime) Run(game IGame) {
	defer sdl.Quit()
	defer r.window.Destroy()

	log.Println("Enter sdl runtime mainloop")
	go func() {
		for {
			game.Update()
		}
	}()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					r.SetKeyState(keyConv[t.Keysym.Sym], true)
				}
				if t.Type == sdl.KEYUP {
					r.SetKeyState(keyConv[t.Keysym.Sym], false)
				}
			}
		}
		//r.renderer.Clear()
		game.Draw()

		err := r.texture.Update(&r.rect, r.rgba_image_buffer.Pix, r.width*4)
		if err != nil {
			panic(err)
		}
		err = r.renderer.Copy(r.texture, &r.rect, &r.rect)
		if err != nil {
			panic(err)
		}
		r.renderer.Present()
	}
}
