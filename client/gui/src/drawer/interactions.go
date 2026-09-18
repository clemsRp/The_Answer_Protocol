package drawer

func (dr *Drawer) DrawGameInteractions() {
	for _, item := range dr.app.Manager.Interactions("Game") {

		// Draw main item
		if item.Emote != nil {
			frame := item.Emote.CurrentFrame(dr.app.Variables.StartTime)
			dr.DrawImage(
				item.Emote.Texture,
				item.Emote.X, item.Emote.Y,
				frame.IndX, frame.IndY,
				frame.RatioX, frame.RatioY,
				item.Emote.Zoom,
				item.Emote.Rotation,
			)
		}

		// Draw buttons only if hovered
		if item.IsHovered() {
			for _, btn := range item.Buttons {
				btn_frame := btn.CurrentFrame()
				dr.DrawImage(
					btn.Texture,
					btn.X, btn.Y,
					btn_frame.IndX, btn_frame.IndY,
					btn_frame.RatioX, btn_frame.RatioY,
					btn.Zoom,
					btn.Rotation,
				)
			}
		}
	}
}
