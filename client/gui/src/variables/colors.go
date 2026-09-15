package variables

import rl "github.com/gen2brain/raylib-go/raylib"

func GetColors() map[string]rl.Color {
	return map[string]rl.Color{
		"panel_text": rl.NewColor(182, 137, 98, 255),
	}
}
