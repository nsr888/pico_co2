package bar

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinydraw"

	"pico_co2/internal/display/font"
)

type TwoSideBar struct {
	display        drivers.Displayer
	radiusFilled   int16
	barSpacing     int16
	radiusEmpty    int16
	color          color.RGBA
	leftBarsCount  int16
	rightBarsCount int16
	font           font.FontPrinter
}

func NewTwoSideBar(
	fd drivers.Displayer,
	radiusFilled,
	barSpacing int16,
	leftBarsCount int16,
	rightBarsCount int16,
	font font.FontPrinter,
) *TwoSideBar {
	return &TwoSideBar{
		display:        fd,
		radiusFilled:   radiusFilled,
		barSpacing:     barSpacing,
		radiusEmpty:    1,
		color:          color.RGBA{1, 1, 1, 255},
		leftBarsCount:  leftBarsCount,
		rightBarsCount: rightBarsCount,
		font:           font,
	}
}

func (cfg *TwoSideBar) Draw(x, y, idxValue int16, label string) int16 {
	if cfg.display == nil {
		return 0
	}

	barY := y + cfg.radiusFilled + 1
	x = cfg.drawLeftBars(x, barY, idxValue)
	labelWidth := cfg.PrintText(x, y, label)
	x += labelWidth + cfg.radiusFilled + cfg.barSpacing
	x = cfg.drawRightBars(x, barY, idxValue)

	return x
}

func (cfg *TwoSideBar) drawLeftBars(x, y, idxValue int16) int16 {
	return cfg.drawBars(x, y, cfg.leftBarsCount, true,
		func(i int16) bool { return idxValue < 0 && i >= 3+idxValue },
	)
}

func (cfg *TwoSideBar) drawRightBars(x, y, idxValue int16) int16 {
	return cfg.drawBars(x, y, cfg.rightBarsCount, false,
		func(i int16) bool { return idxValue > 0 && i < idxValue },
	)
}

func (cfg *TwoSideBar) drawBars(
	x, y, count int16,
	isLeftSide bool,
	isFilled func(i int16) bool,
) int16 {
	if count == 0 {
		return x
	}

	barX := x + cfg.radiusFilled
	for i := int16(0); i < count; i++ {
		if isFilled(i) && isLeftSide {
			cfg.drawLargeDot(barX, y)
		}
		if isFilled(i) && !isLeftSide {
			cfg.drawLargeDot(barX, y)
		}
		if !isFilled(i) {
			cfg.drawSmallDot(barX, y)
		}
		barX += cfg.radiusFilled*2 + cfg.barSpacing
	}

	return barX
}

func (cfg *TwoSideBar) drawLargeDot(x, y int16) {
	tinydraw.FilledCircle(cfg.display, x, y, cfg.radiusFilled, cfg.color)
}

func (cfg *TwoSideBar) drawSmallDot(x, y int16) {
	tinydraw.FilledCircle(cfg.display, x, y, cfg.radiusEmpty, cfg.color)
}

func (cfg *TwoSideBar) drawTriangleUP(x, y int16) {
	tinydraw.FilledTriangle(
		cfg.display,
		x-cfg.radiusFilled, y+cfg.radiusFilled,
		x, y-cfg.radiusFilled,
		x+cfg.radiusFilled, y+cfg.radiusFilled,
		cfg.color,
	)
}

func (cfg *TwoSideBar) drawTriangleDOWN(x, y int16) {
	tinydraw.FilledTriangle(
		cfg.display,
		x-cfg.radiusFilled, y-cfg.radiusFilled,
		x, y+cfg.radiusFilled,
		x+cfg.radiusFilled, y-cfg.radiusFilled,
		cfg.color,
	)
}

func (cfg *TwoSideBar) PrintText(x, y int16, label string) int16 {
	fontHeight := cfg.font.Height()
	y = y + cfg.radiusFilled + 1 - fontHeight/2
	return cfg.font.Print(x, y, label)
}
