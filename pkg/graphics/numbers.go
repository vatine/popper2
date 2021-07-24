package graphics

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)


type segments struct {
	uprights map[float64]*ebiten.Image
	flats map[float64]*ebiten.Image
}

var segCache segments
var segCol = color.RGBA{R: 255, G: 200, B:0, A: 255}
var digitMap map[byte]byte

func init() {
	segCache.uprights = make(map[float64]*ebiten.Image)
	segCache.flats = make(map[float64]*ebiten.Image)
	digitMap = make(map[byte]byte)
	digitMap['0'] = 0x77
	digitMap['1'] = 0x24
	digitMap['2'] = 0x5d
	digitMap['3'] = 0x6d
	digitMap['4'] = 0x2e
	digitMap['5'] = 0x6b
	digitMap['6'] = 0x7b
	digitMap['7'] = 0x25
	digitMap['8'] = 0x7f
	digitMap['9'] = 0x6f
}


func drawUpright(img *ebiten.Image, w, h int) {
	middle := w/2

	for n := 0; n < middle; n++ {
		o1 := middle - n
		o2 := middle + n
		for q := middle; q < (h - middle); q++ {
			img.Set(o1, q, segCol)
			img.Set(o2, q, segCol)
		}
	}
}

func drawFlat(img *ebiten.Image, w, h int) {
	middle := h/2

	for n := 0; n < middle; n++ {
		o1 := middle - n
		o2 := middle + n
		for q := middle; q < (w - middle); q++ {
			img.Set(q, o1, segCol)
			img.Set(q, o2, segCol)
		}
	}
}

func getSegments(width float64) (*ebiten.Image, *ebiten.Image) {
	upright, ok := segCache.uprights[width]
	if !ok {
		hU := int(1.5 * width)
		wU := int(0.25 * width)
		hF := int(0.25 * width)
		wF := int(width)
		upright = ebiten.NewImage(wU, hU)
		flat := ebiten.NewImage(wF, hF)
		segCache.uprights[width] = upright
		segCache.flats[width] = flat

		drawUpright(upright, wU, hU)
		drawFlat(flat, wF, hF)
		
	}
	flat := segCache.flats[width]

	return upright, flat
}

// Render a seven-segment unit, based on the 1 bits being "lit up"
// Segments numbered as follows:
//
//    - 0 -
//   |     |
//   1     2
//   |     |
//    - 3 -
//   |     |
//   4     5
//   |     |
//    - 7 -
func sevenSegment(bg *ebiten.Image, width, x, y float64, segments byte) {
	upright, flat := getSegments(width)
	jig := width / 8

	if (segments & 0x01) != 0 {
		var g ebiten.GeoM
		g.Translate(x, y)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(flat, &opts)
	}
	if (segments & 0x02) != 0 {
		var g ebiten.GeoM
		g.Translate(x - jig, y)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(upright, &opts)
		
	}
	if (segments & 0x04) != 0 {
		var g ebiten.GeoM
		g.Translate(x+width+jig, y)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(upright, &opts)		
	}
	if (segments & 0x08) != 0 {
		var g ebiten.GeoM
		g.Translate(x, y + 1.5*width)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(flat, &opts)
	}
	if (segments & 0x10) != 0 {
		var g ebiten.GeoM
		g.Translate(x - jig, y + 1.5*width)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(upright, &opts)
	}
	if (segments & 0x20) != 0 {
		var g ebiten.GeoM
		g.Translate(x+width+jig, y + 1.5 * width)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(upright, &opts)
	}
	if (segments & 0x40) != 0 {
		var g ebiten.GeoM
		g.Translate(x, y + 3*width)
		opts := ebiten.DrawImageOptions{GeoM: g}
		bg.DrawImage(flat, &opts)
	}
}

func DrawNumber(bg *ebiten.Image, x, y, width float64, n int) {
	s := fmt.Sprintf("%d", n)
	for ix := 0; ix < len(s); ix++ {
		sevenSegment(bg, width, x + 2 * float64(ix)*width, y, digitMap[s[ix]])
	}
}
