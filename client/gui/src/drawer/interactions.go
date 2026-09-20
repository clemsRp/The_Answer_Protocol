package drawer

func (dr *Drawer) DrawGameInteractions() {
	// Draw items
	for _, item := range dr.app.Manager.Interactions("Items") {

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

	// Draw npcs
	for _, npc := range dr.app.Manager.Interactions("Npcs") {

		// Draw main npc
		if npc.Emote != nil {
			frame := npc.Emote.CurrentFrame(dr.app.Variables.StartTime)
			dr.DrawImage(
				npc.Emote.Texture,
				npc.Emote.X, npc.Emote.Y,
				frame.IndX, frame.IndY,
				frame.RatioX, frame.RatioY,
				npc.Emote.Zoom,
				npc.Emote.Rotation,
			)
		}

		// Draw buttons only if hovered
		if npc.IsHovered() {
			for _, btn := range npc.Buttons {
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
