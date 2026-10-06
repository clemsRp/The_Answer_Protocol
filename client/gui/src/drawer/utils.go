package drawer

import (
	"math"
	"strings"
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
	e := float32(zoom)
	cell := e * dr.app.Variables.Tileset_size

	startX := cell * start.X
	startY := cell * start.Y

	dr.DrawRealWoodFrame(
		startX, startY,
		frame_width, frame_height,
		zoom, darkness, empty,
	)
}

func (dr *Drawer) DrawRealWoodFrame(startX, startY float32, frame_width, frame_height int, zoom float32, darkness int, empty bool) {
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
				startX+cell*float32(x)-cell/2,
				startY+cell*float32(y)-cell/2,
				float32(frame_start_x+ind_x), float32(frame_start_y+ind_y+darkness*3),
				1, 1, e*dr.app.Variables.Zoom, 0,
			)
		}
	}
}

func (dr *Drawer) DrawWoodFramePx(x, y, w, h, corner float32, darkness int, empty bool) {
	texture := (*dr.app.Textures)[vars.UI_SPRITE_TEXTURE]

	nx := max(2, int(math.Round(float64(w/corner))))
	ny := max(2, int(math.Round(float64(h/corner))))
	cw := w / float32(nx)
	ch := h / float32(ny)

	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			if empty && i != 0 && j != 0 && i != nx-1 && j != ny-1 {
				continue
			}

			indX, indY := 1, 1
			if i == 0 {
				indX = 0
			} else if i == nx-1 {
				indX = 2
			}
			if j == 0 {
				indY = 0
			} else if j == ny-1 {
				indY = 2
			}

			src := rl.NewRectangle(
				float32((12+indX)*vars.FRAME_WIDTH),
				float32((indY+darkness*3)*vars.FRAME_HEIGHT),
				float32(vars.FRAME_WIDTH), float32(vars.FRAME_HEIGHT),
			)
			dst := rl.NewRectangle(x+cw*float32(i), y+ch*float32(j), cw+1, ch+1)
			rl.DrawTexturePro(texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)
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

func (dr *Drawer) LimitString(text string, limit int) string {
	res := text
	if len(res) > limit {
		res = string([]rune(res)[:limit-2])
		res += "..."
	}

	return res
}

func (dr *Drawer) WrapText(text string, maxLineLen int) (string, int) {
	if maxLineLen <= 0 {
		return text, 1
	}

	// Split sentence
	words := strings.Fields(text)
	if len(words) == 0 {
		return "", 0
	}

	var builder strings.Builder
	nb_line := 1
	currentLineLen := 0

	for i, word := range words {
		// Handle first word
		if i == 0 {
			builder.WriteString(word)
			currentLineLen = len(word)
			continue
		}

		// Add \n if needed
		if currentLineLen+1+len(word) > maxLineLen {
			builder.WriteString("\n")
			builder.WriteString(word)
			currentLineLen = len(word)
			nb_line++

		} else {
			// Add word normally
			builder.WriteString(" ")
			builder.WriteString(word)
			currentLineLen += 1 + len(word)
		}
	}

	return builder.String(), nb_line
}
