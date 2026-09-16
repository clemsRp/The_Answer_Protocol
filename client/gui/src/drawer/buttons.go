package drawer

func (dr *Drawer) DrawConnectButtons() {
	for _, b := range dr.app.Manager.Buttons("Connect") {
		frame := b.CurrentFrame()
		dr.DrawImage(
			b.Texture,
			b.X, b.Y,
			frame.IndX, frame.IndY,
			frame.RatioX, frame.RatioY,
			b.Zoom,
			b.Rotation,
		)
	}
}

func (dr *Drawer) DrawGameButtons() {
	for _, b := range dr.app.Manager.Buttons("Game") {
		frame := b.CurrentFrame()
		dr.DrawImage(
			b.Texture,
			b.X, b.Y,
			frame.IndX, frame.IndY,
			frame.RatioX, frame.RatioY,
			b.Zoom,
			b.Rotation,
		)
	}
}

func (dr *Drawer) DrawChatButtons() {
	for _, b := range dr.app.Manager.Buttons("Chat") {
		frame := b.CurrentFrame()
		dr.DrawImage(
			b.Texture,
			b.X, b.Y,
			frame.IndX, frame.IndY,
			frame.RatioX, frame.RatioY,
			b.Zoom,
			b.Rotation,
		)
	}
}
