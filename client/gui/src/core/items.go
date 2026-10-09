package core

import (
	"fmt"
	"strconv"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	panel "tap/client/tui/panels"
	"tap/protocol"
	pr "tap/protocol"
)

func (app *App) GetNewRoomItems(roomItems []string) []*ui.Interaction {
	items := make([]*ui.Interaction, 0)

	for _, it := range roomItems {
		frame_duration := 500

		text, frames := get_item_datas(it, frame_duration)

		var pos *vars.Position
		var ok bool
		if pos, ok = (*app.Variables.ItemPositions)[it]; !ok {
			if itemDef, exists := vars.ItemConvertor[it]; exists {
				p := &vars.Position{
					X: itemDef.Pos.X * app.Variables.Tileset_size,
					Y: itemDef.Pos.Y * app.Variables.Tileset_size,
				}
				(*app.Variables.ItemPositions)[it] = p
				pos = p
			} else {
				continue
			}
		}

		// Create the visual emote (sprite) for the item
		item_emote := &ui.Emote{
			ID:           "item_" + it,
			Texture:      text,
			X:            pos.X,
			Y:            pos.Y,
			Zoom:         app.Variables.Zoom,
			Rotation:     0,
			AnimDuration: len(frames) * frame_duration,
			Frames:       frames,
		}

		// Create the background bubble for the interaction
		item_bubble := &ui.Button{
			ID:       "item_bubble",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        pos.X - 0.2*app.Variables.Tileset_size,
			Y:        pos.Y - 1.2*app.Variables.Tileset_size,
			Rotation: -90,
			Zoom:     app.Variables.Zoom / 2,
			Normal:   ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			Pressed:  ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			Hover:    ui.Frame{IndX: 28, IndY: 8, RatioX: 3, RatioY: 3},
			OnClick:  func() {},
		}

		// Create the interactive button to pick the item
		item_btn := &ui.Button{
			ID:       "item_btn",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        pos.X + 0.15*app.Variables.Tileset_size,
			Y:        pos.Y - 0.85*app.Variables.Tileset_size,
			Rotation: -90,
			Zoom:     app.Variables.Zoom * 2 / 5,
			Normal:   ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
			Hover:    ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
			Pressed:  ui.Frame{IndX: 54, IndY: 8, RatioX: 2, RatioY: 2},
			OnClick: func() {
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdTake + " " + it}
			},
		}

		// Group the emote and buttons into a single interaction entity
		item := &ui.Interaction{
			ID:      it,
			Emote:   item_emote,
			Buttons: []*ui.Button{item_bubble, item_btn},
			Inspect: func(item string) {
				app.Variables.PanelsVariables.Inspect.LastInspect = "ITEM"
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdInspectItem + " " + item}
				app.Variables.PanelsVariables.Inspect.Open = true
			},
		}

		items = append(items, item)
	}

	return items
}

