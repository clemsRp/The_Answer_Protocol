package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawChat() {
	// Draw panel
	chat_width := float32(10)
	chat_height := float32(14)
	chat_start_x := float32(32 - chat_width - 1)
	chat_start_y := float32(1)

	dr.DrawWoodFrameAt(
		vars.Position{X: chat_start_x, Y: chat_start_y},
		vars.Position{X: chat_start_x + chat_width, Y: chat_start_y + chat_height},
		1, 0, false,
	)

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

	// Draw chats
	scope := dr.app.Variables.PanelsVariables.Chat.CurrentScope
	chats := dr.app.Variables.PanelsVariables.Chat.ScopeChats[scope]

	for ind, chat := range chats {
		dr.DrawChatMsg(
			chat, chat_start_x, chat_start_y,
			float32(ind), false,
		)
	}

	dr.DrawChatEmotes()
	dr.DrawChatButtons()
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
		posX+2.5*dr.app.Variables.Tileset_size,
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
		dr.app.Colors["chat_pseudo_text"],
	)

	// Draw message
	rl.DrawText(
		chat.Msg,
		int32(posX+3*dr.app.Variables.Tileset_size),
		int32(posY+1.5*dr.app.Variables.Tileset_size),
		int32(dr.app.Variables.FontSize),
		dr.app.Colors["panel_text"],
	)

	// Draw border top
	if border_bottom {
		rl.DrawRectangle(
			int32(posX+dr.app.Variables.Tileset_size),
			int32(posY+2.5*dr.app.Variables.Tileset_size),
			int32(8*dr.app.Variables.Tileset_size), 3,
			rl.Black,
		)
	}
}
