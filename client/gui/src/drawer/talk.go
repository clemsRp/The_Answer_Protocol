package drawer

import (
	vars "tap/client/gui/src/variables"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawTalk displays the npc dialogue box centered near the bottom of the screen
func (dr *Drawer) DrawTalk() {
	talk := dr.app.Variables.PanelsVariables.Talk.Talking
	if talk == nil {
		return
	}

	box_start_x := 6.0
	box_start_y := 12.0
	box_width := 20.0
	box_height := 4.5

	dr.DrawWoodFrameAt(
		vars.Position{X: float32(box_start_x), Y: float32(box_start_y)},
		vars.Position{X: float32(box_start_x + box_width), Y: float32(box_start_y + box_height)},
		1, 0, false,
	)

	// Show a placeholder while waiting for the server answer
	const msPerChar = 30
	var text string

	// Animate text
	if !dr.app.Variables.PanelsVariables.Talk.Finished {
		text = "..."
		elapsed := time.Since(talk.Start)

		if talk.Result != "" {
			text = talk.Result
			elapsed = time.Since(talk.ResultStart)
			runes := []rune(text)
			index := int(elapsed.Milliseconds()) / msPerChar
			index = max(0, min(len(runes), index))
			text = string(runes[:index])

			// Check animation end
			if index >= len(talk.Result)-1 {
				dr.app.Variables.PanelsVariables.Talk.Finished = true
			}
		}

		// Don't animate text
	} else {
		text = talk.Result
	}

	rl.DrawText(
		text,
		int32(float32(box_start_x+0.75)*dr.app.Variables.Tileset_size),
		int32(float32(box_start_y+0.75)*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	if talk.Result != "" {
		rl.DrawText(
			"click to continue",
			int32(float32(box_start_x+box_width-6)*dr.app.Variables.Tileset_size),
			int32(float32(box_start_y+box_height-0.9)*dr.app.Variables.Tileset_size),
			int32(dr.app.Variables.FontSize*0.6),
			dr.app.Colors["pseudo_text"],
		)
	}

	dr.DrawTalkButtons()
}
