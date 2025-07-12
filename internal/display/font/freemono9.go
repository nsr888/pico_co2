package font

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

type Freemono9 struct {
	display    drivers.Displayer
	color      color.RGBA
	font       tinyfont.Fonter
	charWidth  int8
	charHeight int8
}

func NewFreemonoRegular9(display drivers.Displayer, color color.RGBA) *Freemono9 {
	if display == nil {
		return nil
	}

	return &Freemono9{
		display:    display,
		color:      color,
		font:       &freemono.Regular9pt7b,
		charWidth:  11,
		charHeight: 11,
	}
}

func (f *Freemono9) Print(x, y int16, text string) int16 {
	if f == nil {
		return 0
	}

	y += int16(f.charHeight)

	tinyfont.WriteLine(f.display, f.font, x, y, text, f.color)

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Freemono9) CalcWidth(text string) int16 {
	if f == nil {
		return 0
	}

	if len(text) == 0 {
		return 0
	}

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Freemono9) Width() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charWidth)
}

func (f *Freemono9) Height() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charHeight)
}
