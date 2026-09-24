package core

import (
	"fmt"
	"slices"
	"strings"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

func (app *App) GetGroupOptions(command string) []*ui.Option {
	options := make([]*ui.Option, 0)

	// Add the current tab's entries
	options = append(options, app.buildTabOptions(command)...)

	// Add the accept/decline promotion prompt if one is pending
	if app.Variables.PanelsVariables.Group.Promote {
		options = append(options, app.buildPromoteOption())
	}

	return options
}

func (app *App) buildTabOptions(command string) []*ui.Option {
	options := make([]*ui.Option, 0)

	// Map each tab name to its underlying data slice
	group_options := map[string]*[]string{
		"Invite":  &app.Variables.PanelsVariables.Group.UnGrouped,
		"Join":    &app.Variables.PanelsVariables.Group.Invitations,
		"Kick":    &app.Variables.PanelsVariables.Group.Grouped,
		"Promote": &app.Variables.PanelsVariables.Group.Grouped,
	}

	// Bail out if the tab is unknown
	target_options, ok := group_options[command]
	if !ok || target_options == nil {
		return options
	}

	start_x := 4.25 * app.Variables.Tileset_size
	start_y := 7.1 * app.Variables.Tileset_size
	end_x := (vars.GROUP_START_X + vars.GROUP_WIDTH - 1.6) * app.Variables.Tileset_size

	for index, command_option := range *target_options {
		// Skip an entry already handled locally
		if app.shouldSkipTabOption(command, command_option) {
			continue
		}

		jump_line := float32(index) * 1.5 * app.Variables.FontSize
		accept_btn := app.buildTabAcceptButton(command, command_option, end_x, start_y+jump_line)

		// Create the row
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

func (app *App) shouldSkipTabOption(command, command_option string) bool {
	gr := app.Variables.PanelsVariables.Group

	// Skip a player already invited to join
	if command == "Join" && slices.Contains(gr.SendInvitations, command_option) {
		return true
	}

	// Skip promoting while another promotion is already pending
	if command == "Promote" && gr.SendPromotion != "" {
		return true
	}

	return false
}

func (app *App) buildTabAcceptButton(command, command_option string, x, y float32) *ui.Button {
	return &ui.Button{
		ID:       fmt.Sprintf("accept_btn-%s-%s", command, command_option),
		Texture:  vars.UI_SPRITE_TEXTURE,
		X:        x,
		Y:        y - 0.18*app.Variables.Tileset_size,
		Zoom:     app.Variables.Zoom * 0.75,
		Rotation: 0,
		Normal:   ui.Frame{IndX: 15, IndY: 4, RatioX: 1, RatioY: 1},
		Hover:    ui.Frame{IndX: 15, IndY: 4, RatioX: 1, RatioY: 1},
		Pressed:  ui.Frame{IndX: 16, IndY: 4, RatioX: 1, RatioY: 1},
		OnClick: func() {
			// Run the local side effect for this command
			app.groupCommandFuncs()[command](command, command_option)

			// Refresh the select and options after the local change
			app.RefreshGroupView(command)

			// Notify the server of the action
			app.ActionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: fmt.Sprintf("GROUP %s %s", strings.ToUpper(command), command_option),
			}
		},
	}
}

func (app *App) groupCommandFuncs() map[string]func(command, command_option string) {
	return map[string]func(command, command_option string){
		"Invite":  app.InviteFunc,
		"Join":    app.JoinFunc,
		"Promote": app.PromoteFunc,
		"Kick":    app.KickFunc,
	}
}

func (app *App) buildPromoteOption() *ui.Option {
	promote_x := (vars.GROUP_START_X + 0.7) * app.Variables.Tileset_size
	promote_y := (vars.GROUP_START_Y+2)*app.Variables.Tileset_size + 2.5*app.Variables.FontSize
	end_x := (vars.GROUP_START_X + vars.GROUP_WIDTH - 1.3) * app.Variables.Tileset_size

	decline_btn := app.buildPromoteButton("decline_btn-promotion", end_x-1.2*app.Variables.Tileset_size, promote_y, 5, pr.CmdDeclinePromoteGroup)
	accept_btn := app.buildPromoteButton("accept_btn-promotion", end_x, promote_y, 4, pr.CmdAcceptPromoteGroup)

	return &ui.Option{
		OptionName: "promote_option",
		X:          int(promote_x),
		Y:          int(promote_y),
		DeclineBtn: decline_btn,
		AcceptBtn:  accept_btn,
	}
}

func (app *App) buildPromoteButton(id string, x, y float32, frame_ind_y float32, payload string) *ui.Button {
	return &ui.Button{
		ID:       id,
		Texture:  vars.UI_SPRITE_TEXTURE,
		X:        x,
		Y:        y - 0.18*app.Variables.Tileset_size,
		Zoom:     app.Variables.Zoom * 0.75,
		Rotation: 0,
		Normal:   ui.Frame{IndX: 15, IndY: frame_ind_y, RatioX: 1, RatioY: 1},
		Hover:    ui.Frame{IndX: 15, IndY: frame_ind_y, RatioX: 1, RatioY: 1},
		Pressed:  ui.Frame{IndX: 16, IndY: frame_ind_y, RatioX: 1, RatioY: 1},
		OnClick: func() {
			if payload == pr.CmdAcceptPromoteGroup {
				app.Variables.PanelsVariables.Group.Leader = app.GetPseudo()
			}
			app.Variables.PanelsVariables.Group.Promote = false
			app.RefreshGroupView("")
			app.ActionsChan <- panel.Action{
				Type:    panel.ActionSendServer,
				Payload: payload,
			}
		},
	}
}

func (app *App) RefreshGroupView(command string) {
	selects := app.GetGroupSelects()
	next_current := ""

	if len(selects) > 0 {
		// Keep the current tab if still valid, else fall back to the first one
		if slices.Contains(selects[0].Options, command) {
			next_current = command
		} else if len(selects[0].Options) > 0 {
			next_current = selects[0].Options[0]
		}
		selects[0].CurrentOption = next_current
	}

	app.Manager.SetViewSelects("Group", selects)
	app.Manager.SetViewOptions("Group", app.GetGroupOptions(next_current))
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
	app.Variables.PanelsVariables.Group.SendPromotion = command_option
}

func (app *App) KickFunc(command, command_option string) {
	if app.Variables.PanelsVariables.Group.SendPromotion == command_option {
		app.Variables.PanelsVariables.Group.SendPromotion = ""
	}
}
