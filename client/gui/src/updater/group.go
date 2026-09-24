package updater

import rl "github.com/gen2brain/raylib-go/raylib"

func (up *Updater) UpdateGroupScroll() {
	if len(up.app.Manager.Options("Group")) < 5 {
		up.app.Variables.PanelsVariables.Group.ScrollBarY = 0
		return
	}

	mouse := rl.GetMousePosition()
	hover := rl.CheckCollisionPointRec(mouse, up.app.Variables.PanelsVariables.Group.Rect)

	scroll := rl.GetMouseWheelMove()
	scroll_speed := 35

	if hover {
		up.app.Variables.PanelsVariables.Group.Scroll += scroll * float32(scroll_speed)
	}

	if up.app.Variables.PanelsVariables.Group.ScrollActive {
		clicked := rl.IsMouseButtonPressed(rl.MouseButtonLeft)
		down := rl.IsMouseButtonDown(rl.MouseButtonLeft)
		scroll_hover := rl.CheckCollisionPointRec(mouse, up.app.Variables.PanelsVariables.Group.ScrollRect)
		last_frame_scroll := up.app.Variables.PanelsVariables.Group.LastFrameScroll

		up.app.Variables.PanelsVariables.Group.LastFrameScroll = false

		if (scroll_hover && (clicked || down)) || (last_frame_scroll && down) {
			up.app.Variables.PanelsVariables.Group.LastFrameScroll = true
			up.app.Variables.PanelsVariables.Group.ScrollBarY = mouse.Y
		}
	}
}
