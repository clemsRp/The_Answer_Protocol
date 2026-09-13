package drawer

import (
	"tap/client/gui/src/core"
	vars "tap/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Drawer struct {
	app *core.App

	gameTexture rl.RenderTexture2D
	blurShader  rl.Shader
}

var (
	Down  = vars.Direction{X: 0, Y: 1}
	Up    = vars.Direction{X: 0, Y: -1}
	Left  = vars.Direction{X: -1, Y: 0}
	Right = vars.Direction{X: 1, Y: 0}

	DownLeft  = vars.Direction{X: -1, Y: 1}
	DownRight = vars.Direction{X: 1, Y: 1}
	UpLeft    = vars.Direction{X: -1, Y: -1}
	UpRight   = vars.Direction{X: 1, Y: -1}
)

func NewDrawer(app *core.App) *Drawer {
	dr := &Drawer{
		app:         app,
		gameTexture: rl.LoadRenderTexture(int32(app.ScreenWidth), int32(app.ScreenHeight)),
		blurShader:  rl.LoadShader("", "./client/gui/src/drawer/blur.fs"),
	}

	rl.SetShaderValue(
		dr.blurShader,
		rl.GetShaderLocation(dr.blurShader, "renderWidth"),
		[]float32{float32(app.ScreenWidth)},
		rl.ShaderUniformFloat,
	)
	rl.SetShaderValue(
		dr.blurShader,
		rl.GetShaderLocation(dr.blurShader, "renderHeight"),
		[]float32{float32(app.ScreenHeight)},
		rl.ShaderUniformFloat,
	)
	rl.SetShaderValue(
		dr.blurShader,
		rl.GetShaderLocation(dr.blurShader, "blurStrength"),
		[]float32{3.0},
		rl.ShaderUniformFloat,
	)

	return dr
}
