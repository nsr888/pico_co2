package font

import (
	"github.com/Nondzu/ssd1306_font"
	"tinygo.org/x/drivers/ssd1306"
)

type Font7 struct {
	font       *ssd1306font.Display
	charWidth  int8
	charHeight int8
}

func NewFont7(display *ssd1306.Device) *Font7 {
	if display == nil {
		return nil
	}

	fontLib := ssd1306font.NewDisplay(*display)

	return &Font7{
		font:       &fontLib,
		charWidth:  7,
		charHeight: 8,
	}
}

func (f *Font7) Print(x, y int16, text string) int16 {
	if f == nil {
		return 0
	}

	f.font.Configure(ssd1306font.Config{
		FontType: ssd1306font.FONT_7x10,
	})
	f.font.XPos = x
	f.font.YPos = y
	f.font.PrintText(text)

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Font7) CalcWidth(text string) int16 {
	if f == nil {
		return 0
	}

	if len(text) == 0 {
		return 0
	}

	return int16(len(text)) * int16(f.charWidth)
}

func (f *Font7) Width() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charWidth)
}

func (f *Font7) Height() int16 {
	if f == nil {
		return 0
	}

	return int16(f.charHeight)
}
