package game

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"mazox/runtime"
	"strings"
)

type ResourceBuffer = runtime.Image

type ResourceLoader struct {
	resources map[string]ResourceBuffer
}

func (loader *ResourceLoader) loadImage(filename string) ResourceBuffer {
	if loader.resources == nil {
		loader.resources = make(map[string]ResourceBuffer)
	}
	for k := range loader.resources {
		if k == filename {
			return loader.resources[k]
		}
	}

	reader := gm.rtm.LoadAsset(filename)
	defer reader.Close()
	decoded_img, _, err := image.Decode(reader)
	check(err)
	img := image.NewRGBA(decoded_img.Bounds())
	draw.Draw(img, img.Bounds(), decoded_img, decoded_img.Bounds().Min, draw.Src)

	loader.resources[filename] = img
	return img
}

func (loader *ResourceLoader) loadImageSequence(filename string, num_frames int) []runtime.Image {
	var images []runtime.Image

	basename := strings.Split(filename, ".")[0]
	extension := strings.Split(filename, ".")[1]
	for i := 0; i < num_frames; i++ {
		images = append(images, loader.loadImage(fmt.Sprintf("%s_f%d.%s", basename, i, extension)))
	}
	return images
}

func (loader *ResourceLoader) loadTiledImages(filename string, cols int, rows int, width int, height int) []runtime.Image {
	var images []runtime.Image

	tileImage := loader.loadImage(filename)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			img := image.NewRGBA(image.Rect(0, 0, width, height))
			draw.Draw(img,
				image.Rect(0, 0, width, height),
				tileImage, image.Point{col * width, row * height},
				draw.Src,
			)
			images = append(images, img)
		}
	}
	return images
}

func (loader *ResourceLoader) clear() {
	loader.resources = nil
}
