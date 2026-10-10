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
	{localID: 8, textureName: "Chicken_Baby_Red"}: {
		{Step: 200, Id: 8},
		{Step: 200, Id: 9},
		{Step: 200, Id: 10},
		{Step: 1500, Id: 11},
		{Step: 750, Id: 12},
		{Step: 1500, Id: 13},
		{Step: 750, Id: 14},
	},
	{localID: 11, textureName: "Chicken_Baby_Red"}: {
		{Step: 1500, Id: 11},
		{Step: 750, Id: 12},
		{Step: 1500, Id: 13},
		{Step: 750, Id: 14},
		{Step: 200, Id: 8},
		{Step: 200, Id: 9},
		{Step: 200, Id: 10},
	},

	// Cows
	{localID: 0, textureName: "Free Cow Sprites"}:                    {{Step: 1500, Id: 0}, {Step: 100, Id: 2}, {Step: 1500, Id: 4}},
	{localID: 1, textureName: "Free Cow Sprites"}:                    {{Step: 1500, Id: 1}, {Step: 100, Id: 3}, {Step: 1500, Id: 5}},
	{localID: 6, textureName: "Free Cow Sprites"}:                    {{Step: 1500, Id: 6}, {Step: 100, Id: 8}, {Step: 1500, Id: 10}},
	{localID: 7, textureName: "Free Cow Sprites"}:                    {{Step: 1500, Id: 7}, {Step: 100, Id: 9}, {Step: 1500, Id: 11}},
	{localID: 176, textureName: "baby light cow animations sprites"}: {{Step: 2000, Id: 176}, {Step: 2000, Id: 178}},
	{localID: 177, textureName: "baby light cow animations sprites"}: {{Step: 2000, Id: 177}, {Step: 2000, Id: 179}},

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
	{localID: 0, textureName: "big fish 2 swimming in cirkels"}: makeRotatedSequence(200, 15, 0),
	{localID: 4, textureName: "big fish 2 swimming in cirkels"}: makeRotatedSequence(200, 15, 4),
	{localID: 8, textureName: "big fish 2 swimming in cirkels"}: makeRotatedSequence(200, 15, 8),

	// Boats
	{localID: 0, textureName: "Boats"}:  {{Step: 500, Id: 0}, {Step: 500, Id: 3}},
	{localID: 1, textureName: "Boats"}:  {{Step: 500, Id: 1}, {Step: 500, Id: 4}},
	{localID: 2, textureName: "Boats"}:  {{Step: 500, Id: 2}, {Step: 500, Id: 5}},
	{localID: 9, textureName: "Boats"}:  {{Step: 500, Id: 9}, {Step: 500, Id: 12}},
	{localID: 10, textureName: "Boats"}: {{Step: 500, Id: 10}, {Step: 500, Id: 13}},
	{localID: 11, textureName: "Boats"}: {{Step: 500, Id: 11}, {Step: 500, Id: 14}},

	{localID: 18, textureName: "Boats"}: {{Step: 500, Id: 18}, {Step: 500, Id: 21}},
	{localID: 19, textureName: "Boats"}: {{Step: 500, Id: 19}, {Step: 500, Id: 22}},
	{localID: 20, textureName: "Boats"}: {{Step: 500, Id: 20}, {Step: 500, Id: 23}},
	{localID: 27, textureName: "Boats"}: {{Step: 500, Id: 27}, {Step: 500, Id: 30}},
	{localID: 28, textureName: "Boats"}: {{Step: 500, Id: 28}, {Step: 500, Id: 31}},
	{localID: 29, textureName: "Boats"}: {{Step: 500, Id: 29}, {Step: 500, Id: 32}},

	// Mails Boxs
	{localID: 11, textureName: "Mailbox Animation Frames"}: {{Step: 500, Id: 11}, {Step: 500, Id: 44}},
	{localID: 12, textureName: "Mailbox Animation Frames"}: {{Step: 500, Id: 12}, {Step: 500, Id: 45}},
	{localID: 22, textureName: "Mailbox Animation Frames"}: {{Step: 500, Id: 22}, {Step: 500, Id: 55}},
	{localID: 23, textureName: "Mailbox Animation Frames"}: {{Step: 500, Id: 23}, {Step: 500, Id: 56}},
}

func (dr *Drawer) GetLocalID(localID int, textureName string) int {
	if textureName == "WaterDark" || textureName == "WaterLight" {
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
