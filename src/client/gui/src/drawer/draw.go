package drawer

import (
	"embed"
	"fmt"
	"tap/src/client/gui/src/core"
	vars "tap/src/client/gui/src/variables"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed shader/*.fs
var shadersFS embed.FS

type Drawer struct {
	app *core.App

	gameTexture         rl.RenderTexture2D
	combatResultTexture rl.RenderTexture2D
	blurShader          rl.Shader
	darkenShader        rl.Shader

	chatTexture  rl.RenderTexture2D
	groupTexture rl.RenderTexture2D
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
	blurBytes, errBlur := shadersFS.ReadFile("shader/blur.fs")
	if errBlur != nil {
		fmt.Println("Warning: Failed to load blur shader from embed:", errBlur)
	}
	
	darkenBytes, errDarken := shadersFS.ReadFile("shader/darken.fs")
	if errDarken != nil {
		fmt.Println("Warning: Failed to load darken shader from embed:", errDarken)
	}

	dr := &Drawer{
		app:                 app,
		gameTexture:         rl.LoadRenderTexture(int32(app.ScreenWidth), int32(app.ScreenHeight)),
		combatResultTexture: rl.LoadRenderTexture(int32(app.ScreenWidth), int32(app.ScreenHeight)),
		blurShader:          rl.LoadShaderFromMemory("", string(blurBytes)),
		darkenShader:        rl.LoadShaderFromMemory("", string(darkenBytes)),
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
		[]float32{5.0},
		rl.ShaderUniformFloat,
	)

	return dr
}
