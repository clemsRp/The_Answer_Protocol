package updater

import (
	"fmt"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdateCombatView() {
	up.UpdateGameView()

	up.app.Variables.PanelsVariables.Chat.Open = true
	if up.app.Variables.PanelsVariables.Chat.CurrentScope != "COMBAT" {
		up.app.Variables.PanelsVariables.Chat.CurrentScope = "COMBAT"
	}

	if len(up.app.Manager.Buttons("CombatChat")) == 0 {
		up.buildCombatChatButtons()
	}

	up.app.Manager.Update("Combat")
	up.app.Manager.Update("CombatChat")

	up.UpdateCombatChat()
}

func (up *Updater) UpdateCombatChat() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		up.SendCombatChat()
	}

	up.UpdateChatMsg() // Réutilisation de la logique de frappe standard
	up.UpdateCombatChatScroll()
}

func (up *Updater) SendCombatChat() {
	msg := up.app.Variables.PanelsVariables.Chat.Msg
	if msg == "" {
		return
	}

	up.app.Variables.PanelsVariables.Chat.Msg = ""
	up.app.Variables.PanelsVariables.Chat.LastMsgText = msg
	up.app.Variables.PanelsVariables.Chat.LastMsgScope = "COMBAT"

	up.actionsChan <- panel.Action{
		Type:    panel.ActionSendServer,
		Payload: fmt.Sprintf("%s %s", pr.CmdChatCombat, msg),
	}
}

func (up *Updater) UpdateCombatChatScroll() {
	mouse := rl.GetMousePosition()

	chat_start_x := float32(vars.COMBAT_CHAT_START_X)
	chat_width := float32(vars.COMBAT_CHAT_END_X - vars.COMBAT_CHAT_START_X + 1)
	chat_start_y := float32(vars.COMBAT_START_Y)
	chat_height := float32(vars.COMBAT_END_Y - vars.COMBAT_START_Y)

	rect := rl.NewRectangle(
		chat_start_x*up.app.Variables.Tileset_size,
		chat_start_y*up.app.Variables.Tileset_size,
		chat_width*up.app.Variables.Tileset_size,
		chat_height*up.app.Variables.Tileset_size,
	)

	scrollRect := rl.NewRectangle(
		(chat_start_x+chat_width-0.6)*up.app.Variables.Tileset_size,
		(chat_start_y+0.25)*up.app.Variables.Tileset_size,
		0.25*up.app.Variables.Tileset_size,
		(chat_height-0.5)*up.app.Variables.Tileset_size,
	)

	hover := rl.CheckCollisionPointRec(mouse, rect)
	scroll := rl.GetMouseWheelMove()
	scroll_speed := 35

	if hover {
		up.app.Variables.PanelsVariables.Chat.Scroll += scroll * float32(scroll_speed)
	}

	if up.app.Variables.PanelsVariables.Chat.ScrollActive {
		clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
		down := rl.IsMouseButtonDown(rl.MouseButtonLeft)
		scroll_hover := rl.CheckCollisionPointRec(mouse, scrollRect)
		last_frame_scroll := up.app.Variables.PanelsVariables.Chat.LastFrameScroll

		up.app.Variables.PanelsVariables.Chat.LastFrameScroll = false

		if (scroll_hover && (clicked || down)) || (last_frame_scroll && down) {
			up.app.Variables.PanelsVariables.Chat.LastFrameScroll = true
			up.app.Variables.PanelsVariables.Chat.ScrollBarY = mouse.Y
		}
	}
}

func (up *Updater) buildCombatChatButtons() {
	chat_start_x := float32(vars.COMBAT_CHAT_START_X)
	chat_width := float32(vars.COMBAT_CHAT_END_X - vars.COMBAT_CHAT_START_X + 1)
	chat_start_y := float32(vars.COMBAT_START_Y)
	chat_height := float32(vars.COMBAT_END_Y - vars.COMBAT_START_Y)

	sendchatBtn := &ui.Button{
		ID:      "send_combat_chat",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       (chat_start_x + chat_width - 1.4) * up.app.Variables.Tileset_size,
		Y:       (chat_start_y + chat_height - 1.4) * up.app.Variables.Tileset_size,
		Zoom:    up.app.Variables.Zoom * 0.5,
		Normal:  ui.Frame{IndX: 48, IndY: 0, RatioX: 2, RatioY: 2},
		Pressed: ui.Frame{IndX: 50, IndY: 0, RatioX: 2, RatioY: 2},
		OnClick: up.SendCombatChat,
	}

	up.app.Manager.SetViewButtons("CombatChat", []*ui.Button{sendchatBtn})
}
