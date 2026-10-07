package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

		npc_name := dr.app.Variables.NpcConvertor[npc.ID]

		center := rl.MeasureText(
			npc_name, int32(dr.app.Variables.FontSize),
		) / 2

		ratioX := float32(1)
		if npc.Emote != nil && len(npc.Emote.Frames) > 0 {
			ratioX = npc.Emote.Frames[0].Frame.RatioX
		}

		spriteWidth := float32(vars.FRAME_WIDTH) * npc.Emote.Zoom * ratioX
		x := npc.Emote.X + (spriteWidth / 2) - float32(center)
		y := npc.Emote.Y - float32(dr.app.Variables.FontSize)

		// Draw npc name
		rl.DrawText(
			npc_name,
			int32(x),
			int32(y),
			int32(dr.app.Variables.FontSize),
			rl.Yellow,
		)

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
