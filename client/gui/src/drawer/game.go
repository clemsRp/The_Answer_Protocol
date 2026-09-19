package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGameView() {
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
		if layer.Name == "InFrontOfPlayer" {
			dr.DrawLayer(layer, dr.app.Rooms[dr.app.Variables.Current_room].Tilesets)
		}
	}

	dr.DrawPlayerPanel()

	dr.DrawGameEmotes()
	dr.DrawGameButtons()
	dr.DrawInventoryEmotes()
	dr.DrawInventoryButtons()
}

func (dr *Drawer) DrawPlayerPanel() {
	dr.drawPlayerInfoBox()
	dr.drawPlayerLives()

	// Draw chat
	if dr.app.Variables.PanelsVariables.Chat.Open {
		dr.DrawChat()
	}

	// Draw left panel
	if dr.app.Variables.PanelsVariables.LeftPanel.Open {
		dr.DrawLeftPanel()
	}

	dr.DrawInventory()
}

func (dr *Drawer) drawPlayerInfoBox() {
	// Draw Player datas up/left corner
	pseudo := dr.app.Variables.Player.Pseudo
	box_type := " big"
	box_width := float32(19)

	// Draw box/frame
	if len(pseudo) <= vars.PSEUDO_MAX_CHAR/3 {
		box_type = "small"
		box_width = 11
	} else if len(pseudo) <= vars.PSEUDO_MAX_CHAR*2/3 {
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
	live := float64(dr.app.Variables.Player.Hp) / float64(dr.app.Variables.Player.MaxHp)

	for nb_live := 0; nb_live < 5; nb_live++ {
		indX := 0
		indY := 3
		ratio := 2

		heart_start := 0.2 * float64(nb_live)
		heart_half := heart_start + 0.1
		heart_full := heart_start + 0.2

		if live < heart_full {
			indX += 2
		}
		if live < heart_half {
			indX += 2
		}

		if live < 0.1 && nb_live == 0 {
			indX = 2
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

	if len(pseudo) <= vars.PSEUDO_MAX_CHAR/3 {
		start_x = 11
	} else if len(pseudo) <= vars.PSEUDO_MAX_CHAR*2/3 {
		start_x = 15
	}
	start_x += 2
	start_x *= dr.app.Variables.Zoom / 2 * vars.FRAME_WIDTH
	start_x -= 0.1 * dr.app.Variables.Tileset_size

	if dr.app.Variables.PanelsVariables.Inventory.Open {
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
}
