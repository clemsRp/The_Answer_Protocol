package core

import (
	"fmt"
	"slices"
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
		// Skip already invited players
		if slices.Contains(app.Variables.PanelsVariables.Group.SendInvitations, command_option) && command == "Join" {
			continue
		}

		// Skip already invited players
		if app.Variables.PanelsVariables.Group.SendPromotion && command == "Promote" {
			continue
		}

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

				// Refresh select
				selects := app.GetGroupSelects()
				next_current := ""
				if len(selects) > 0 {
					valid := false
					for _, opt := range selects[0].Options {
						if opt == command {
							valid = true
							break
						}
					}
					if valid {
						next_current = command
					} else if len(selects[0].Options) > 0 {
						next_current = selects[0].Options[0]
					}
					selects[0].CurrentOption = next_current
				}
				app.Manager.SetViewSelects("Group", selects)
				app.Manager.SetViewOptions("Group", app.GetGroupOptions(next_current))

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

	app.Variables.PanelsVariables.Group.SendInvitations = append(app.Variables.PanelsVariables.Group.SendInvitations, command_option)
}

func (app *App) JoinFunc(command, command_option string) {
	app.Variables.PanelsVariables.Group.Leader = command_option
}

func (app *App) PromoteFunc(command, command_option string) {

}

func (app *App) KickFunc(command, command_option string) {

}