func (app *App) GetNewInventory(inventory []string) ([]*ui.Button, []*ui.Emote) {
	// Initialize inventory UI element arrays
	invent_buttons := make([]*ui.Button, 0)
	invent_emotes := make([]*ui.Emote, 0)

	pseudo := app.Variables.Player.Pseudo
	start_x := float32(19)

	tier := app.GetPseudoBoxTier(pseudo)
	if tier == 1 {
		start_x = 11
	} else if tier == 2 {
		start_x = 15
	}

	// Calculate final UI coordinates based on zoom and tileset scale
	start_x += 2
	start_x *= app.Variables.Zoom / 2 * vars.FRAME_WIDTH
	start_x -= 0.1 * app.Variables.Tileset_size
	start_y := app.Variables.Tileset_size

	for ind, it := range inventory {
		frame_duration := 500

		// Retrieve texture and frames for the current inventory item
		text, frames := get_item_datas(it, frame_duration)

		// Create the visual representation of the inventory item
		item_emote := &ui.Emote{
			ID:           it,
			Texture:      text,
			X:            start_x + app.Variables.Tileset_size*(float32(ind)+0.2),
			Y:            0.47*app.Variables.Tileset_size + start_y,
			Zoom:         app.Variables.Zoom / 1.5,
			Rotation:     0,
			AnimDuration: 2 * frame_duration,
			Frames:       frames,
			Inspect: func(item string) {
				app.Variables.PanelsVariables.Inspect.LastInspect = "ITEM"
				app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdInspectInventoryItem + " " + item}
				app.Variables.PanelsVariables.Inspect.Open = true
			},
		}

		// Create the use button for the inventory item
		if app.IsItemUsable(it) {
			item_use_btn := &ui.Button{
				ID:       "item_use_btn",
				Texture:  vars.UI_SPRITE_TEXTURE,
				X:        start_x + app.Variables.Tileset_size*float32(ind),
				Y:        0.2*app.Variables.Tileset_size + start_y,
				Rotation: 0,
				Zoom:     app.Variables.Zoom / 5,
				Normal:   ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
				Hover:    ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
				Pressed:  ui.Frame{IndX: 54, IndY: 8, RatioX: 2, RatioY: 2},
				OnClick: func() {
					app.QueueUpdate(func() {
						app.ActionsChan <- panel.Action{
							Type:    panel.ActionSendServer,
							Payload: pr.CmdUseItem + " " + it,
						}
						if app.IsItemWeapon(it) {
							app.Variables.Player.EquipedItems = append(app.Variables.Player.EquipedItems, it)
						}
					})
				},
			}
			invent_buttons = append(invent_buttons, item_use_btn)
		}

		// Create the drop button for the inventory item
		item_drop_btn := &ui.Button{
			ID:       "item_drop_btn",
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        start_x + app.Variables.Tileset_size*(float32(ind)+0.645),
			Y:        0.2*app.Variables.Tileset_size + start_y,
			Rotation: 0,
			Zoom:     app.Variables.Zoom / 5,
			Normal:   ui.Frame{IndX: 52, IndY: 10, RatioX: 2, RatioY: 2},
			Hover:    ui.Frame{IndX: 52, IndY: 10, RatioX: 2, RatioY: 2},
			Pressed:  ui.Frame{IndX: 54, IndY: 10, RatioX: 2, RatioY: 2},
			OnClick: func() {
				app.DropItem(it)
			},
		}

		invent_buttons = append(invent_buttons, item_drop_btn)
		invent_emotes = append(invent_emotes, item_emote)
	}

	return invent_buttons, invent_emotes
}

func (app *App) GetNewCombatInventory(inventory []string) []*ui.Button {
	use_buttons := make([]*ui.Button, 0)

	pseudo := app.Variables.Player.Pseudo
	start_x := float32(19)

	tier := app.GetPseudoBoxTier(pseudo)
	if tier == 1 {
		start_x = 11
	} else if tier == 2 {
		start_x = 15
	}

	start_x += 2
	start_x *= app.Variables.Zoom / 2 * vars.FRAME_WIDTH
	start_x -= 0.1 * app.Variables.Tileset_size
	start_y := app.Variables.Tileset_size

	usableIndex := 0
	for ind, it := range inventory {
		if !app.IsItemUsable(it) {
			continue
		}
		itemName := it
		item_btn := &ui.Button{
			ID:       "combat_use_" + itemName + "_" + strconv.Itoa(ind),
			Texture:  vars.UI_SPRITE_TEXTURE,
			X:        start_x + app.Variables.Tileset_size*(float32(ind)+0.645),
			Y:        0.2*app.Variables.Tileset_size + start_y,
			Rotation: 0,
			Zoom:     app.Variables.Zoom / 5,
			Normal:   ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
			Hover:    ui.Frame{IndX: 52, IndY: 8, RatioX: 2, RatioY: 2},
			Pressed:  ui.Frame{IndX: 54, IndY: 8, RatioX: 2, RatioY: 2},
			OnClick: func() {
				app.ActionsChan <- panel.Action{
					Type:    panel.ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdUseItem, itemName),
				}
				if app.IsItemWeapon(it) {
					app.Variables.Player.EquipedItems = append(app.Variables.Player.EquipedItems, it)
				}
			},
		}
		use_buttons = append(use_buttons, item_btn)
		usableIndex++
	}

	return use_buttons
}

