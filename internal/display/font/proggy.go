package font

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

type Proggy struct {
	display    drivers.Displayer
	color      color.RGBA
	font       tinyfont.Fonter
	charWidth  int8
	charHeight int8
}

func NewProggy(display drivers.Displayer, color color.RGBA) *Proggy {
	if display == nil {
		return nil
	}

	return &Proggy{
		display:    display,
		color:      color,
		font:       &proggy.TinySZ8pt7b,
		charWidth:  6,
		charHeight: 6,
	}
}

func (f *Proggy) Print(x, y int16, text string) int16 {
	if f == nil {
		return 0
	}

	y += int16(f.charHeight)

	tinyfont.WriteLine(f.display, f.font, x, y, text, f.color)

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Proggy) CalcWidth(text string) int16 {
	if f == nil {
		return 0
	}

	if len(text) == 0 {
		return 0
	}

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Proggy) Width() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charWidth)
}

func (f *Proggy) Height() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charHeight)
}
