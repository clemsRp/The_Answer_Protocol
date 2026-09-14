package drawer

func (dr *Drawer) DrawGameView() {
	dr.DrawMap()
	dr.DrawPlayers()

	for _, layer := range dr.app.Rooms[dr.app.Variables.Current_room].Layers {
		if layer.Name == "InFrontOfPlayer" {
			dr.DrawLayer(layer, dr.app.Rooms[dr.app.Variables.Current_room].Tilesets)
		}
	}

	dr.DrawPseudos()
}
