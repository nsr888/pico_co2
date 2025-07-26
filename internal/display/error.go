package display

import (
	"pico_co2/internal/types"
)

func RenderError(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	if r != nil && r.Warning != "" {
		renderer.DrawLongText(0, 0, r.Warning)
	} else {
		renderer.DrawLongText(0, 0, "No error message available")
	}

	renderer.Display()
}
