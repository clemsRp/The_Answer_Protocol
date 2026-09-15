package drawer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	vars "tap/client/gui/src/variables"
)

func (dr *Drawer) DrawGameView() {
	dr.DrawMap()
	dr.DrawPlayers()

	for _, layer := range dr.app.Rooms[dr.app.Variables.Current_room].Layers {
		if layer.Name == "InFrontOfPlayer" {
			dr.DrawLayer(layer, dr.app.Rooms[dr.app.Variables.Current_room].Tilesets)
		}
	}

	dr.DrawPseudos()

	dr.DrawPlayerPanel()

	dr.DrawGameEmotes()
	dr.DrawGameButtons()
}

func (dr *Drawer) DrawPlayerPanel() {
	// Draw Player datas up/right corner
	pseudo := dr.app.Variables.Player.Pseudo
	box_type := " big"
	var box_width float32
	box_width = 19

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
