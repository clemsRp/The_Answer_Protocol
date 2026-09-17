package drawer

import (
	"math"
	"time"

	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawImage(texture_name string, posX, posY, indX, indY, ratioX, ratioY, zoom, rotation float32) {
	texture := (*dr.app.Textures)[texture_name]

	sourceW := ratioX * float32(vars.FRAME_WIDTH)
	sourceH := ratioY * float32(vars.FRAME_HEIGHT)

	sourceRec := rl.NewRectangle(
		indX*float32(vars.FRAME_WIDTH), indY*float32(vars.FRAME_HEIGHT),
		sourceW, sourceH,
	)

	destW := float32(vars.FRAME_WIDTH) * zoom * ratioX
	destH := float32(vars.FRAME_HEIGHT) * zoom * ratioY

	destRec := rl.NewRectangle(posX+(destW/2), posY+(destH/2), destW, destH)

	origin := rl.NewVector2(destW/2, destH/2)

	rl.DrawTexturePro(texture, sourceRec, destRec, origin, rotation, rl.White)
}

func (dr *Drawer) DrawWoodFrameAt(visualStart vars.Position, visualEnd vars.Position, zoom float32, darkness int, empty bool) {
	e := float32(zoom)

	start := vars.Position{
		X: visualStart.X / e,
		Y: visualStart.Y / e,
	}

	frame_width := int(math.Round(float64(visualEnd.X-visualStart.X)/float64(e))) + 1
	frame_height := int(math.Round(float64(visualEnd.Y-visualStart.Y)/float64(e))) + 1

	dr.DrawWoodFrame(start, frame_width, frame_height, zoom, darkness, empty)
}

func (dr *Drawer) DrawWoodFrame(start vars.Position, frame_width, frame_height int, zoom float32, darkness int, empty bool) {
	frame_start_x := 12
	frame_start_y := 0

	e := float32(zoom)
	cell := e * dr.app.Variables.Tileset_size

	for x := range frame_width {
		for y := range frame_height {
			if x != 0 && y != 0 && x != frame_width-1 && y != frame_height-1 && empty {
				continue
			}

			var ind_x, ind_y int
			switch x {
			case 0:
				ind_x = 0
			case frame_width - 1:
				ind_x = 2
			default:
				ind_x = 1
			}
			switch y {
			case 0:
				ind_y = 0
			case frame_height - 1:
				ind_y = 2
			default:
				ind_y = 1
			}

			dr.DrawImage(
				vars.UI_SPRITE_TEXTURE,
				cell*(float32(x)+start.X)-cell/2,
				cell*(float32(y)+start.Y)-cell/2,
				float32(frame_start_x+ind_x), float32(frame_start_y+ind_y+darkness*3),
				1, 1, e*dr.app.Variables.Zoom, 0,
			)
		}
	}
}

func (dr *Drawer) DrawBorderedInput(x, y, width, height, border int32, border_color, input_color rl.Color) {
	// Border
	rl.DrawRectangle(x-border, y, width+(2*border), height, border_color)
	rl.DrawRectangle(x, y-border, width, height+(2*border), border_color)

	// Input
	rl.DrawRectangle(x, y+border, width, height-(2*border), input_color)
	rl.DrawRectangle(x+border, y, width-(2*border), height, input_color)
}

func (dr *Drawer) DrawCursor(x, y, width, height int32, text string, color rl.Color) {
	cursor_duration := 750
	elapsed := int(time.Since(dr.app.Variables.StartTime).Milliseconds())
	too_late := time.Since(dr.app.Variables.Player.LastTimeTyped).Milliseconds() > 300
	text_size := rl.MeasureText(text, int32(dr.app.Variables.FontSize))

	if (elapsed/cursor_duration)%2 == 0 && too_late {
		rl.DrawRectangle(
			x+text_size, y, width, height, color,
		)
	}
}
