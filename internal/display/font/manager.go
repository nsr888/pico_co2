package font

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/freesans"
)

// FontType represents different font categories available in the system
type FontType int

const (
	FremonoRegular18 FontType = iota
	FreemonoBold18
	FreemonoRegular12
	FreemonoBold12
	FreesansRegular12
	FreesansBold12
	FreesansRegular9
	FreesansBold9
	ProggySZ8
)

// FontRegistry manages font instances and provides centralized access
type FontRegistry struct {
	display drivers.Displayer
	color   color.RGBA
	fonts   map[FontType]FontPrinter
}

// NewFontRegistry creates a new font registry with the given display and color
func NewFontRegistry(
	display drivers.Displayer,
	textColor color.RGBA,
) *FontRegistry {
	if display == nil {
		return nil
	}

	registry := &FontRegistry{
		display: display,
		color:   textColor,
		fonts:   make(map[FontType]FontPrinter),
	}

	registry.initializeFonts()
	return registry
}

// GetFont returns a font printer for the specified font type
func (fr *FontRegistry) GetFont(fontType FontType) FontPrinter {
	if fr == nil {
		return nil
	}

	font, exists := fr.fonts[fontType]
	if !exists {
		// Return small font as fallback
		return fr.fonts[ProggySZ8]
	}

	return font
}

// SetColor updates the text color for all fonts
func (fr *FontRegistry) SetColor(c color.RGBA) {
	if fr == nil {
		return
	}

	fr.color = c

	// Update all fonts with new color
	for fontType := range fr.fonts {
		fr.fonts[fontType] = fr.createFont(fontType)
	}
}

// initializeFonts creates and registers all available fonts
func (fr *FontRegistry) initializeFonts() {
	fontTypes := []FontType{
		FremonoRegular18,
		FreemonoBold18,
		FreemonoRegular12,
		FreemonoBold12,
		FreesansRegular12,
		FreesansBold12,
		FreesansRegular9,
		FreesansBold9,
		ProggySZ8,
	}

	for _, fontType := range fontTypes {
		fr.fonts[fontType] = fr.createFont(fontType)
	}
}

// createFont creates a font instance for the specified type
func (fr *FontRegistry) createFont(fontType FontType) FontPrinter {
	switch fontType {
	case FremonoRegular18:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freemono.Regular18pt7b,
		}
	case FreemonoBold18:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freemono.Bold18pt7b,
		}
	case FreemonoRegular12:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freemono.Regular12pt7b,
		}
	case FreemonoBold12:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freemono.Bold12pt7b,
		}
	case FreesansRegular12:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freesans.Regular12pt7b,
		}
	case FreesansBold12:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freesans.Bold12pt7b,
		}
	case FreesansRegular9:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freesans.Regular9pt7b,
		}
	case FreesansBold9:
		return &TinyFontWrapper{
			display: fr.display,
			color:   fr.color,
			font:    &freesans.Bold9pt7b,
		}
	case ProggySZ8:
		// Use existing Proggy implementation for small text
		return NewProggy(fr.display, fr.color)
	default:
		return NewProggy(fr.display, fr.color)
	}
}
