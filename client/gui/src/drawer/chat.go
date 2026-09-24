package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawChat() {
	chat_width := float32(vars.CHAT_WIDTH)
	chat_height := float32(vars.CHAT_HEIGHT)
	chat_start_x := float32(vars.CHAT_START_X)
	chat_start_y := float32(vars.CHAT_START_Y)

	dr.drawChatPanel(chat_start_x, chat_start_y, chat_width, chat_height)
	dr.drawChatInput(chat_start_x, chat_start_y, chat_width, chat_height)

	// Draw chats
	dr.DrawChats(chat_start_x, chat_start_y, chat_height)

	dr.DrawChatEmotes()
	dr.DrawChatButtons()

	dr.DrawScope()
}

func (dr *Drawer) RenderChatTexture() {
	chat := dr.app.Variables.PanelsVariables.Chat
	if !chat.Open {
		return
	}
	dr.initChatTexture(float32(vars.CHAT_WIDTH))

	chats := chat.ScopeChats[chat.CurrentScope]
	rl.BeginTextureMode(dr.chatTexture)
	rl.ClearBackground(rl.Blank)
	for ind, c := range chats {
		dr.DrawChatMsg(c, 0, 0, float32(ind), false)
	}
	rl.EndTextureMode()
}

func (dr *Drawer) initChatTexture(chat_width float32) {
	chat := dr.app.Variables.PanelsVariables.Chat
	chats := chat.ScopeChats[chat.CurrentScope]

	if chat.LastNbChats != len(chats) || chat.LastMsgScope != chat.CurrentScope || dr.chatTexture.ID == 0 {
		if dr.chatTexture.ID != 0 {
			rl.UnloadRenderTexture(dr.chatTexture)
		}
		dr.chatTexture = rl.LoadRenderTexture(
			int32(chat_width*dr.app.Variables.Tileset_size),
			int32((2*float32(len(chats))+4.5)*dr.app.Variables.Tileset_size),
		)
		chat.LastNbChats = len(chats)
		chat.LastMsgScope = chat.CurrentScope
	}
}

func (dr *Drawer) drawChatPanel(chat_start_x, chat_start_y, chat_width, chat_height float32) {
	// Draw panel
	dr.DrawWoodFrameAt(
		vars.Position{X: chat_start_x, Y: chat_start_y},
		vars.Position{X: chat_start_x + chat_width, Y: chat_start_y + chat_height},
		1, 0, false,
	)
}

func (dr *Drawer) drawChatInput(chat_start_x, chat_start_y, chat_width, chat_height float32) {
	// Draw input
	dr.DrawWoodFrameAt(
		vars.Position{X: chat_start_x, Y: chat_start_y + chat_height},
		vars.Position{X: chat_start_x + chat_width, Y: 17},
		1, 0, false,
	)

	baseX := int32((chat_start_x + 0.63) * dr.app.Variables.Tileset_size)
	baseY := int32((chat_start_y + chat_height + 0.63) * dr.app.Variables.Tileset_size)

	input_width := (chat_width - 2) * dr.app.Variables.Tileset_size
	input_height := dr.app.Variables.Tileset_size * 0.75
	border := input_height * 0.1

	dr.DrawBorderedInput(
		baseX,
		baseY,
		int32(input_width),
		int32(input_height),
		int32(border),
		rl.NewColor(196, 154, 108, 255),
		rl.NewColor(232, 207, 166, 255),
	)

	// Draw message
	rl.DrawText(
		dr.app.Variables.PanelsVariables.Chat.Msg,
		baseX+int32(0.4*dr.app.Variables.FontSize),
		baseY+int32(0.4*dr.app.Variables.FontSize),
		int32(dr.app.Variables.FontSize), dr.app.Colors["panel_text"],
	)

	// Draw cursor
	dr.DrawCursor(
		int32(baseX+int32(0.4*dr.app.Variables.FontSize))+int32(border),
		int32(baseY+int32(0.4*dr.app.Variables.FontSize)),
		int32(border), int32(dr.app.Variables.FontSize),
		dr.app.Variables.PanelsVariables.Chat.Msg, dr.app.Colors["panel_text"],
	)
}

