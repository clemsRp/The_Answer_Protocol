package drawer

// import (
// 	"tap/client/gui/src/ui"
// 	vars "tap/client/gui/src/variables"
// )

// func (dr *Drawer) buildConnectEmotes() {
// 	connect_start_x := 8.0
// 	connect_start_y := 4.5
// 	connect_width := 16.0
// 	connect_height := 9.0

// 	emote_x := float32(connect_start_x+connect_width/2) * dr.app.Variables.Tileset_size
// 	emote_y := float32(connect_start_y+connect_height) * dr.app.Variables.Tileset_size

// 	lengths := []int{1, 2, 5, 4, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1}
// 	emotes := make([]*ui.Emote, 0)

// 	for ind_y, ind_x_max := range lengths {
// 		emote_frames := make([]*ui.EmoteFrame, 0)

// 		for ind_x := range ind_x_max + 1 {
// 			frame := ui.Frame{IndX: float32(ind_x), IndY: float32(ind_y), RatioX: 1, RatioY: 1}
// 			emote_frame := ui.EmoteFrame{
// 				Frame:    &frame,
// 				Duration: 100,
// 			}

// 			emote_frames = append(emote_frames, &emote_frame)
// 		}

// 		emote := ui.Emote{
// 			ID:       "emote" + string(ind_y),
// 			Texture:  vars.EMOTES_TEXTURE,
// 			X:        emote_x,
// 			Y:        emote_y,
// 			Zoom:     dr.app.Variables.Zoom,
// 			Rotation: 0,
// 			Frames:   emote_frames,
// 		}

// 		emotes = append(emotes, &emote)
// 	}

// 	dr.app.Manager.SetViewEmotes("Connect", emotes)
// }

// func (dr *Drawer) DrawEmotes(view string) {
// 	for _, e := range dr.app.Manager.Emotes(view) {
// 		frame := e.CurrentFrame()
// 		dr.DrawImage(
// 			e.Texture,
// 			e.X, e.Y,
// 			frame.IndX, frame.IndY,
// 			frame.RatioX, frame.RatioY,
// 			e.Zoom,
// 			e.Rotation,
// 		)
// 	}
// }
