package variables

import rl "github.com/gen2brain/raylib-go/raylib"

func GetColors() map[string]rl.Color {
	return map[string]rl.Color{
		"panel_text":  rl.NewColor(182, 137, 98, 255),
		"pseudo_text": rl.NewColor(232, 207, 166, 255),
		"group_text":  rl.NewColor(107, 75, 91, 255),
	}
}
