package updater

import (
	"fmt"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	"time"

	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (up *Updater) UpdateChat() {
	if len(up.app.Manager.Buttons("Chat")) == 0 {
		up.buildChatButtons()
	}
	if len(up.app.Manager.Emotes("Chat")) == 0 {
		up.buildChatEmotes()
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		up.SendChat()
	}

	up.app.Manager.Update("Chat")

	up.UpdateChatMsg()
	up.UpdateChatScroll()
	up.UpdateScope()
}

func (up *Updater) UpdateChatMsg() {
	msg := up.app.Variables.PanelsVariables.Chat.Msg
	if rl.IsKeyPressed(rl.KeyBackspace) || rl.IsKeyPressedRepeat(rl.KeyBackspace) {
		up.app.Variables.Player.LastTimeTyped = time.Now()
		if len(msg) > 0 {
			runes := []rune(msg)
			msg = string(runes[:len(runes)-1])
		}
	}

	char := rl.GetCharPressed()

	for char > 0 {
		up.app.Variables.Player.LastTimeTyped = time.Now()
		if char >= 32 && len(msg) < vars.MSG_MAX_CHAR {
			msg += string(char)
		}

		char = rl.GetCharPressed()
	}

	up.app.Variables.PanelsVariables.Chat.Msg = msg
}

func (up *Updater) SendChat() {
	scope := up.app.Variables.PanelsVariables.Chat.CurrentScope
	msg := up.app.Variables.PanelsVariables.Chat.Msg

	up.app.Variables.PanelsVariables.Chat.Msg = ""
	up.app.Variables.PanelsVariables.Chat.LastMsgText = msg
	up.app.Variables.PanelsVariables.Chat.LastMsgScope = scope

	up.app.ActionsChan <- panel.Action{
		Type:    panel.ActionSendServer,
		Payload: fmt.Sprintf("%s %s %s", pr.CmdChat, scope, msg),
	}
}

func (up *Updater) UpdateChatScroll() {
	mouse := rl.GetMousePosition()
	hover := rl.CheckCollisionPointRec(mouse, up.app.Variables.PanelsVariables.Chat.Rect)

	scroll := rl.GetMouseWheelMove()
	scroll_speed := 35

	if hover {
		up.app.Variables.PanelsVariables.Chat.Scroll += scroll * float32(scroll_speed)
	}

	if up.app.Variables.PanelsVariables.Chat.ScrollActive {
		clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
		down := rl.IsMouseButtonDown(rl.MouseButtonLeft)
		scroll_hover := rl.CheckCollisionPointRec(mouse, up.app.Variables.PanelsVariables.Chat.ScrollRect)
		last_frame_scroll := up.app.Variables.PanelsVariables.Chat.LastFrameScroll

		up.app.Variables.PanelsVariables.Chat.LastFrameScroll = false

		if (scroll_hover && (clicked || down)) || (last_frame_scroll && down) {
			up.app.Variables.PanelsVariables.Chat.LastFrameScroll = true
			up.app.Variables.PanelsVariables.Chat.ScrollBarY = mouse.Y
		}
	}
}

func (up *Updater) UpdateScope() {

}
