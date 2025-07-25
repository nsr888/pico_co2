package font

import (
	"github.com/Nondzu/ssd1306_font"
	"tinygo.org/x/drivers/ssd1306"
)

type Font11 struct {
	font       *ssd1306font.Display
	charWidth  int8
	charHeight int8
}

func NewFont11(display *ssd1306.Device) *Font11 {
	if display == nil {
		return nil
	}

	fontLib := ssd1306font.NewDisplay(*display)

	return &Font11{
		font:       &fontLib,
		charWidth:  11,
		charHeight: 16,
	}
}

func (f *Font11) Print(x, y int16, text string) int16 {
	if f == nil {
		return 0
	}

	f.font.Configure(ssd1306font.Config{
		FontType: ssd1306font.FONT_11x18,
	})
	f.font.XPos = x
	f.font.YPos = y
	f.font.PrintText(text)

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Font11) CalcWidth(text string) int16 {
	if f == nil {
		return 0
	}

	if len(text) == 0 {
		return 0
	}

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Font11) Width() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charWidth)
}

func (f *Font11) Height() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charHeight)
}
