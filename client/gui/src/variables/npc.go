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
	NpcConvertor = map[string]Npc{
		"guide": {
			Texture: NPC_TEXTURE,
			RatioX:  1,
			RatioY:  1.5,
			IndXs:   []float32{4, 5, 6, 7, 0, 1, 2, 3},
			IndYs:   []float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			Pos:     Position{X: 17, Y: 7},
		},
		"maelis": {
			Texture: "bat_animations",
			RatioX:  2,
			RatioY:  2,
			IndXs:   []float32{4, 6, 8, 0, 2},
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
	}
)
