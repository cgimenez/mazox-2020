package game

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"mazox/runtime"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	name       string
	size       float64
	DPI        float64
	image      runtime.Image
	advances   []int
	lineHeight int
	drawer     *font.Drawer
}

const (
	TEXT_NO_STYLE = 0
	TEXT_ALIGN_H  = 1
)

// backtick replaced with 0
var chars = ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_0abcdefghijklmnopqrstuvwxyz{|}~ ¡¢£¤¥¦§¨©ª«¬­®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþ`

func newFont(name string, size float64) Font {
	f := Font{name: name, size: size}
	f.DPI = 72
	return f
}

func (fnt *Font) load() {
	file := gm.rtm.ReadFile("assets/" + fnt.name + ".ttf")
	defer file.Close()

	ttf_data, _ := ioutil.ReadAll(file)
	f, err := truetype.Parse([]byte(ttf_data))
	if err != nil {
		log.Fatal(err)
	}

	fnt.drawer = &font.Drawer{
		Dst: fnt.image,
		Src: image.Black,
		Face: truetype.NewFace(f, &truetype.Options{
			Size: fnt.size,
			DPI:  fnt.DPI,
		}),
	}

	metrics := fnt.drawer.Face.Metrics()
	fnt.lineHeight = metrics.Ascent.Floor() + metrics.Descent.Floor() // Descent is negative

	fnt.drawer.Dst = fnt.image
	fnt.drawer.Dot = fixed.P(0, metrics.Ascent.Floor())
	fnt.image = image.NewRGBA(image.Rect(0, 0, 4096, fnt.lineHeight))
	draw.Draw(fnt.image, fnt.image.Bounds(), image.White, image.Point{}, draw.Src)
	fnt.drawer.DrawString(chars)

	var advance fixed.Int26_6
	prevRune := rune(-1)
	for i, runeValue := range chars {
		fnt.advances = append(fnt.advances, advance.Floor())
		if prevRune >= 0 {
			advance += fnt.drawer.Face.Kern(prevRune, runeValue)
		}
		adv, _ := fnt.drawer.Face.GlyphAdvance(runeValue)
		advance += adv
		if chars[i] == 32 {
			advance += 3
		}
		prevRune = runeValue
	}
}

func (fnt *Font) DrawString(s string, maxWidth int, minHeight int, padding int, flags int) runtime.Image {
	var height int
	var width int
	var x, y int

	// First we need to compute the bouding rect of the resulting image
	for i, runeValue := range s {
		if s[i] >= 32 {
			runeIndex := runeValue - 32
			runeWidth := fnt.advances[runeIndex+1] - fnt.advances[runeIndex]
			x += runeWidth
		}
		if (x >= maxWidth && s[i] == 32) || s[i] == 10 {
			if x >= width {
				width = x
			}
			x = 0
			height += fnt.lineHeight
		}
	}

	if width == 0 {
		width = x
	}
	height += fnt.lineHeight

	// Now we know, let's build the bitmap
	x = 0
	y = 0
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	for i, runeValue := range s {
		if s[i] >= 32 {
			runeIndex := runeValue - 32
			advance := fnt.advances[runeIndex]
			runeWidth := fnt.advances[runeIndex+1] - advance
			r := image.Rect(x, y, x+runeWidth, y+fnt.lineHeight)
			draw.Draw(img, r, fnt.image, image.Point{advance, 0}, draw.Src)
			x += runeWidth
		}
		if (x >= maxWidth && s[i] == 32) || s[i] == 10 {
			x = 0
			y += fnt.lineHeight
		}
	}
	return img
}
