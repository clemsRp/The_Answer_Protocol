package drawer

import (
	"time"
)

type AnimKey struct {
	localID     int
	textureName string
}

type LocalIDS struct {
	Step int
	Id   int
}

func makeSequence(step, count int) []LocalIDS {
	seq := make([]LocalIDS, count)
	for i := 0; i < count; i++ {
		seq[i] = LocalIDS{Step: step, Id: i}
	}
	return seq
}

func makeRotatedSequence(step, count, offset int) []LocalIDS {
	seq := make([]LocalIDS, count)
	for i := 0; i < count; i++ {
		seq[i] = LocalIDS{Step: step, Id: (i + offset) % count}
	}
	return seq
}

var animationConvertor = map[AnimKey][]LocalIDS{
	// Chickens
	{localID: 0, textureName: "Free Chicken Sprites"}: {
		{Step: 2000, Id: 0},
		{Step: 100, Id: 1},
	},
	{localID: 1, textureName: "Free Chicken Sprites"}: {
		{Step: 1000, Id: 0},
		{Step: 100, Id: 1},
		{Step: 1000, Id: 0},
	},

	// Cows
	{localID: 0, textureName: "Free Cow Sprites"}: {{Step: 1500, Id: 0}, {Step: 100, Id: 2}, {Step: 1500, Id: 4}},
	{localID: 1, textureName: "Free Cow Sprites"}: {{Step: 1500, Id: 1}, {Step: 100, Id: 3}, {Step: 1500, Id: 5}},
	{localID: 6, textureName: "Free Cow Sprites"}: {{Step: 1500, Id: 6}, {Step: 100, Id: 8}, {Step: 1500, Id: 10}},
	{localID: 7, textureName: "Free Cow Sprites"}: {{Step: 1500, Id: 7}, {Step: 100, Id: 9}, {Step: 1500, Id: 11}},

	// Flowers
	{localID: 24, textureName: "Basic_Grass_Biom_things"}: {{Step: 750, Id: 24}, {Step: 750, Id: 25}},
	{localID: 25, textureName: "Basic_Grass_Biom_things"}: {{Step: 750, Id: 25}, {Step: 750, Id: 24}},
	{localID: 33, textureName: "Basic_Grass_Biom_things"}: {{Step: 750, Id: 33}, {Step: 750, Id: 34}},
	{localID: 34, textureName: "Basic_Grass_Biom_things"}: {{Step: 750, Id: 34}, {Step: 750, Id: 33}},

	// Water
	{localID: 0, textureName: "Water"}: makeRotatedSequence(300, 4, 0),
	{localID: 1, textureName: "Water"}: makeRotatedSequence(300, 4, 1),
	{localID: 2, textureName: "Water"}: makeRotatedSequence(300, 4, 2),
	{localID: 3, textureName: "Water"}: makeRotatedSequence(300, 4, 3),

	// Fishes
	{localID: 0, textureName: "big fish 2 swimming in cirkels"}: makeRotatedSequence(200, 15, 0), // Démarre à l'index 0
	{localID: 4, textureName: "big fish 2 swimming in cirkels"}: makeRotatedSequence(200, 15, 4), // Démarre à l'index 4
	{localID: 8, textureName: "big fish 2 swimming in cirkels"}: makeRotatedSequence(200, 15, 8), // Démarre à l'index 8
}

func (dr *Drawer) GetLocalID(localID int, textureName string) int {
	if textureName == "WaterDark" {
		textureName = "Water"

	} else if textureName == "Basic Grass Biom things 1" {
		textureName = "Basic_Grass_Biom_things"
	}

	animKey := AnimKey{localID: localID, textureName: textureName}
	frames, ok := animationConvertor[animKey]
	if !ok {
		return localID
	}

	var animDuration int64
	for _, frame := range frames {
		animDuration += int64(frame.Step)
	}

	if animDuration == 0 {
		return localID
	}

	elapsedMillis := time.Since(dr.app.Variables.StartTime).Milliseconds()
	curTime := elapsedMillis % animDuration

	var totalAnimTime int64
	for _, frame := range frames {
		totalAnimTime += int64(frame.Step)
		if totalAnimTime >= curTime {
			return frame.Id
		}
	}

	return localID
}
