package drawer

func (dr *Drawer) DrawConnectEmotes() {
	e := dr.app.Manager.Emotes("Connect")[dr.app.Variables.Player.EmoteIndex]
	frame := e.CurrentFrame(dr.app.Variables.StartTime)
	dr.DrawImage(
		e.Texture,
		e.X, e.Y,
		frame.IndX, frame.IndY,
		frame.RatioX, frame.RatioY,
		e.Zoom,
		e.Rotation,
	)
}

func (dr *Drawer) DrawGameEmotes() {
	e := dr.app.Manager.Emotes("Game")[dr.app.Variables.Player.EmoteIndex]
	frame := e.CurrentFrame(dr.app.Variables.StartTime)
	dr.DrawImage(
		e.Texture,
		e.X, e.Y,
		frame.IndX, frame.IndY,
		frame.RatioX, frame.RatioY,
		e.Zoom,
		e.Rotation,
	)
}
