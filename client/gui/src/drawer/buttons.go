package drawer

import (
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

var nb_click = 0

func (dr *Drawer) buildConnectButtons() {
	connect_start_x := 8.0
	connect_start_y := 4.5
	connect_width := 16.0
	connect_height := 9.0

	play_x := float32(connect_start_x+connect_width/2) * dr.app.Variables.Tileset_size
	play_y := float32(connect_start_y+connect_height) * dr.app.Variables.Tileset_size

	center_x := 3 * dr.app.Variables.Tileset_size
	center_y := 2.5 * dr.app.Variables.Tileset_size

	playBtn := &ui.Button{
		ID:      "play",
		Texture: vars.PLAY_TEXTURE,
		X:       play_x - center_x,
		Y:       play_y - center_y,
		Zoom:    dr.app.Variables.Zoom,
		Normal:  ui.Frame{IndX: 0, IndY: 2, RatioX: 6, RatioY: 2},
		Pressed: ui.Frame{IndX: 6, IndY: 2, RatioX: 6, RatioY: 2},
		OnClick: func() {
			nb_click++
			if nb_click < 5 {
				return
			}
			dr.app.ActionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: pr.CmdConnect + " clement",
			}
		},
	}

	dr.app.Manager.SetViewButtons("Connect", []*ui.Button{playBtn})
}

func (dr *Drawer) DrawButtons(view string) {
	for _, b := range dr.app.Manager.Buttons(view) {
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
