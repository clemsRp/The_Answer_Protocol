package drawer

import (
	"strconv"
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGroupPanel() {
	group_width := float32(vars.GROUP_WIDTH)
	group_height := float32(vars.GROUP_HEIGHT)
	group_start_x := float32(vars.GROUP_START_X)
	group_start_y := float32(vars.GROUP_START_Y)
	group := dr.app.Variables.PanelsVariables.Group

	// Draw frame
	dr.DrawWoodFrameAt(
		vars.Position{X: group_start_x, Y: group_start_y},
		vars.Position{X: group_start_x + group_width, Y: group_start_y + group_height},
		1, 0, false,
	)

	tile := dr.app.Variables.Tileset_size
	font := dr.app.Variables.FontSize

	// Draw group management
	rl.DrawText(
		"Group",
		int32((group_start_x+0.7)*tile),
		int32((group_start_y+0.7)*tile),
		int32(1.5*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	dr.DrawGroupButtons()
	dr.DrawCommonButton()
	dr.DrawGroupOptions()

	if !group.InGroup {
		return
	}

	dr.DrawGroupDatas(group, group_start_x, group_start_y, group_width, tile, font)
}

func (dr *Drawer) DrawCommonButton() {
	// Get payload
	var payload string
	if dr.app.Variables.PanelsVariables.Group.InGroup {
		payload = "Leave"
	} else {
		payload = "Create"
	}

	// Calculate center_text
	group_width := vars.GROUP_WIDTH * dr.app.Variables.Tileset_size
	title_width := rl.MeasureText(
		"Group", int32(1.5*dr.app.Variables.FontSize),
	)
	payload_width := rl.MeasureText(
		payload, int32(dr.app.Variables.FontSize),
	)
	center_text := (group_width + float32(title_width-payload_width)) / 2

	// Draw button text
	rl.DrawText(
		payload,
		int32((vars.GROUP_START_X)*dr.app.Variables.Tileset_size+center_text),
		int32((vars.GROUP_START_Y+0.85)*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)
}

func (dr *Drawer) DrawGroupDatas(group *vars.GroupPanel, group_start_x, group_start_y, group_width, tile, font float32) {
	infos := map[string]string{
		"Leader":  group.Leader,
		"Members": strconv.Itoa(len(group.Grouped)),
	}
	infos_keys := []string{"Leader", "Members"}

	if group.InGroup {
		num, _ := strconv.Atoi(infos["Members"])
		infos["Members"] = strconv.Itoa(num + 1)
	}

	middle := (group_start_x + group_width/2) * tile
	index := 0
	for _, info_key := range infos_keys {
		info_value := infos[info_key]
		// Measure text
		align_info := rl.MeasureText(
			info_key+"   ", int32(font),
		)

		// Draw info key
		rl.DrawText(
			info_key,
			int32(middle)-align_info,
			int32((group_start_y+2)*float32(tile)+float32(index*int(font))),
			int32(font), dr.app.Colors["pseudo_text"],
		)

		// Draw info value
		rl.DrawText(
			dr.LimitString(info_value, 10),
			int32(middle),
			int32((group_start_y+2)*float32(tile)+float32(index*int(font))),
			int32(font), dr.app.Colors["group_text"],
		)

		index++
	}
}
