package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGroupPanel() {
	group_width := float32(vars.GROUP_WIDTH)
	group_height := float32(vars.GROUP_HEIGHT)
	group_start_x := float32(vars.GROUP_START_X)
	group_start_y := float32(vars.GROUP_START_Y)

	// Draw frame
	dr.DrawWoodFrameAt(
		vars.Position{X: group_start_x, Y: group_start_y},
		vars.Position{X: group_start_x + group_width, Y: group_start_y + group_height},
		1, 0, false,
	)

	// Draw group management
	rl.DrawText(
		"Group",
		int32((group_start_x+0.7)*dr.app.Variables.Tileset_size),
		int32((group_start_y+0.7)*dr.app.Variables.Tileset_size),
		int32(1.5*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	dr.DrawGroupOptions()
}
