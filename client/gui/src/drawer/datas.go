package drawer

import (
	"strconv"
	vars "tap/client/gui/src/variables"
	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawDatasPanel() {
	datas_start_x := float32(vars.DATAS_START_X)
	datas_start_y := float32(vars.DATAS_START_Y)
	datas_width := float32(vars.DATAS_WIDTH)
	datas_height := float32(vars.DATAS_HEIGHT)

	// Draw frame
	dr.DrawWoodFrameAt(
		vars.Position{X: datas_start_x, Y: datas_start_y},
		vars.Position{X: datas_start_x + datas_width, Y: datas_start_y + datas_height},
		1, 0, false,
	)

	tile := dr.app.Variables.Tileset_size
	font := dr.app.Variables.FontSize

	// Draw datas management
	rl.DrawText(
		"Datas",
		int32((datas_start_x+0.7)*tile),
		int32((datas_start_y+0.7)*tile),
		int32(1.5*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	dr.DrawDatasButtons()

	player_datas_width := datas_width * tile

	// Calculate center_text
	title_width := rl.MeasureText(
		"Group", int32(1.5*dr.app.Variables.FontSize),
	)
	payload_width := rl.MeasureText(
		"Update", int32(dr.app.Variables.FontSize),
	)
	center_text := (player_datas_width + float32(title_width-payload_width)) / 2

	// Draw button text
	rl.DrawText(
		"Update",
		int32((vars.GROUP_START_X)*dr.app.Variables.Tileset_size+center_text),
		int32((vars.GROUP_START_Y+0.85)*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)

	datas := dr.app.Variables.PanelsVariables.Datas

	dr.DrawPlayersDatas(datas, datas_start_x, datas_start_y, datas_width, tile, font)
	dr.DrawQuestsDatas(datas, datas_start_x, datas_start_y, datas_width, tile, font)
}

func (dr *Drawer) DrawPlayersDatas(datas *vars.DatasPanel, datas_start_x, datas_start_y, datas_width, tile, font float32) {
	infos := map[string]string{
		"Room":  strconv.Itoa(len(*dr.app.Variables.RemotePlayers) + 1),
		"Total": strconv.Itoa(datas.NbServerPlayers),
	}
	infos_keys := []string{"Room", "Total"}

	// Draw datas management
	rl.DrawText(
		"Players",
		int32((datas_start_x+0.7)*tile),
		int32((datas_start_y+1.75)*float32(tile)),
		int32(1.2*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	middle := (datas_start_x + datas_width/2) * tile
	index := 1
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
			int32((datas_start_y+2)*float32(tile)+float32(index*int(font))),
			int32(font), dr.app.Colors["pseudo_text"],
		)

		// Draw info value
		rl.DrawText(
			dr.LimitString(info_value, 10),
			int32(middle),
			int32((datas_start_y+2)*float32(tile)+float32(index*int(font))),
			int32(font), dr.app.Colors["group_text"],
		)

		index++
	}
}

func (dr *Drawer) DrawQuestsDatas(datas *vars.DatasPanel, datas_start_x, datas_start_y, datas_width, tile, font float32) {
	rl.DrawText(
		"Quests",
		int32((datas_start_x+0.7)*tile),
		int32((datas_start_y+4)*float32(tile)),
		int32(1.2*dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	quests := dr.app.Variables.PanelsVariables.Quests

	if len(*quests) == 0 {
		rl.DrawText(
			"None",
			int32((datas_start_x+4)*tile),
			int32((datas_start_y+4)*float32(tile)),
			int32(1.2*dr.app.Variables.FontSize),
			dr.app.Colors["group_text"],
		)

		return
	}

	actives := get_active_quests(*quests)
	finished := get_finished_quests(*quests)

	indY := int32((datas_start_y + 5) * float32(tile))

	if len(actives) != 0 {
		dr.DrawQuests("Actives", actives, indY, datas_start_x, tile, font, rl.Green)
		indY += int32(len(actives)+2) * int32(font)
	}

	if len(finished) != 0 {
		dr.DrawQuests("Finished", finished, indY, datas_start_x, tile, font, rl.Red)
		indY += int32(len(finished)) * int32(font)
	}
}

func (dr *Drawer) DrawQuests(title string, quests []pr.TrackedQuestData, start_y int32, datas_start_x, tile, font float32, color rl.Color) {
	indX := int32((datas_start_x + 0.7) * tile)

	// Draw Title
	rl.DrawText(
		title,
		indX,
		start_y,
		int32(dr.app.Variables.FontSize),
		color,
	)

	index := 0
	for _, quest := range quests {
		index++
		// Draw Quest ID
		rl.DrawText(
			quest.Id,
			indX+int32(0.33*tile),
			start_y+int32(index)*int32(font),
			int32(font), dr.app.Colors["group_text"],
		)
		index++

		description, nb_line := dr.WrapText(quest.Description, 24)

		// Draw Quest Description
		rl.DrawText(
			description,
			indX+int32(0.66*tile),
			start_y+int32(index)*int32(font),
			int32(font), dr.app.Colors["pseudo_text"],
		)

		index += nb_line
	}
}

func get_active_quests(all_quests []pr.TrackedQuestData) []pr.TrackedQuestData {
	active_quests := make([]pr.TrackedQuestData, 0)

	for _, quest := range all_quests {
		if quest.Status == "active" {
			active_quests = append(active_quests, quest)
		}
	}

	return active_quests
}

func get_finished_quests(all_quests []pr.TrackedQuestData) []pr.TrackedQuestData {
	finished_quests := make([]pr.TrackedQuestData, 0)

	for _, quest := range all_quests {
		if quest.Status != "active" {
			finished_quests = append(finished_quests, quest)
		}
	}

	return finished_quests
}
