package drawer

import rl "github.com/gen2brain/raylib-go/raylib"

func (dr *Drawer) DrawGameView() {
	// Init Buttons/Emotes
	// if len(dr.app.Manager.Buttons("Game")) == 0 {
	// 	dr.buildGameButtons()
	// }
	if len(dr.app.Manager.Emotes("Game")) == 0 {
		dr.buildGameEmotes()
	}

	dr.DrawMap()
	dr.DrawPlayers()

	for _, layer := range dr.app.Rooms[dr.app.Variables.Current_room].Layers {
		if layer.Name == "InFrontOfPlayer" {
			dr.DrawLayer(layer, dr.app.Rooms[dr.app.Variables.Current_room].Tilesets)
		}
	}

	dr.DrawPseudos()

	// Draw Player datas up/right corner
	dr.DrawImage(
		"Premade dialog box "box_type,
		dr.app.Variables.Tileset_size,
		dr.app.Variables.Tileset_size,
		0, 0, 11, 4,
		dr.app.Variables.Zoom/2, 0,
	)

	rl.DrawText(
		dr.app.Variables.Player.Pseudo,
		int32(3.25*dr.app.Variables.Tileset_size),
		int32(1.83*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)

	dr.DrawEmotes("Game")
	dr.DrawButtons("Game")
}
