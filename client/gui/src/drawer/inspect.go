package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawInspectPanel() {
	inspect_width := float32(vars.INSPECT_WIDTH)
	inspect_height := float32(vars.INSPECT_HEIGHT)
	inspect_start_x := float32(vars.INSPECT_START_X)
	inspect_start_y := float32(vars.INSPECT_START_Y)

	// Draw frame
	dr.DrawWoodFrameAt(
		vars.Position{X: inspect_start_x, Y: inspect_start_y},
		vars.Position{X: inspect_start_x + inspect_width, Y: inspect_start_y + inspect_height},
		1, 0, false,
	)

	tile := dr.app.Variables.Tileset_size

	// Draw inspect management
	rl.DrawText(
		"Inspect",
		int32((inspect_start_x+0.7)*tile),
		int32((inspect_start_y+0.7)*tile),
		int32(1.5*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	// Calculate center_text
	inspect_title_width := vars.INSPECT_WIDTH * dr.app.Variables.Tileset_size
	title_width := rl.MeasureText(
		"Inspect", int32(1.5*dr.app.Variables.FontSize),
	)
	payload_width := rl.MeasureText(
		"Self", int32(dr.app.Variables.FontSize),
	)
	center_text := (inspect_title_width + float32(title_width-payload_width)) / 2

	dr.DrawInspectButtons()

	// Draw button text
	rl.DrawText(
		"Self",
		int32((vars.INSPECT_START_X)*dr.app.Variables.Tileset_size+center_text),
		int32((vars.INSPECT_START_Y+0.85)*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)

	// Draw inspect datas
	rl.DrawText(
		dr.app.Variables.PanelsVariables.Inspect.Datas,
		int32((inspect_start_x+0.7)*tile),
		int32((inspect_start_y+1)*tile)+int32(2*dr.app.Variables.FontSize),
		int32(0.85*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)
}
