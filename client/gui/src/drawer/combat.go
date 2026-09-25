package drawer

import (
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (dr *Drawer) DrawCombatView() {
	dr.DrawBlurredGame()

	dr.DrawCombatChatZone()
	dr.DrawCombatFightersZone()
	dr.DrawCombatActionsZone()

	dr.DrawPlayerPanel()
	dr.DrawGameEmotes()
	dr.DrawInventoryEmotes()
	dr.DrawCombatInventoryButtons()
}

func (dr *Drawer) DrawBlurredGame() {
	rl.BeginTextureMode(dr.gameTexture)
	dr.DrawGameView()
	rl.EndTextureMode()

	rl.BeginShaderMode(dr.blurShader)
	rl.DrawTextureRec(
		dr.gameTexture.Texture,
		rl.NewRectangle(
			0, 0,
			float32(dr.gameTexture.Texture.Width),
			-float32(dr.gameTexture.Texture.Height),
		),
		rl.NewVector2(0, 0),
		rl.White,
	)
	rl.EndShaderMode()
}

func (dr *Drawer) DrawCombatFrame(startX, startY, endX, endY float32) {
	dr.DrawWoodFrameAt(
		vars.Position{X: startX, Y: startY},
		vars.Position{X: endX, Y: endY},
		1, 0, false,
	)
}
