package font

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
)

// TinyFontWrapper wraps TinyGo font implementations to implement FontPrinter interface
type TinyFontWrapper struct {
	display drivers.Displayer
	color   color.RGBA
	font    tinyfont.Fonter
}

// NewTinyFontWrapper creates a new wrapper for the given TinyGo font
func NewTinyFontWrapper(
	display drivers.Displayer,
	textColor color.RGBA,
	font tinyfont.Fonter,
) *TinyFontWrapper {
	if display == nil || font == nil {
		return nil
	}

	return &TinyFontWrapper{
		display: display,
		color:   textColor,
		font:    font,
	}
}

// Print draws the text at the specified x, y coordinates
func (fw *TinyFontWrapper) Print(x, y int16, text string) int16 {
	if fw == nil {
		return 0
	}

	y += fw.Height()

	tinyfont.WriteLine(fw.display, fw.font, x, y, text, fw.color)

	_, width := tinyfont.LineWidth(fw.font, text)
	return int16(width)
}

// CalcWidth calculates the width of the text using the font's metrics
func (fw *TinyFontWrapper) CalcWidth(text string) int16 {
	if fw == nil || fw.font == nil {
		return 0
	}

	_, width := tinyfont.LineWidth(fw.font, text)
	return int16(width)
}

// Height returns the height of the font
func (fw *TinyFontWrapper) Height() int16 {
	if fw == nil || fw.font == nil {
		return 0
	}

	// Get the height from a sample character (e.g., '0')
	glyph := fw.font.GetGlyph('0')
	if glyph == nil {
		return 0
	}

	return int16(glyph.Info().Height)
}

// Width returns the average width of a character in the font
func (fw *TinyFontWrapper) Width() int16 {
	if fw == nil || fw.font == nil {
		return 0
	}

	// Get the width from a sample character (e.g., '0')
	glyph := fw.font.GetGlyph('0')
	if glyph == nil {
		return 0
	}

	return int16(glyph.Info().Width)
}

// SetColor updates the text color
func (fw *TinyFontWrapper) SetColor(c color.RGBA) {
	if fw == nil {
		return
	}

	fw.color = c
}

func (fw *TinyFontWrapper) GetFont() tinyfont.Fonter {
	if fw == nil {
		return nil
	}

	return fw.font
}
