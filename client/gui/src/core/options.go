package core

import (
	"fmt"
	"strings"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
)

func (app *App) GetGroupOptions(command string) []*ui.Option {
	// Set group options datas
	group_options := map[string]*[]string{
		"Invite":  &app.Variables.PanelsVariables.Group.UnGrouped,
		"Join":    &app.Variables.PanelsVariables.Group.Invitations,
		"Kick":    &app.Variables.PanelsVariables.Group.Grouped,
		"Promote": &app.Variables.PanelsVariables.Group.Grouped,
	}

	group_options_funcs := map[string]func(command, command_option string){
		"Invite":  app.InviteFunc,
		"Join":    app.JoinFunc,
		"Promote": app.PromoteFunc,
		"Kick":    app.KickFunc,
	}

	start_x := 4.25 * app.Variables.Tileset_size
	start_y := 7.1 * app.Variables.Tileset_size
	end_x := (vars.GROUP_START_X + vars.GROUP_WIDTH - 1.3) * app.Variables.Tileset_size

	// Create options
	options := make([]*ui.Option, 0)

	target_options, ok := group_options[command]
	if !ok || target_options == nil {
		return options
	}

	for index, command_option := range *group_options[command] {
		jump_line := float32(index) * 1.5 * app.Variables.FontSize

		// Declare accept button
		accept_btn := &ui.Button{
			ID:       fmt.Sprintf("accept_btn-%s-%s", command, command_option),
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        end_x,
			Y:        start_y + jump_line - 0.18*app.Variables.Tileset_size,
			Zoom:     app.Variables.Zoom * 0.75,
			Rotation: 0,
			Normal:   ui.Frame{IndX: 15, IndY: 4, RatioX: 1, RatioY: 1},
			Hover:    ui.Frame{IndX: 15, IndY: 4, RatioX: 1, RatioY: 1},
			Pressed:  ui.Frame{IndX: 16, IndY: 4, RatioX: 1, RatioY: 1},
			OnClick: func() {
				group_options_funcs[command](command, command_option)

				app.ActionsChan <- panel.Action{
					Type:    panel.ActionSendServer,
					Payload: fmt.Sprintf("GROUP %s %s", strings.ToUpper(command), command_option),
				}
			},
		}

		// Create Option
		option := &ui.Option{
			OptionName: command_option,
			X:          int(start_x),
			Y:          int(start_y + jump_line),
			DeclineBtn: nil,
			AcceptBtn:  accept_btn,
		}

		options = append(options, option)
	}

	return options
}

func (app *App) InviteFunc(command, command_option string) {
	gr := app.Variables.PanelsVariables.Group

	new_ungrouped := make([]string, 0)
	for _, player := range gr.UnGrouped {
		if player != command_option {
			new_ungrouped = append(new_ungrouped, player)
		}
	}

	gr.UnGrouped = new_ungrouped
}

func (app *App) JoinFunc(command, command_option string) {
	app.Variables.PanelsVariables.Group.Leader = command_option
}

func (app *App) PromoteFunc(command, command_option string) {

}

func (app *App) KickFunc(command, command_option string) {

}
