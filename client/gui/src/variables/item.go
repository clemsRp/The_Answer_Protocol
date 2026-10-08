package variables

type ItemDatas struct {
	Hostile        bool
	HasQuest       bool
	RequestedQuest bool
	CompletedQuest bool
	QuestID        string
}

type Item struct {
	Texture        string
	RatioX, RatioY float32
	IndXs, IndYs   []float32
	Pos            Position
}

var (
	ItemConvertor = map[string]Item{
		"mais_sucre": {
			Texture: "All items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{1},
			IndYs:   []float32{1},
			Pos:     Position{X: 9, Y: 7},
		},
		"flocons_d_avoine": {
			Texture: "All items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0},
			IndYs:   []float32{0},
			Pos:     Position{X: 9, Y: 12},
		},
		"acier": {
			Texture: "dungeon_items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{3},
			IndYs:   []float32{0},
			Pos:     Position{X: 8, Y: 3},
		},
		"branche_solide": {
			Texture: "All items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{2},
			IndYs:   []float32{2},
			Pos:     Position{X: 8, Y: 11},
		},
		"canne_a_peche": {
			Texture: "fishing_rod",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0},
			IndYs:   []float32{0},
			Pos:     Position{X: 8, Y: 11},
		},
		"cancoillotte": {
			Texture: "All items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{7},
			IndYs:   []float32{10},
			Pos:     Position{X: -10, Y: -10},
		},
		"graines": {
			Texture: "fishing_rod",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0},
			IndYs:   []float32{8},
			Pos:     Position{X: -10, Y: -10},
		},
		"pioche": {
			Texture: "dungeon_items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{0},
			IndYs:   []float32{0},
			Pos:     Position{X: 27, Y: 14},
		},
		"marshmallow": {
			Texture: "All items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{3},
			IndYs:   []float32{11},
			Pos:     Position{X: 22, Y: 13},
		},
		"fricadelle": {
			Texture: "dungeon_items",
			RatioX:  1,
			RatioY:  1,
			IndXs:   []float32{1},
			IndYs:   []float32{13},
			Pos:     Position{X: -10, Y: -10},
		},
	}
)
