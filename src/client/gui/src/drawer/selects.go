package drawer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawGameSelects() {
	if !dr.app.Variables.PanelsVariables.Group.Open {
		return
	}

	fontSize := int32(dr.app.Variables.FontSize)
	color := dr.app.Colors["panel_text"]

	for _, s := range dr.app.Manager.Selects("Group") {
		if len(s.Options) == 0 || s.CurrentOption == "" {
			continue
		}

		// Main box (closed state)
		main := s.MainRect()
		frame := s.CurrentMainFrame()
		dr.DrawImage(
			s.Texture,
			s.X, s.Y,
			frame.IndX, frame.IndY,
			frame.RatioX, frame.RatioY,
			s.Zoom, s.Rotation,
		)

		// Current option, centered in the main box
		textW := rl.MeasureText(s.CurrentOption, fontSize)
		rl.DrawText(
			s.CurrentOption,
			int32(main.X+(main.Width-float32(textW))/2),
			int32(main.Y+(main.Height-float32(fontSize))/2),
			fontSize, color,
		)

		// Options list (open state)
		if !s.Open {
			continue
		}
		for i, opt := range s.Options {
			item := s.ItemRect(i)
			iframe := s.ItemFrame(i)
			dr.DrawImage(
				s.Texture,
				item.X, item.Y,
				iframe.IndX, iframe.IndY,
				iframe.RatioX, iframe.RatioY,
				s.Zoom, s.Rotation,
			)

			tw := rl.MeasureText(opt, fontSize)
			rl.DrawText(
				opt,
				int32(item.X+(item.Width-float32(tw))/2),
				int32(item.Y+(item.Height-float32(fontSize))/2),
				fontSize, color,
			)
		}
	}
}
