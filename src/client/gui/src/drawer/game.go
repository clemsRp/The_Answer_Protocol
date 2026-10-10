package drawer

import (
	vars "tap/src/client/gui/src/variables"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGameView() {
	rl.BeginTextureMode(dr.gameTexture)
	dr.DrawGame()
	rl.EndTextureMode()

	talking := dr.app.Variables.PanelsVariables.Talk.Talking != nil
	if talking {
		rl.BeginShaderMode(dr.darkenShader)
	}

	rl.DrawTextureRec(
		dr.gameTexture.Texture,
		rl.NewRectangle(
			0, 0,
			float32(dr.gameTexture.Texture.Width),
			-float32(dr.gameTexture.Texture.Height),
		),
		rl.NewVector2(0, 0),
		rl.White,
	)

	if talking {
		rl.EndShaderMode()
		dr.DrawTalk()
	}
}

func (dr *Drawer) DrawGame() {
	dr.DrawMap()
	dr.DrawGameInteractions()

	dr.DrawPlayers()
	dr.DrawPseudos()

	dr.DrawWoodFrameAt(
		vars.Position{X: 0, Y: 0},
		vars.Position{X: 32, Y: 18},
		2, 0, true,
	)

	for _, layer := range dr.app.Rooms[dr.app.Variables.Current_room].Layers {
		if layer.Name == "InFrontOfPlayer" && layer.Visible {
			dr.DrawLayer(layer, dr.app.Rooms[dr.app.Variables.Current_room].Tilesets)
		}
	}

	dr.DrawPlayerPanel()

	dr.DrawGameEmotes()
	dr.DrawGameButtons()
	dr.DrawInventoryEmotes()
	dr.DrawInventoryButtons()
	dr.DrawGameSelects()

	dr.HideBorder()
}

func (dr *Drawer) DrawPlayerPanel() {
	dr.drawPlayerInfoBox()
	dr.drawPlayerLives()

	// Draw panels
	if dr.app.Variables.Current_view == "Game" {
		// Draw chat
		if dr.app.Variables.PanelsVariables.Chat.Open {
			dr.DrawChat()
		}

		// Draw group
		if dr.app.Variables.PanelsVariables.Group.Open {
			dr.DrawGroupPanel()
		}

		// Draw inspect
		if dr.app.Variables.PanelsVariables.Inspect.Open {
			dr.DrawInspectPanel()
		}

		// Draw datas
		if dr.app.Variables.PanelsVariables.Datas.Open {
			dr.DrawDatasPanel()
		}
	}

	dr.DrawInventory()
}

func (dr *Drawer) drawPlayerInfoBox() {
	// Draw Player datas up/left corner
	pseudo := dr.app.Variables.Player.Pseudo
	box_type := " big"
	box_width := float32(19)

	tier := dr.app.GetPseudoBoxTier(pseudo)
	if tier == 1 {
		box_type = "small"
		box_width = 11
	} else if tier == 2 {
		box_type = "medium"
		box_width = 15
	}

	dr.DrawImage(
		"Premade dialog box "+box_type,
		dr.app.Variables.Tileset_size,
		dr.app.Variables.Tileset_size,
		0, 0, box_width, 4,
		dr.app.Variables.Zoom/2, 0,
	)

	// Draw pseudo
	rl.DrawText(
		pseudo,
		int32(3.25*dr.app.Variables.Tileset_size),
		int32(1.83*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)
}

func (dr *Drawer) drawPlayerLives() {
	// Draw lives
	hp := dr.app.Variables.Player.Hp
	maxHp := dr.app.Variables.Player.MaxHp

	half_hearts := 0
	if maxHp > 0 {
		half_hearts = (hp*10 + maxHp - 1) / maxHp
	}
	if half_hearts > 10 {
		half_hearts = 10
	}
	if half_hearts < 0 {
		half_hearts = 0
	}

	for nb_live := 0; nb_live < 5; nb_live++ {
		indX := 0
		indY := 3
		ratio := 2

		heart_threshold := 2 * nb_live
		if half_hearts <= heart_threshold {
			indX = 4
		} else if half_hearts == heart_threshold+1 {
			indX = 2
		} else {
			indX = 0
		}

		// Draw heart
		dr.DrawImage(
			vars.HEART_TEXTURE,
			(float32(nb_live)*0.6+2.9)*dr.app.Variables.Tileset_size,
			1.03*dr.app.Variables.Tileset_size,
			float32(indX), float32(indY), float32(ratio), float32(ratio),
			dr.app.Variables.Zoom/3, 0,
		)
	}
}

func (dr *Drawer) DrawInventory() {
	pseudo := dr.app.Variables.Player.Pseudo
	start_x := float32(19)

	tier := dr.app.GetPseudoBoxTier(pseudo)
	if tier == 1 {
		start_x = 11
	} else if tier == 2 {
		start_x = 15
	}
	start_x += 2
	start_x *= dr.app.Variables.Zoom / 2 * vars.FRAME_WIDTH
	start_x -= 0.1 * dr.app.Variables.Tileset_size

	for index := range *dr.app.Variables.PanelsVariables.InventoryItems {
		indX := 12.4
		if index == len(*dr.app.Variables.PanelsVariables.InventoryItems)-1 {
			indX += 3
		}

		// Draw frame
		dr.DrawImage(
			vars.INVENTORY_TEXTURE,
			start_x+float32(index)*dr.app.Variables.Tileset_size,
			1.1*dr.app.Variables.Tileset_size,
			float32(indX), 1, 3.3, 5,
			dr.app.Variables.Zoom/3, 0,
		)
	}
}

func (dr *Drawer) DrawRoomName() {
	cur_room := dr.app.Variables.PanelsVariables.Room.Room.Name

	center := rl.MeasureText(
		cur_room,
		int32(dr.app.Variables.FontSize),
	) / 2

	align := float32(center)/dr.app.Variables.Tileset_size + 1

	// Draw frame
	dr.DrawWoodFrameAt(
		vars.Position{X: 15.5 - align, Y: 0},
		vars.Position{X: 15.5 + align, Y: 1.5},
		0.5, 1, false,
	)

	// Draw room name
	rl.DrawText(
		cur_room,
		int32(15.5*dr.app.Variables.Tileset_size)-int32(center),
		int32(0.75*dr.app.Variables.Tileset_size)-int32(dr.app.Variables.FontSize/2),
		int32(dr.app.Variables.FontSize), dr.app.Colors["pseudo_text"],
	)

	cur_description := dr.app.Variables.PanelsVariables.Room.Room.Description

	center = rl.MeasureText(
		cur_description,
		int32(dr.app.Variables.FontSize),
	) / 2

	align = float32(center)/dr.app.Variables.Tileset_size + 1

	// Draw frame
	dr.DrawWoodFrameAt(
		vars.Position{X: 15.5 - align, Y: 16},
		vars.Position{X: 15.5 + align, Y: 17.5},
		0.5, 1, false,
	)

	// Draw room name
	rl.DrawText(
		cur_description,
		int32(15.5*dr.app.Variables.Tileset_size)-int32(center),
		int32(16.75*dr.app.Variables.Tileset_size)-int32(dr.app.Variables.FontSize/2),
		int32(dr.app.Variables.FontSize), dr.app.Colors["pseudo_text"],
	)
}

func (dr *Drawer) HideBorder() {
	panels := dr.app.Variables.PanelsVariables
	datas := panels.Datas.Open
	inspect := panels.Inspect.Open
	group := panels.Group.Open
	chat := panels.Chat.Open
	nothing_open := !datas && !inspect && !group && !chat

	in_time := time.Since(dr.app.Variables.LastRoomChange) <= 5*time.Second

	rl.DrawRectangle(
		int32(6*dr.app.Variables.Tileset_size),
		int32(0),
		int32(18*dr.app.Variables.Tileset_size),
		int32(0.39*dr.app.Variables.Tileset_size),
		rl.Black,
	)

	if dr.app.Variables.Current_view == "Game" && nothing_open && in_time {
		dr.DrawRoomName()
	}
}
