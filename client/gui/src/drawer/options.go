package drawer

import (
	"fmt"
	"tap/client/gui/src/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGroupOptions() {
	// Draw nothing if panel not open
	if !dr.app.Variables.PanelsVariables.Group.Open {
		return
	}

	for _, opt := range dr.app.Manager.Options("Group") {
		opt_name := dr.LimitString(opt.OptionName, 10)

		// Draw option name
		rl.DrawText(
			fmt.Sprintf("- %s:", opt_name),
			int32(opt.X),
			int32(opt.Y),
			int32(dr.app.Variables.FontSize),
			dr.app.Colors["pseudo_text"],
		)

		// Draw Decline/Accept button
		for _, b := range []*ui.Button{opt.DeclineBtn, opt.AcceptBtn} {
			if b == nil {
				continue
			}

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

}
