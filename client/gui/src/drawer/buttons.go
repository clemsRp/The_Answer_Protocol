package drawer

import (
	"tap/client/gui/src/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawConnectButtons() {
	for _, b := range dr.app.Manager.Buttons("Connect") {
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

func (dr *Drawer) DrawGameButtons() {
	for _, b := range dr.app.Manager.Buttons("Game") {
		frame := b.CurrentFrame()

		dr.DrawImage(
			b.Texture,
			b.X,
			b.Y,
			frame.IndX,
			frame.IndY,
			frame.RatioX,
			frame.RatioY,
			b.Zoom,
			b.Rotation,
		)

		dr.DrawPanelNotifications(b)

	}
}

func (dr *Drawer) DrawChatButtons() {
	for _, b := range dr.app.Manager.Buttons("Chat") {
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

func (dr *Drawer) DrawInventoryButtons() {
	for _, b := range dr.app.Manager.Buttons("Inventory") {
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

func (dr *Drawer) DrawGroupButtons() {
	for _, b := range dr.app.Manager.Buttons("Group") {
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

func (dr *Drawer) DrawTalkButtons() {
	for _, b := range dr.app.Manager.Buttons("Talk") {
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

func (dr *Drawer) DrawInspectButtons() {
	for _, b := range dr.app.Manager.Buttons("Inspect") {
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

func (dr *Drawer) DrawCombatActionsButtons() {
	for _, b := range dr.app.Manager.Buttons("CombatActions") {
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

func (dr *Drawer) DrawCombatInventoryButtons() {
	for _, b := range dr.app.Manager.Buttons("CombatInventory") {
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

func (dr *Drawer) DrawDatasButtons() {
	for _, b := range dr.app.Manager.Buttons("Datas") {
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

func (dr *Drawer) DrawPanelNotifications(button *ui.Button) {
	group := dr.app.Variables.PanelsVariables.Group
	chat := dr.app.Variables.PanelsVariables.Chat

	var showNotification bool

	switch button.ID {
	case "open_chat":
		showNotification = !chat.Open && chat.Unread

	case "open_group":
		showNotification =
			!group.Open &&
				((!group.InGroup && len(group.Invitations) > 0) ||
					group.Promote)
	}

	if !showNotification {
		return
	}

	rect := button.Rect()

	radius := 0.10 * dr.app.Variables.Tileset_size

	x := rect.X + rect.Width - 1.75*radius
	y := rect.Y + 1.75*radius

	rl.DrawCircle(
		int32(x),
		int32(y),
		radius,
		rl.Red,
	)
}
