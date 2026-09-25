package drawer

import (
	"sort"
	"strings"
	"tap/client/gui/src/ui"
	vars "tap/client/gui/src/variables"
	pr "tap/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type fightersPanelGeometry struct {
	tile          float32
	zoom          float32
	frameSize     int
	baseFrameSize float32
	curTurn       string

	panelStartX float32
	panelStartY float32
	panelWidth  float32
	panelHeight float32
}

func turnZoomFor(name, curTurn string) float32 {
	if curTurn != "" && curTurn == name {
		return 1.25
	}
	return 1
}

func (dr *Drawer) ComputeFightersPanelGeometry() fightersPanelGeometry {
	tile := dr.app.Variables.Tileset_size
	zoom := dr.app.Variables.Zoom
	frameSize := 4

	panelStartX := vars.COMBAT_START_X * tile
	panelStartY := vars.COMBAT_START_Y * tile
	panelWidth := (vars.COMBAT_LEFT_END_X - vars.COMBAT_START_X) * tile
	panelHeight := (vars.COMBAT_END_Y - vars.COMBAT_START_Y) * tile

	return fightersPanelGeometry{
		tile:          tile,
		zoom:          zoom,
		frameSize:     frameSize,
		baseFrameSize: 0.5 * float32(frameSize) * vars.FRAME_WIDTH * zoom,
		curTurn:       dr.app.Variables.PanelsVariables.CombatState.CurrentTurn,
		panelStartX:   panelStartX,
		panelStartY:   panelStartY,
		panelWidth:    panelWidth,
		panelHeight:   panelHeight,
	}
}

