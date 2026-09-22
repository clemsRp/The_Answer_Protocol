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
	group_selects = SelectDatas{
		X: 2, Y: 7, Options: make([]string, 0),
	}
)

func (app *App) GetGroupSelects() []*ui.Select {
	options := make([]string, 0)
	gr := app.Variables.PanelsVariables.Group

	if len(gr.UnGrouped) > 0 && gr.InGroup && gr.Leader == app.Variables.Player.Pseudo {
		options = append(options, "Invite")
	}
	if len(gr.Grouped) > 0 && gr.InGroup && gr.Leader == app.Variables.Player.Pseudo {
		options = append(options, "Kick")
		options = append(options, "Promote")
	}
	if len(gr.Invitations) > 0 && !gr.InGroup {
		options = append(options, "Join")
	}

	group_selects.Options = options

	cur_option := ""
	if len(group_selects.Options) > 0 {
		cur_option = group_selects.Options[0]
	}

	// Create Select
	select_element := &ui.Select{
		ID:            "group_options_select",
		Texture:       vars.UI_SPRITE_TEXTURE,
		X:             group_selects.X * app.Variables.Tileset_size,
		Y:             group_selects.Y * app.Variables.Tileset_size,
		Zoom:          app.Variables.Zoom / 3,
		Rotation:      0,
		Normal:        ui.Frame{IndX: 10, IndY: 11, RatioX: 6, RatioY: 2},
		Hover:         ui.Frame{IndX: 16, IndY: 11, RatioX: 6, RatioY: 2},
		ItemNormal:    ui.Frame{IndX: 10, IndY: 11, RatioX: 6, RatioY: 2},
		ItemHover:     ui.Frame{IndX: 16, IndY: 11, RatioX: 6, RatioY: 2},
		Options:       group_selects.Options,
		CurrentOption: cur_option,
		Open:          false,
		OnChange: func(option string) {
			app.Manager.SetViewOptions("Group", app.GetGroupOptions(option))
		},
	}

	return []*ui.Select{select_element}
}
