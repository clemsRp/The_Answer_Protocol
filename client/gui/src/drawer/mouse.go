package drawer

import vars "tap/client/gui/src/variables"

func (dr *Drawer) DrawMouse() {
	indX := 21
	indY := 8

	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		dr.app.Variables.Mouse.X,
		dr.app.Variables.Mouse.Y,
		float32(indX), float32(indY), 1, 1,
		dr.app.Variables.Zoom/2, 0,
	)
}
