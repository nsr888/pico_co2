package display

import (
	"fmt"
	"image/color"

	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayComfortIndexTinyFont(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	black := color.RGBA{1, 1, 1, 255}
	status := r.ComfortStatus()
	if r.ValidityError != "" {
		status = r.ValidityError
	}
	font9 := font.NewFreemonoRegular9(f.display, black)
	xPos := int16(0)
	yPos := int16(0)
	font9.Print(xPos, yPos, status)

	font12 := font.NewFreemonoRegular12(f.display, black)
	if r.ValidityError == "" {
		xPos = int16(0)
		yPos = int16(16)
		co2 := fmt.Sprintf("%d", r.CO2)
		font12.Print(xPos, yPos, co2)
	}

	hum := fmt.Sprintf("%.0f", r.Humidity)
	xPos = 128 - font12.CalcWidth(hum)
	yPos = int16(16)
	font12.Print(xPos, yPos, hum)

	temp := fmt.Sprintf("%.0f", r.Temperature)
	xPos = int16(128 - font12.CalcWidth(hum) - 10 - font12.CalcWidth(temp))
	yPos = int16(16)
	font12.Print(xPos, yPos, temp)

	f.display.Display()
}
