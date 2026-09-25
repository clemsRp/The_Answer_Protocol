package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawCombatChatZone() {
	chat_width := float32(vars.COMBAT_CHAT_END_X - vars.COMBAT_CHAT_START_X + 1)
	chat_height := float32(vars.COMBAT_END_Y - vars.COMBAT_START_Y)
	chat_start_x := float32(vars.COMBAT_CHAT_START_X)
	chat_start_y := float32(vars.COMBAT_START_Y)

	dr.DrawCombatFrame(chat_start_x, chat_start_y, chat_start_x+chat_width, chat_start_y+chat_height-2)
	dr.DrawCombatChatInput(chat_start_x, chat_start_y, chat_width, chat_height-2)
	dr.DrawCombatChats(chat_start_x, chat_start_y, chat_height-2)

	for _, b := range dr.app.Manager.Buttons("CombatChat") {
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

func (dr *Drawer) DrawCombatChatInput(chat_start_x, chat_start_y, chat_width, chat_height float32) {
	dr.DrawWoodFrameAt(
		vars.Position{X: chat_start_x, Y: chat_start_y + chat_height},
		vars.Position{X: chat_start_x + chat_width, Y: chat_start_y + chat_height + 2},
		1, 0, false,
	)

	baseX := int32((chat_start_x + 0.63) * dr.app.Variables.Tileset_size)
	baseY := int32((chat_start_y + chat_height + 0.63) * dr.app.Variables.Tileset_size)

	input_width := (chat_width - 2) * dr.app.Variables.Tileset_size
	input_height := dr.app.Variables.Tileset_size * 0.75
	border := input_height * 0.1

	dr.DrawBorderedInput(
		baseX, baseY,
		int32(input_width), int32(input_height), int32(border),
		rl.NewColor(196, 154, 108, 255),
		rl.NewColor(232, 207, 166, 255),
	)

	rl.DrawText(
		dr.app.Variables.PanelsVariables.Chat.Msg,
		baseX+int32(0.4*dr.app.Variables.FontSize),
		baseY+int32(0.4*dr.app.Variables.FontSize),
		int32(dr.app.Variables.FontSize), dr.app.Colors["panel_text"],
	)

	dr.DrawCursor(
		int32(baseX+int32(0.4*dr.app.Variables.FontSize))+int32(border),
		int32(baseY+int32(0.4*dr.app.Variables.FontSize)),
		int32(border), int32(dr.app.Variables.FontSize),
		dr.app.Variables.PanelsVariables.Chat.Msg, dr.app.Colors["panel_text"],
	)
}

func (dr *Drawer) DrawCombatChats(chat_start_x, chat_start_y, chat_height float32) {
	tex := dr.chatTexture.Texture
	viewHeight := (chat_height - 1) * dr.app.Variables.Tileset_size

	displayHeight := viewHeight
	maxY := float32(tex.Height) - displayHeight
	minY := float32(0)

	if float32(tex.Height) < viewHeight {
		displayHeight = float32(tex.Height)
		maxY = 0
		dr.app.Variables.PanelsVariables.Chat.ScrollActive = false
	} else {
		dr.DrawCombatChatScrollBar(maxY, chat_start_x, chat_start_y, chat_height)
		dr.app.Variables.PanelsVariables.Chat.ScrollActive = true
	}

	sourceY := maxY
	scroll := dr.app.Variables.PanelsVariables.Chat.Scroll
	finalY := sourceY + scroll

	if finalY > maxY {
		finalY = maxY
		dr.app.Variables.PanelsVariables.Chat.Scroll = finalY - sourceY
	} else if finalY < minY {
		finalY = minY
		dr.app.Variables.PanelsVariables.Chat.Scroll = finalY - sourceY
	}

	sourceRec := rl.NewRectangle(0, finalY, float32(tex.Width), -displayHeight)
	position := rl.NewVector2(
		chat_start_x*dr.app.Variables.Tileset_size,
		(chat_start_y+0.5)*dr.app.Variables.Tileset_size,
	)

	rl.DrawTextureRec(tex, sourceRec, position, rl.White)
}

func (dr *Drawer) DrawCombatChatScrollBar(maxY, chat_start_x, chat_start_y, chat_height float32) {
	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32((chat_start_x+vars.CHAT_WIDTH-2)*dr.app.Variables.Tileset_size),
		float32((chat_start_y)*dr.app.Variables.Tileset_size),
		20, 8, 1, 1, dr.app.Variables.Zoom, 0,
	)

	limit := int(chat_height - 2)
	for mid := 0; mid < limit; mid++ {
		dr.DrawImage(
			vars.UI_SPRITE_TEXTURE,
			float32((chat_start_x+vars.CHAT_WIDTH-2)*dr.app.Variables.Tileset_size),
			float32((chat_start_y+float32(mid)+1)*dr.app.Variables.Tileset_size),
			20, 9, 1, 1, dr.app.Variables.Zoom, 0,
		)
	}

	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32((chat_start_x+vars.CHAT_WIDTH-2)*dr.app.Variables.Tileset_size),
		float32((chat_start_y+chat_height-1)*dr.app.Variables.Tileset_size),
		20, 10, 1, 1, dr.app.Variables.Zoom, 0,
	)

	scroll_bar_start_y := (chat_start_y + 0.5) * dr.app.Variables.Tileset_size
	scroll_bar_end_y := dr.app.Variables.Tileset_size * (chat_start_y + chat_height - 0.25)

	cursor_y := dr.app.Variables.PanelsVariables.Chat.ScrollBarY
	current_scroll := dr.app.Variables.PanelsVariables.Chat.Scroll

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		cursor_y = max(scroll_bar_start_y, cursor_y)
		cursor_y = min(cursor_y, scroll_bar_end_y)
		dr.app.Variables.PanelsVariables.Chat.ScrollBarY = cursor_y

		if scroll_bar_end_y > scroll_bar_start_y {
			percent := (cursor_y - scroll_bar_start_y) / (scroll_bar_end_y - scroll_bar_start_y)
			dr.app.Variables.PanelsVariables.Chat.Scroll = -maxY * percent
		}
	} else {
		if maxY > 0 {
			percent := current_scroll / -maxY
			cursor_y = scroll_bar_start_y + percent*(scroll_bar_end_y-scroll_bar_start_y)
			cursor_y = max(scroll_bar_start_y, cursor_y)
			cursor_y = min(cursor_y, scroll_bar_end_y)
			dr.app.Variables.PanelsVariables.Chat.ScrollBarY = cursor_y
		}
	}

	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32((chat_start_x+vars.CHAT_WIDTH-2)*dr.app.Variables.Tileset_size),
		cursor_y-dr.app.Variables.Tileset_size,
		19, 8, 1, 2, dr.app.Variables.Zoom, 0,
	)
}