func (dr *Drawer) DrawCombatFightersZone() {
	dr.DrawCombatFrame(vars.COMBAT_START_X, vars.COMBAT_START_Y, vars.COMBAT_LEFT_END_X, vars.COMBAT_FIGHTERS_END_Y)

	geo := dr.ComputeFightersPanelGeometry()
	state := dr.app.Variables.PanelsVariables.CombatState

	oppBaseY := geo.panelStartY + geo.tile
	teamBaseY := geo.panelStartY + geo.panelHeight - (float32(geo.frameSize)+4.5)*geo.tile

	dr.DrawOpponentFrames(state.Opponents, oppBaseY, geo)
	dr.DrawTeamFrames(state.Team, teamBaseY, geo)
	dr.DrawCombatInteractions(geo)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (dr *Drawer) DrawOpponentFrames(opponents map[string]pr.CombatPersonData, baseY float32, geo fightersPanelGeometry) {
	caseWidth := geo.panelWidth / float32(len(opponents))
	for index, key := range sortedKeys(opponents) {
		dr.DrawFighterFrame(opponents[key].Name, index, caseWidth, baseY, geo)
	}
}

func (dr *Drawer) DrawTeamFrames(team map[string]pr.CombatPersonData, baseY float32, geo fightersPanelGeometry) {
	caseWidth := geo.panelWidth / float32(len(team))
	for index, key := range sortedKeys(team) {
		dr.DrawFighterFrame(team[key].Name, index, caseWidth, baseY-0.15*geo.tile, geo)
	}
}

func (dr *Drawer) DrawFighterFrame(name string, index int, caseWidth, baseY float32, geo fightersPanelGeometry) {
	turnZoom := turnZoomFor(name, geo.curTurn)

	drawnSize := geo.baseFrameSize * turnZoom
	offsetX := float32(index)*caseWidth + caseWidth/2 - drawnSize/2
	offsetY := (drawnSize - geo.baseFrameSize) / 2

	dr.DrawRealWoodFrame(
		geo.panelStartX+offsetX,
		baseY-offsetY,
		geo.frameSize, geo.frameSize,
		0.5*turnZoom, 2, false,
	)
}

func (dr *Drawer) DrawCombatInteractions(geo fightersPanelGeometry) {
	for _, item := range dr.app.Manager.Interactions("Combat") {
		if item.Emote == nil {
			continue
		}

		cleanID := strings.TrimPrefix(item.Name, "opp_")
		cleanID = strings.TrimPrefix(cleanID, "team_")

		turnZoom := turnZoomFor(cleanID, geo.curTurn)

		zoomAlign := float32(0)
		if turnZoom != 1 {
			zoomAlign = 0.1*geo.tile - geo.zoom*0.6
		}

		shift := (geo.baseFrameSize*turnZoom - geo.baseFrameSize) / 2

		dr.DrawFighterEmote(item, shift, zoomAlign, turnZoom, geo)
		dr.DrawFighterPseudo(item, cleanID, shift, turnZoom, geo)
		dr.DrawFighterLive(item, cleanID, shift, turnZoom, geo)
	}
}

func (dr *Drawer) DrawFighterEmote(item *ui.Interaction, shift, zoomAlign, turnZoom float32, geo fightersPanelGeometry) {
	frame := item.Emote.CurrentFrame(dr.app.Variables.StartTime)

	dr.DrawImage(item.Emote.Texture,
		item.Emote.X-0.1*geo.tile-shift+zoomAlign,
		item.Emote.Y-0.1*geo.tile-shift+zoomAlign,
		frame.IndX, frame.IndY,
		frame.RatioX, frame.RatioY,
		item.Emote.Zoom/1.5*turnZoom,
		item.Emote.Rotation,
	)
}

func (dr *Drawer) DrawFighterPseudo(item *ui.Interaction, cleanID string, shift, turnZoom float32, geo fightersPanelGeometry) {
	pseudo := dr.LimitString(cleanID, 12)
	fontSize := int32(dr.app.Variables.FontSize * turnZoom)
	pseudoLen := rl.MeasureText(pseudo, fontSize)

	center := (float32(geo.frameSize) * vars.FRAME_WIDTH * geo.zoom * turnZoom * 1.5 / 4) - float32(pseudoLen)/2

	x := item.Emote.X - geo.tile - 2*shift + center
	y := item.Emote.Y + 0.55*float32(geo.baseFrameSize)

	if geo.curTurn == item.Name {
		y += 0.1 * geo.tile
	}

	rl.DrawText(pseudo, int32(x), int32(y), fontSize, dr.app.Colors["pseudo_text"])
}

func (dr *Drawer) DrawFighterLive(item *ui.Interaction, cleanID string, shift, turnZoom float32, geo fightersPanelGeometry) {
	// Define live asset indexs
	lives_index := []vars.Position{
		{X: 30, Y: 0},
		{X: 29, Y: 1},
		{X: 30, Y: 1},
		{X: 31, Y: 1},
		{X: 29, Y: 2},
		{X: 30, Y: 2},
	}

	// Get corresponding live indexs
	fighter_live := dr.GetFighterLive(item)
	live_index := vars.Position{X: 31, Y: 2}
	for index, pos := range lives_index {
		if fighter_live < float32(index)*0.2 {
			live_index = pos
			break
		}
	}
	if fighter_live == 0 {
		live_index = vars.Position{X: 30, Y: 0}
	}

	// Calculate coordinates
	zoom := dr.app.Variables.Zoom / 1.5
	if geo.curTurn == item.Name {
		zoom *= 1.25
	}
	live_drawn_width := float32(geo.tile) * zoom

	posX := item.Emote.X + 1.67*float32(geo.tile) - (live_drawn_width / 4) - vars.FRAME_WIDTH*zoom/2
	posY := item.Emote.Y - shift - 0.8*float32(geo.tile)

	if geo.curTurn == item.Name {
		posX += 0.2 * geo.tile
		posY -= 0.15 * geo.tile
	}

	// Draw live
	dr.DrawImage(
		vars.UI_SPRITE_TEXTURE,
		float32(posX), float32(posY),
		live_index.X, live_index.Y,
		1, 1, zoom, 0,
	)
}

func (dr *Drawer) GetFighterLive(item *ui.Interaction) float32 {
	// Get fighter type
	split_id := strings.SplitN(item.ID, "-", 2)
	fighter_type, fighter_name := split_id[0], split_id[1]

	fighter_list := dr.app.Variables.PanelsVariables.CombatState.Team
	if fighter_type == "npc" {
		fighter_list = dr.app.Variables.PanelsVariables.CombatState.Opponents
	}

	// Get fighter datas
	datas, ok := fighter_list[fighter_name]
	if !ok {
		return 0.5
	}

	// Update player life
	if fighter_type == "player" && fighter_name == dr.app.GetPseudo() {
		dr.app.Variables.Player.Hp = datas.Hp
	}

	return float32(datas.Hp) / float32(datas.MaxHp)
}
