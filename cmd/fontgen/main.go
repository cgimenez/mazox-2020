// Copyright 2010 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

// https://www.dafont.com/fr/minecraftia.font

package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type fontData struct {
	base     int
	size     int
	runes    string
	advances []int
}

var text = ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_0abcdefghijklmnopqrstuvwxyz{|}~ ¡¢£¤¥¦§¨©ª«¬­®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþ`

var fontname = flag.String("fontname", "Minecraftia-Regular", "font to convert")
var fontsize = flag.Float64("fontsize", 12, "font size")
var fontbase = flag.Int("fontbase", 80, "font base")

func main() {
	font_data := &fontData{base: *fontbase, size: int(*fontsize), runes: text}

	flag.Parse()
	fmt.Printf("Converting font %s\n", *fontname)

	ttfData, err := ioutil.ReadFile("fonts/" + *fontname + ".ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(ttfData)
	if err != nil {
		log.Fatal(err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, 3072, *fontbase))
	draw.Draw(dst, dst.Bounds(), image.White, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst: dst,
		Src: image.Black,
		Face: truetype.NewFace(f, &truetype.Options{
			Size: *fontsize,
			DPI:  72,
		}),
		Dot: fixed.P(0, *fontbase),
	}
	d.DrawString(text)

	var advance fixed.Int26_6
	prevRune := rune(-1)
	for _, runeValue := range text {
		font_data.advances = append(font_data.advances, advance.Floor())
		//fmt.Println(string(runeValue), advance.Floor())
		if prevRune >= 0 {
			advance += d.Face.Kern(prevRune, runeValue)
		}
		adv, _ := d.Face.GlyphAdvance(runeValue)
		advance += adv
		prevRune = runeValue
	}

	fmt.Println(d.Face.Metrics())
	//fmt.Println(d.Dot.X)

	out, err := os.Create(strings.ToLower(*fontname) + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	err = png.Encode(out, dst)
	if err != nil {
		log.Fatal(err)
	}

}
