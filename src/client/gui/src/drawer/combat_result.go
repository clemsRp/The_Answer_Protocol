package drawer

import (
	"strings"
	vars "tap/src/client/gui/src/variables"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawCombatResultView() {
	dr.DrawBlurredGame()

	// Draw Frame
	margin := float32(3)
	dr.DrawWoodFrameAt(
		vars.Position{X: margin, Y: margin},
		vars.Position{X: 32 - margin, Y: 18 - margin},
		1, 1, false,
	)

	tile := dr.app.Variables.Tileset_size
	font := dr.app.Variables.FontSize
	res := dr.app.Variables.PanelsVariables.CombatResult

	// Draw result
	dr.DrawResultTitle(res, margin, tile, font)
	dr.DrawResultDatas(res, tile, font)

	if time.Since(res.LastTime) >= 5*time.Second {
		dr.DrawContinueText(margin)
	}
}

func (dr *Drawer) DrawResultTitle(res *vars.CombatResultPanel, margin, tile, font float32) {
	color := rl.Green
	if res.Result == "DEFEAT" {
		color = rl.Red
	}

	center := rl.MeasureText(res.Result, int32(4*font)) / 2

	rl.DrawText(
		res.Result,
		int32(16*tile)-center,
		int32((1.5*margin+1)*tile),
		int32(4*font), color,
	)
}

func (dr *Drawer) DrawContinueText(margin float32) {
	tile := dr.app.Variables.Tileset_size
	font := dr.app.Variables.FontSize

	text := "Click to continue ..."
	align_text := float32(rl.MeasureText(
		text, int32(font),
	))

	rl.DrawText(
		text,
		int32((31-margin)*tile-align_text),
		int32((17-margin)*tile-font),
		int32(font), dr.app.Colors["pseudo_text"],
	)
}

func (dr *Drawer) DrawResultDatas(res *vars.CombatResultPanel, tile, font float32) {
	type text struct {
		color rl.Color
		text  string
	}

	middle := 16 * tile
	start_y := 9 * tile

	// Defeat case
	textList := []text{
		{
			color: rl.Red,
			text:  "That's a shame, ",
		},
		{
			color: rl.White,
			text:  "but next time will be the one !",
		},
	}
	option := ""

	// Victory case
	if res.Result != "DEFEAT" {

		if len(res.Rewards) == 0 {
			textList = []text{}

		} else {
			true_rewards := make([]string, 0)
			for _, r := range res.Rewards {
				true_rewards = append(true_rewards, dr.app.Variables.ItemConvertor[r])
			}

			rewards := strings.Join(true_rewards, ", ")

			textList = []text{
				{
					color: rl.Yellow,
					text:  "Rewards: ",
				},
				{
					color: rl.White,
					text:  rewards,
				},
			}
		}
		option = "You deserve it !"
	}

	// Draw case text
	total_width := int32(0)
	for _, t := range textList {
		total_width += rl.MeasureText(t.text, int32(font))
	}

	start_x := middle - float32(total_width)/2

	offset_x := 0
	for _, t := range textList {
		rl.DrawText(
			t.text,
			int32(start_x+float32(offset_x)),
			int32(start_y),
			int32(font), t.color,
		)

		offset_x += int(rl.MeasureText(t.text, int32(font)))
	}

	// Draw case option
	if option == "" {
		return
	}

	center_option := rl.MeasureText(option, int32(font)) / 2
	rl.DrawText(
		option,
		int32(middle)-center_option,
		int32(start_y+tile),
		int32(font), dr.app.Colors["pseudo_text"],
	)
}
