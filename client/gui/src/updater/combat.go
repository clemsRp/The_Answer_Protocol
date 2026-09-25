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

	if len(up.app.Manager.Buttons("CombatActions")) == 0 {
		up.buildCombatActionsButtons()
	}

	if up.app.Variables.PanelsVariables.InventoryItems != nil {
		combat_use_buttons := up.app.GetNewCombatInventory(*up.app.Variables.PanelsVariables.InventoryItems)
		up.app.Manager.SetViewButtons("CombatInventory", combat_use_buttons)
	}

	up.app.Manager.Update("Combat")

	if up.app.Variables.PanelsVariables.CombatState.CurrentTurn == up.app.GetPseudo() {
		up.app.Manager.Update("CombatActions")
	}

	up.app.Manager.Update("CombatChat")
	up.app.Manager.Update("CombatInventory")

	up.UpdateCombatChat()
}

func (up *Updater) UpdateCombatChat() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		up.SendCombatChat()
	}

	up.UpdateChatMsg()
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

func (up *Updater) buildCombatActionsButtons() {
	splitX := float32(vars.COMBAT_START_X) + float32(vars.COMBAT_LEFT_END_X-vars.COMBAT_START_X)/2.0
	tile := up.app.Variables.Tileset_size
	zoom := up.app.Variables.Zoom * 11 / 20

	leftCenterX := (float32(vars.COMBAT_START_X) + splitX) / 2 * tile
	btnWidth := 6 * float32(vars.FRAME_WIDTH) * zoom
	btnX := leftCenterX - btnWidth/2

	btnAttackY := (float32(vars.COMBAT_ACTIONS_START_Y) + 0.8) * tile
	btnFleeY := (float32(vars.COMBAT_ACTIONS_START_Y) + 3.2) * tile

	attackBtn := &ui.Button{
		ID:      "combat_attack",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       btnX,
		Y:       btnAttackY,
		Zoom:    zoom,
		Normal:  ui.Frame{IndX: 10, IndY: 11, RatioX: 6, RatioY: 2},
		Pressed: ui.Frame{IndX: 16, IndY: 11, RatioX: 6, RatioY: 2},
		OnClick: func() {
			cs := up.app.Variables.PanelsVariables.CombatState
			target := ""
			if cs != nil {
				target = cs.SelectedPerson
				if target == "" && len(cs.Opponents) > 0 {
					for k := range cs.Opponents {
						target = k
						break
					}
				}
			}
			payload := pr.CmdAttack
			if target != "" {
				payload = fmt.Sprintf("%s %s", pr.CmdAttack, target)
			}
			up.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: payload,
			}
		},
	}

	fleeBtn := &ui.Button{
		ID:      "combat_flee",
		Texture: vars.UI_SPRITE_TEXTURE,
		X:       btnX,
		Y:       btnFleeY,
		Zoom:    zoom,
		Normal:  ui.Frame{IndX: 10, IndY: 11, RatioX: 6, RatioY: 2},
		Pressed: ui.Frame{IndX: 16, IndY: 11, RatioX: 6, RatioY: 2},
		OnClick: func() {
			up.actionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdFlee,
			}
		},
	}

	up.app.Manager.SetViewButtons("CombatActions", []*ui.Button{attackBtn, fleeBtn})
}
