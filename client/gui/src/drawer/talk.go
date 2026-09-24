package drawer

import (
	vars "tap/client/gui/src/variables"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawTalk() {
	talk := dr.app.Variables.PanelsVariables.Talk.Talking
	if talk == nil {
		return
	}

	startX := float32(vars.TALK_START_X)
	startY := float32(vars.TALK_START_Y)
	width := float32(vars.TALK_WIDTH)
	height := float32(vars.TALK_HEIGHT)

	dr.drawTalkFrame(startX, startY, width, height)

	text := dr.getAnimatedTalkText()
	dr.drawTalkContent(text, talk.Result, startX, startY, width, height)

	dr.DrawTalkButtons()
	dr.DrawTalkButtonTitle()
}

func (dr *Drawer) drawTalkFrame(startX, startY, width, height float32) {
	dr.DrawWoodFrameAt(
		vars.Position{X: startX, Y: startY},
		vars.Position{X: startX + width, Y: startY + height},
		1, 0, false,
	)
}

func (dr *Drawer) getAnimatedTalkText() string {
	talk := dr.app.Variables.PanelsVariables.Talk.Talking

	if dr.app.Variables.PanelsVariables.Talk.Finished {
		return talk.Result
	}

	if talk.Result == "" {
		return "..."
	}

	const msPerChar = 30
	elapsed := time.Since(talk.ResultStart)
	runes := []rune(talk.Result)
	index := int(elapsed.Milliseconds()) / msPerChar
	index = max(0, min(len(runes), index))

	if index >= len(talk.Result)-1 {
		dr.app.Variables.PanelsVariables.Talk.Finished = true
	}

	return string(runes[:index])
}

func (dr *Drawer) drawTalkContent(text, result string, startX, startY, width, height float32) {
	tileSize := float32(dr.app.Variables.Tileset_size)
	fontSize := dr.app.Variables.FontSize
	textColor := dr.app.Colors["pseudo_text"]

	rl.DrawText(
		text,
		int32((startX+0.75)*tileSize),
		int32((startY+0.75)*tileSize),
		int32(fontSize),
		textColor,
	)

	if result != "" {
		rl.DrawText(
			"click to continue",
			int32((startX+width-6)*tileSize),
			int32((startY+height-0.9)*tileSize),
			int32(float32(fontSize)*0.6),
			textColor,
		)
	}
}

func (dr *Drawer) DrawTalkButtonTitle() {
	// Get npc datas
	npc := dr.app.Variables.PanelsVariables.Talk.Talking.NpcID
	datas := dr.app.Variables.Npcs[npc]

	// Chose title
	var title string
	if datas.Hostile {
		title = "Attack"

	} else if !datas.RequestedQuest {
		title = "Get Quest"

	} else if !datas.CompletedQuest {
		title = "Validate Quest"
	}

	// Calculate coordinates
	tile := dr.app.Variables.Tileset_size
	font := dr.app.Variables.FontSize

	center := float32(vars.TALK_START_X+vars.TALK_WIDTH/2) * tile
	align_title := rl.MeasureText(
		title, int32(font),
	) / 2

	posX := center - float32(align_title)
	posY := float32(vars.TALK_START_Y)*tile - 1.8*font

	// Draw title
	rl.DrawText(
		title,
		int32(posX), int32(posY),
		int32(font), dr.app.Colors["panel_text"],
	)
}
