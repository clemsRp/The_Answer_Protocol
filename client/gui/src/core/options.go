package core

import (
	"fmt"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
)

var ()

func (app *App) GetGroupOptions(command string) []*ui.Option {
	// Set group options datas
	group_options := map[string]*[]string{
		"Invite":  &app.Variables.PanelsVariables.GroupState.UnGrouped,
		"Join":    &app.Variables.PanelsVariables.GroupState.Invitations,
		"Kick":    &app.Variables.PanelsVariables.GroupState.Grouped,
		"Promote": &app.Variables.PanelsVariables.GroupState.Grouped,
	}

	start_x := 2 * app.Variables.Tileset_size
	start_y := 12 * app.Variables.Tileset_size

	// Create options
	options := make([]*ui.Option, 0)

	for _, command_option := range *group_options[command] {
		// Declare accept button
		accept_btn := &ui.Button{
			ID:       fmt.Sprintf("accept_btn-%s-%s", command, command_option),
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        start_x,
			Y:        start_y,
			Zoom:     app.Variables.Zoom,
			Rotation: 0,
			Normal:   ui.Frame{IndX: 15, IndY: 4, RatioX: 1, RatioY: 1},
			Hover:    ui.Frame{IndX: 15, IndY: 4, RatioX: 1, RatioY: 1},
			Pressed:  ui.Frame{IndX: 16, IndY: 4, RatioX: 1, RatioY: 1},
			OnClick: func() {
				app.ActionsChan <- panel.Action{
					Type:    panel.ActionSendServer,
					Payload: fmt.Sprintf("%s %s", command, command_option),
				}
			},
		}

		// Create Option
		option := &ui.Option{
			OptionName: command_option,
			X:          int(start_x),
			Y:          int(start_y),
			DeclineBtn: nil,
			AcceptBtn:  accept_btn,
		}

		options = append(options, option)
	}

	return options
}