func (dr *Drawer) DrawChats(chat_start_x, chat_start_y, chat_height float32) {
	// Display needed part of the chat texture for scroll system
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
		dr.DrawChatScrollBar(maxY)
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

func (dr *Drawer) DrawChatMsg(chat vars.Chat, start_x, start_y, index float32, border_bottom bool) {
	posX := float32(start_x * dr.app.Variables.Tileset_size)
	posY := float32((start_y + 2*index) * dr.app.Variables.Tileset_size)

	// Draw Woodframe
	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		posX, posY,
		12, 6, 3, 3,
		dr.app.Variables.Zoom, 0,
	)

	// Draw dialogue bubble
	box_type := "big"
	var box_width float32
	box_width = 11

	// Draw box/frame
	if len(chat.Msg) <= vars.MSG_MAX_CHAR/3 {
		box_type = "small"
		box_width = 7

	} else if len(chat.Msg) <= vars.MSG_MAX_CHAR*2/3 {
		box_type = "medium"
		box_width = 8
	}

	dr.DrawImage(
		"dialog box "+box_type,
		posX+2.4*dr.app.Variables.Tileset_size,
		posY+dr.app.Variables.Tileset_size,
		0, 0, box_width, 3,
		dr.app.Variables.Zoom/2, 0,
	)

	// Draw emote
	e := dr.app.Manager.Emotes("Connect")[chat.EmoteIndex]
	frame := e.CurrentFrame(dr.app.Variables.StartTime)
	dr.DrawImage(
		e.Texture,
		posX+dr.app.Variables.Tileset_size,
		posY+dr.app.Variables.Tileset_size,
		frame.IndX, frame.IndY,
		frame.RatioX, frame.RatioY,
		e.Zoom/2,
		e.Rotation,
	)

	// Draw Pseudo
	rl.DrawText(
		chat.Pseudo,
		int32(posX+2.7*dr.app.Variables.Tileset_size),
		int32(posY+0.7*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["pseudo_text"],
	)

	// Draw message
	rl.DrawText(
		chat.Msg,
		int32(posX+3*dr.app.Variables.Tileset_size),
		int32(posY+1.5*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)

	// Draw border bottom
	if border_bottom {
		rl.DrawRectangle(
			int32(posX+dr.app.Variables.Tileset_size),
			int32(posY+2*dr.app.Variables.Tileset_size),
			int32(8*dr.app.Variables.Tileset_size), 3,
			rl.Black,
		)
	}
}

func (dr *Drawer) DrawChatScrollBar(maxY float32) {
	// Draw Top part
	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32((vars.CHAT_START_X+vars.CHAT_WIDTH-1)*dr.app.Variables.Tileset_size),
		float32((vars.CHAT_START_Y)*dr.app.Variables.Tileset_size),
		20, 8, 1, 1, dr.app.Variables.Zoom, 0,
	)

	// Draw Middle part
	limit := int(vars.CHAT_HEIGHT - 2)
	for mid := 0; mid < limit; mid++ {
		dr.DrawImage(
			vars.UI_SPRITE_TEXTURE,
			float32((vars.CHAT_START_X+vars.CHAT_WIDTH-1)*dr.app.Variables.Tileset_size),
			float32((vars.CHAT_START_Y+float32(mid)+1)*dr.app.Variables.Tileset_size),
			20, 9, 1, 1, dr.app.Variables.Zoom, 0,
		)
	}

	// Draw Bottom part
	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32((vars.CHAT_START_X+vars.CHAT_WIDTH-1)*dr.app.Variables.Tileset_size),
		float32((vars.CHAT_START_Y+vars.CHAT_HEIGHT-1)*dr.app.Variables.Tileset_size),
		20, 10, 1, 1, dr.app.Variables.Zoom, 0,
	)

	scroll_bar_start_y := float32(1.5) * dr.app.Variables.Tileset_size
	scroll_bar_end_y := dr.app.Variables.Tileset_size * (float32(vars.CHAT_START_Y) + float32(vars.CHAT_HEIGHT) - 0.25)

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

	// Draw cursor
	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32((vars.CHAT_START_X+vars.CHAT_WIDTH-1)*dr.app.Variables.Tileset_size),
		cursor_y-dr.app.Variables.Tileset_size,
		19, 8, 1, 2, dr.app.Variables.Zoom, 0,
	)
}

func (dr *Drawer) DrawScope() {
	// Draw Frame
	dr.DrawWoodFrameAt(
		vars.Position{X: vars.CHAT_START_X + 2, Y: 0.5},
		vars.Position{X: vars.CHAT_START_X + vars.CHAT_WIDTH - 2, Y: 2},
		0.5, 2, false,
	)

	// Draw scope
	scope := dr.app.Variables.PanelsVariables.Chat.CurrentScope
	font_size := 3 * int32(dr.app.Variables.FontSize)
	center_text := float32(rl.MeasureText(scope, font_size) / 2)

	rl.DrawText(
		scope,
		int32((vars.CHAT_START_X+vars.CHAT_WIDTH/2)*dr.app.Variables.Tileset_size-center_text),
		int32(vars.CHAT_START_Y+font_size*13/22),
		font_size, rl.NewColor(232, 207, 166, 255),
	)
}
