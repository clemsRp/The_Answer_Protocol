package tui

import (
	"encoding/json"
	"strings"
	panel "tap/client/tui/panels"
	pr "tap/protocol"
)

var (
	// Items
	items = make([]string, 0)

	// Interactions
	npcs    = make([]string, 0)
	players = make([]string, 0)

	// Group
	group       = ""
	invitations = make([]string, 0)
)

func (m *MyApp) NavListenOutputs(res pr.ServerResponse) {
	// Remove previous nav panel
	m.grid.RemoveItem(m.Navigation.Layout)

	// Get new rooms
	opts := make(map[string]string)

	raw, err := json.Marshal(res.Datas)
	if err == nil {
		var data pr.LookCommandData
		if err := json.Unmarshal(raw, &data); err == nil {
			exits := data.Room.Exits
			if exits.North != "" {
				opts["north"] = exits.North
			}
			if exits.East != "" {
				opts["east"] = exits.East
			}
			if exits.West != "" {
				opts["west"] = exits.West
			}
			if exits.South != "" {
				opts["south"] = exits.South
			}
		}
	}

	// Create new panel
	m.Navigation = panel.NewNavigationComponent(
		m.app,
		m.popup,
		opts,
		m.router.Inputs,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	// Add new panel
	m.grid.AddItem(m.Navigation.Layout, 0, 0, 1, 1, 0, 0, false)
	m.setupMatrix()
	m.ShowGamePage()
}

func (m *MyApp) ItemListenOutputs(res pr.ServerResponse) {
	// Remove previous item panel
	m.grid.RemoveItem(m.Items.Layout)

	var opts []string

	if strings.HasPrefix(res.Msg, "EVT ITEM DROPPED") {
		new_item := strings.SplitN(res.Msg, "DROPPED ", 2)[1]
		items = append(items, new_item)
		opts = items

	} else if strings.HasPrefix(res.Msg, "EVT ITEM TOOK") {
		previous_item := strings.SplitN(res.Msg, "TOOK ", 2)[1]

		for i, p := range items {
			if p == previous_item {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}

		opts = items

	} else {
		// Get new items
		raw, err := json.Marshal(res.Datas)
		if err == nil {
			var data pr.LookCommandData
			if err := json.Unmarshal(raw, &data); err == nil {
				opts = data.Room.Items
			}
		}
	}

	// Create new panel
	m.Items = panel.NewItemsComponent(
		m.app,
		m.popup,
		opts,
		m.router.Inputs,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	// Add new panel
	m.grid.AddItem(m.Items.Layout, 0, 1, 2, 1, 0, 0, false)
	m.setupMatrix()
	m.ShowGamePage()
}

func (m *MyApp) InteractionListenOutputs(res pr.ServerResponse) {
	// Remove previous item panel
	m.grid.RemoveItem(m.Interaction.Layout)

	if strings.HasPrefix(res.Msg, "EVT ROOM PRESENCE ENTER") {
		new_player := strings.SplitN(res.Msg, "ENTER ", 2)[1]
		players = append(players, new_player)

	} else if strings.HasPrefix(res.Msg, "EVT ROOM PRESENCE LEAVE") {
		previous_player := strings.SplitN(res.Msg, "LEAVE ", 2)[1]

		for i, p := range players {
			if p == previous_player {
				players = append(players[:i], players[i+1:]...)
				break
			}
		}

	} else {
		// Get new items
		npcs = []string{}
		players = []string{}
		raw, err := json.Marshal(res.Datas)
		if err == nil {
			var data pr.LookCommandData
			if err := json.Unmarshal(raw, &data); err == nil {
				npcs = data.Room.Npcs
				for _, player := range data.Room.Players {
					if player != m.pseudo {
						players = append(players, player)
					}
				}
			}
		}
	}

	// Create new panel
	m.Interaction = panel.NewInteractionComponent(
		m.app,
		m.popup,
		npcs,
		players,
		m.router.Inputs,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	// Add new panel
	m.grid.AddItem(m.Interaction.Layout, 0, 2, 2, 1, 0, 0, false)
	m.setupMatrix()
	m.ShowGamePage()
}

func (m *MyApp) GroupListenOutputs(res pr.ServerResponse) {
	// Remove previous item panel
	m.grid.RemoveItem(m.Group.Layout)

	if strings.HasPrefix(res.Msg, "EVT GROUP INVITE") {
		invitation := strings.SplitN(res.Msg, "INVITE ", 2)[1]
		invitations = append(invitations, invitation)

	} else if strings.HasPrefix(res.Msg, "OK group=") {
		group = strings.Split(res.Msg, "group=")[1]
	}

	// Create new panel
	m.Group = panel.NewGroupComponent(
		m.app,
		m.popup,
		group,
		&invitations,
		m.router.Inputs,
		m.OnOpenPopup,
		m.ShowGamePage,
	)

	// Add new panel
	m.grid.AddItem(m.Group.Layout, 1, 0, 1, 1, 0, 0, false)
	m.setupMatrix()
	m.ShowGamePage()
}
