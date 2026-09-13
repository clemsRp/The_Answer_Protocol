package drawer

import (
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

func (dr *Drawer) DrawWoodFrameAt(visualStart vars.Position, visualEnd vars.Position) {
	start := vars.Position{
		X: visualStart.X / 2,
		Y: visualStart.Y / 2,
	}
	end := vars.Position{
		X: visualEnd.X - start.X + 1,
		Y: visualEnd.Y - start.Y + 1,
	}
	dr.DrawWoodFrame(start, end)
}

func (dr *Drawer) DrawWoodFrame(start vars.Position, end vars.Position) {
	frame_start_x := 12
	frame_start_y := 0

	frame_width := int((end.X-start.X)/2 + 1)
	frame_height := int((end.Y-start.Y)/2 + 1)

	for x := range frame_width {
		for y := range frame_height {
			if x != 0 && y != 0 && x != frame_width-1 && y != frame_height-1 {
				continue
			}

			var ind_x int
			var ind_y int
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
				vars.WOOD_FRAME_TEXTURE,
				2*dr.app.Variables.Tileset_size*(float32(x)+start.X)-dr.app.Variables.Tileset_size,
				2*dr.app.Variables.Tileset_size*(float32(y)+start.Y)-dr.app.Variables.Tileset_size,
				float32(frame_start_x+ind_x), float32(frame_start_y+ind_y),
				1, 1, 2*dr.app.Variables.Zoom, 0,
			)
		}
	}
}
