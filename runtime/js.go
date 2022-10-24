// +build runtime_js

// https://github.com/gohxs/folder2go

// https://github.com/mattn/golang-wasm-example/blob/master/main.go
// https://jlongster.com/Making-Sprite-based-Games-with-Canvas
// https://www.aaron-powell.com/posts/2019-02-06-golang-wasm-3-interacting-with-js-from-go/
// https://github.com/Yaoir
// https://github.com/markfarnan/go-canvas/tree/master/canvas

// Obsolète ?
// https://medium.zenika.com/go-1-11-webassembly-for-the-gophers-ae4bb8b1ee03
//
// https://www.google.com/search?q=%22js.FuncOf%22&oq=%22js.FuncOf%22&aqs=chrome..69i57j0l2.1663j0j7&sourceid=chrome&ie=UTF-8

package runtime

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"syscall/js"
	"time"
)

type Runtime struct {
	_runtime
	document    js.Value
	canvas      js.Value
	canvasImage js.Value
	ctx         js.Value
	imageBuffer js.Value
}

func check(e error) {
	if e != nil {
		//log.Fatal(e)
		panic(e)
	}
}

func (r *Runtime) ReadFile(filename string) io.ReadCloser {
	log.Printf("Loading %s", filename)
	resp, err := http.Get(filename)
	check(err)
	return resp.Body
}

func (r *Runtime) LoadAsset(filename string) io.ReadCloser {
	return r.ReadFile("/assets/" + filename)
}

func (r *Runtime) Init(width int, height int) {
	log.Println("Init JS runtime")
	r._init(width, height)

	r.document = js.Global().Get("document")
	r.canvas = r.document.Call("createElement", "canvas")
	r.ctx = r.canvas.Call("getContext", "2d")

	r.canvas.Set("width", js.ValueOf(width))
	r.canvas.Set("height", js.ValueOf(height))
	r.document.Get("body").Call("appendChild", r.canvas)

	r.canvasImage = r.ctx.Call("createImageData", width, height)
	r.imageBuffer = js.Global().Get("Uint8Array").New(len(r.rgbaImageBuffer.Pix))

	js.Global().Set("onload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Println("In load2")
		return nil
	}))

	kb_handler := func(this js.Value, args []js.Value) interface{} {
		ev_type := args[0].Get("type").String()
		if ev_type == "keydown" {
			k := args[0].Get("keyCode").Int()
			r.SetKeyState(k, true)
		}
		if ev_type == "keyup" {
			k := args[0].Get("keyCode").Int()
			r.SetKeyState(k, false)
		}
		return nil
	}
	r.document.Call("addEventListener", "keydown", js.FuncOf(kb_handler))
	r.document.Call("addEventListener", "keyup", js.FuncOf(kb_handler))
}

func (r *Runtime) Run(game IGame) {
	log.Println("Enter wasm runtime mainloop")
	//c = make(chan bool)

	var interval int64 = 16
	var _start_at time.Time
	var updateFn func(this js.Value, args []js.Value) interface{}
	updateFn = func(this js.Value, args []js.Value) interface{} {
		game.Update()
		_drift := time.Since(_start_at).Milliseconds() - interval
		//log.Println(_drift)
		_sinterval := strconv.Itoa(int(interval - _drift))
		//log.Println(_drift, _sinterval)
		js.Global().Call("setTimeout", js.FuncOf(updateFn), _sinterval)
		_start_at = time.Now()
		return nil
	}
	_start_at = time.Now()
	js.Global().Call("setTimeout", js.FuncOf(updateFn), strconv.Itoa(int(interval)))

	var renderFn func(this js.Value, args []js.Value) interface{}
	renderFn = func(this js.Value, args []js.Value) interface{} {
		game.Draw()
		js.CopyBytesToJS(r.imageBuffer, r.rgbaImageBuffer.Pix)
		r.canvasImage.Get("data").Call("set", r.imageBuffer)
		r.ctx.Call("putImageData", r.canvasImage, 0, 0)
		elapsed := time.Since(r.frameCountStart).Milliseconds()
		if elapsed >= 1000 {
			//log.Printf("%d", 1000*e.frame_count/time.Since(e.started_at).Milliseconds())
			r.frameCountStart = time.Now()
		}
		r.frame_count++
		js.Global().Call("requestAnimationFrame", js.FuncOf(renderFn))
		return nil
	}
	js.Global().Call("requestAnimationFrame", js.FuncOf(renderFn))
	select {} // ou <-c
}
