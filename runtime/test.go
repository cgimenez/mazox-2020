// +build runtime_test

package runtime

import (
	"io"
	"log"
	"os"
)

type Runtime struct {
	_runtime
}

func (r *Runtime) Init(width int, height int)                   {}
func (r *Runtime) Run(game IGame)                               {}
func (r *Runtime) DrawBitmap(img Image, x int, y int, dop bool) {}
func (r *Runtime) HasKeyPressed(k int) bool                     { return false }
func (r *Runtime) LoadAsset(filename string) io.ReadCloser {
	return r.LoadFile(filename)
}
func (r *Runtime) LoadFile(filename string) io.ReadCloser {
	filename = "../fixtures/" + filename
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
