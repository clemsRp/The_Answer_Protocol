package variables

type NpcDatas struct {
	Hostile        bool
	HasQuest       bool
	RequestedQuest bool
	CompletedQuest bool
	QuestID        string
}

type Npc struct {
	Texture        string
	RatioX, RatioY float32
	IndXs, IndYs   []float32
	Pos            Position
}

var (
	// TODO Define all npcs
	NpcConvertor = map[string]Npc{
		"thomas_charbonneau": {
			Texture: "free_character_spritesheet_by-cupnooble",
			RatioX:  1,
			RatioY:  2,
			IndXs:   []float32{0, 1, 2, 3, 4, 5, 6, 7},
			IndYs:   []float32{0, 0, 0, 0, 0, 0, 0, 0},
			Pos:     Position{X: 12, Y: 6},
		},
		"maelis": {
			Texture: "bat_animations",
			RatioX:  2,
			RatioY:  2,
			IndXs:   []float32{0, 2, 4, 6, 8},
			IndYs:   []float32{0, 0, 0, 0, 0},
			Pos:     Position{X: 13, Y: 13},
		},
		"le_prince": {
			Texture: "frog_spritesheet",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0, 1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			IndYs:   []float32{4, 4, 4, 4, 4, 4, 4, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
			Pos:     Position{X: 20, Y: 9},
		},
		"maya": {
			Texture: "bee_spritesheet",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0, 1, 2, 3, 4, 5, 6, 7},
			IndYs:   []float32{0, 0, 0, 0, 0, 0, 0, 0},
			Pos:     Position{X: 18, Y: 6},
		},
		"bob": {
			Texture: "BouncingSpritesheet",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0, 1, 2, 0, 1, 2, 0},
			IndYs:   []float32{0, 0, 0, 1, 1, 1, 2},
			Pos:     Position{X: 20, Y: 6},
		},
		"gabinap": {
			Texture: "gabinap",
			RatioX:  8,
			RatioY:  13,
			IndXs:   []float32{10, 20},
			IndYs:   []float32{7, 7},
			Pos:     Position{X: 21, Y: 6},
		},
		// TODO
		// "bernard": {
		// 	Texture: "frog_spritesheet",
		// 	RatioX:  1,
		// 	RatioY:  1,
		// 	IndXs:   []float32{0, 1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		// 	IndYs:   []float32{4, 4, 4, 4, 4, 4, 4, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		// 	Pos:     Position{X: 20, Y: 9},
		// },
		"salomon": {
			Texture: "salomon",
			RatioX:  8,
			RatioY:  13,
			IndXs:   []float32{10, 20},
			IndYs:   []float32{7, 7},
			Pos:     Position{X: 28, Y: 12},
		},
		// "franck": {
		// 	Texture: "BouncingSpritesheet",
		// 	RatioX:  1,
		// 	RatioY:  1,
		// 	IndXs:   []float32{0, 1, 2, 0, 1, 2, 0},
		// 	IndYs:   []float32{0, 0, 0, 1, 1, 1, 2},
		// 	Pos:     Position{X: 20, Y: 6},
		// },
		// "roger": {
		// 	Texture: "gabinap",
		// 	RatioX:  2.3,
		// 	RatioY:  5,
		// 	IndXs:   []float32{1.6, 5.5, 9.7},
		// 	IndYs:   []float32{0, 0, 0},
		// 	Pos:     Position{X: 21, Y: 6},
		// },
	}
)
