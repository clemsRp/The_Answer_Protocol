package drawer

import (
	"fmt"
	"tap/client/gui/src/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGroupOptions(group_start_x, group_start_y, group_height float32) {
	if !dr.app.Variables.PanelsVariables.Group.Open {
		return
	}

	offsetX := group_start_x * dr.app.Variables.Tileset_size
	offsetY := dr.app.GroupContentStartY()

	for _, opt := range dr.app.Manager.Options("Group") {
		if opt.OptionName == "promote_option" && !dr.app.Variables.PanelsVariables.Group.Promote {
			continue
		}

		opt_name := dr.LimitString(opt.OptionName, 8)

		rl.DrawText(
			fmt.Sprintf("- %s:", opt_name),
			int32(opt.X)-int32(offsetX),
			int32(opt.Y)-int32(offsetY),
			int32(dr.app.Variables.FontSize),
			dr.app.Colors["pseudo_text"],
		)

		for _, b := range []*ui.Button{opt.DeclineBtn, opt.AcceptBtn} {
			if b == nil {
				continue
			}

			frame := b.CurrentFrame()

			dr.DrawImage(
				b.Texture,
				b.X-offsetX, b.Y-offsetY,
				frame.IndX, frame.IndY,
				frame.RatioX, frame.RatioY,
				b.Zoom,
				b.Rotation,
			)
		}
	}
}

func (dr *Drawer) DrawGroupOptionsToScreen(group_start_x, group_start_y, group_height float32) {
	tex := dr.groupTexture.Texture
	tile := dr.app.Variables.Tileset_size

	offsetX := group_start_x * tile
	// On découpe à partir du bas de la ligne "Members", pas du haut du panel.
	offsetY := dr.app.GroupContentStartY()
	viewHeight := (group_start_y+group_height)*tile - offsetY - 0.5*tile

	displayHeight := viewHeight
	maxY := float32(tex.Height) - displayHeight
	minY := float32(0)

	if float32(tex.Height) < viewHeight {
		displayHeight = float32(tex.Height)
		maxY = 0
		dr.app.Variables.PanelsVariables.Group.ScrollActive = false
	} else if len(dr.app.Manager.Options("Group")) >= 5 {
		dr.DrawGroupScrollBar(maxY)
		dr.app.Variables.PanelsVariables.Group.ScrollActive = true
	}

	sourceY := maxY
	scroll := dr.app.Variables.PanelsVariables.Group.Scroll
	finalY := sourceY + scroll

	if finalY > maxY {
		finalY = maxY
		dr.app.Variables.PanelsVariables.Group.Scroll = finalY - sourceY
	} else if finalY < minY {
		finalY = minY
		dr.app.Variables.PanelsVariables.Group.Scroll = finalY - sourceY
	}

	sourceRec := rl.NewRectangle(0, finalY, float32(tex.Width), -displayHeight)
	position := rl.NewVector2(offsetX, offsetY)

	rl.DrawTextureRec(tex, sourceRec, position, rl.White)
}