func (app *App) DropItem(item string) {
	// Get offset datas
	dirX := app.Variables.Player.Direction.X
	dirY := app.Variables.Player.Direction.Y
	if dirX == 0 && dirY == 0 {
		dirY = 1
	}

	tileset_size := app.Variables.Tileset_size

	// Find closest free tile around the player
	newPosX, newPosY, found := app.FindNearestFreeTile(
		app.Variables.Player.Position.X+dirX*tileset_size,
		app.Variables.Player.Position.Y+dirY*tileset_size,
	)
	if !found {
		return
	}

	// Update position
	if _, exist := (*app.Variables.ItemPositions)[item]; !exist {
		(*app.Variables.ItemPositions)[item] = &vars.Position{
			X: app.Variables.Player.Position.X,
			Y: app.Variables.Player.Position.Y,
		}
	}

	(*app.Variables.ItemPositions)[item].X = newPosX
	(*app.Variables.ItemPositions)[item].Y = newPosY

	newTileX := newPosX / app.Variables.Tileset_size
	newTileY := newPosY / app.Variables.Tileset_size

	payload := fmt.Sprintf("%s %f %f %f %f %s", protocol.CmdNotifyItemPosition, newPosX, newPosY, newTileX, newTileY, item)

	// Send notifs
	app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: pr.CmdDrop + " " + item}
	app.ActionsChan <- panel.Action{Type: panel.ActionSendServer, Payload: payload}
}

func (app *App) FindNearestFreeTile(startX, startY float32) (float32, float32, bool) {
	tileSize := int(app.Variables.Tileset_size)
	if tileSize <= 0 {
		return startX, startY, false
	}

	room := app.Rooms[app.Variables.Current_room]
	if room == nil || len(room.Collisions) == 0 || len(room.Collisions[0]) == 0 {
		return startX, startY, false
	}

	// Get starting tile datas
	startTileX := int(startX) / tileSize
	startTileY := int(startY) / tileSize
	maxRadius := len(room.Collisions) + len(room.Collisions[0])

	// Expand ring by ring around the starting tile
	for radius := 1; radius <= maxRadius; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				// Only check the ring border, skip cells already checked at smaller radius
				if radius > 0 && absInt(dx) != radius && absInt(dy) != radius {
					continue
				}

				tileX := startTileX + dx
				tileY := startTileY + dy

				// Skip tiles outside the map
				if tileY < 0 || tileY >= len(room.Collisions) || tileX < 0 || tileX >= len(room.Collisions[0]) {
					continue
				}

				// Skip tiles with a collision
				if room.Collisions[tileY][tileX] != 0 {
					continue
				}

				worldX := float32(tileX * tileSize)
				worldY := float32(tileY * tileSize)

				// Skip tiles already holding an item
				if app.isTileOccupiedByItem(worldX, worldY) {
					continue
				}

				return worldX, worldY, true
			}
		}
	}

	return startX, startY, false
}

func (app *App) isTileOccupiedByItem(x, y float32) bool {
	for _, pos := range *app.Variables.ItemPositions {
		if pos.X == x && pos.Y == y {
			return true
		}
	}
	return false
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func get_item_datas(item string, frame_duration int) (string, []*ui.EmoteFrame) {
	frames := make([]*ui.EmoteFrame, 0)

	// Define needed datas
	var texture string
	var ratioX, ratioY float32
	var indXs, indYs []float32

	// Get datas depending on item
	if d, ok := vars.ItemConvertor[item]; ok {
		texture = d.Texture
		indXs = d.IndXs
		indYs = d.IndYs
		ratioX = d.RatioX
		ratioY = d.RatioY

	} else {
		texture = vars.UI_SPRITE_TEXTURE
		indXs = []float32{17, 18}
		indYs = []float32{0, 0}
		ratioX = 1
		ratioY = 1
	}

	// Create emote frames
	for index := range indXs {
		frame := &ui.EmoteFrame{
			Frame:    &ui.Frame{IndX: indXs[index], IndY: indYs[index], RatioX: ratioX, RatioY: ratioY},
			Duration: frame_duration,
		}

		frames = append(frames, frame)
	}

	return texture, frames
}
