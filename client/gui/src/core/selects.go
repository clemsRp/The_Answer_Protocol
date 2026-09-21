package core

import (
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
)

type SelectDatas struct {
	X       float32
	Y       float32
	Options []string
}

var (
	group_selects = []SelectDatas{
		{X: 2, Y: 6, Options: []string{"Invite", "Promote", "Kick"}},
	}
)

func (app *App) GetGroupSelects() []*ui.Select {
	selects_elements := make([]*ui.Select, 0)

	for _, s := range group_selects {
		// Create Select
		select_element := &ui.Select{
			ID:            "group_options_select",
			Texture:       vars.UI_SPRITE_TEXTURE,
			X:             s.X * app.Variables.Tileset_size,
			Y:             s.Y * app.Variables.Tileset_size,
			Zoom:          app.Variables.Zoom / 3,
			Rotation:      0,
			Normal:        ui.Frame{IndX: 10, IndY: 11, RatioX: 6, RatioY: 2},
			Hover:         ui.Frame{IndX: 16, IndY: 11, RatioX: 6, RatioY: 2},
			ItemNormal:    ui.Frame{IndX: 10, IndY: 11, RatioX: 6, RatioY: 2},
			ItemHover:     ui.Frame{IndX: 16, IndY: 11, RatioX: 6, RatioY: 2},
			Options:       s.Options,
			CurrentOption: s.Options[0],
			Open:          false,
			OnChange: func(option string) {
				app.Manager.SetViewOptions("Group", app.GetGroupOptions(option))
			},
		}

		selects_elements = append(selects_elements, select_element)
	}

	return selects_elements
}
