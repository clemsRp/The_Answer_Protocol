package gui

import (
	"time"
)

type AnimKey struct {
	localID     int
	textureName string
}

type localIDS struct {
	Step int
	Id   int
}

var (
	water1 = localIDS{Step: 300, Id: 0}
	water2 = localIDS{Step: 300, Id: 1}
	water3 = localIDS{Step: 300, Id: 2}
	water4 = localIDS{Step: 300, Id: 3}

	animation_convertor = map[AnimKey][]localIDS{
		// Chickens
		{0, "Free Chicken Sprites"}: []localIDS{
			{Step: 2000, Id: 0},
			{Step: 100, Id: 1},
		},
		{1, "Free Chicken Sprites"}: []localIDS{
			{Step: 1000, Id: 0},
			{Step: 100, Id: 1},
			{Step: 1000, Id: 0},
		},

		// Flowers
		{24, "Basic Grass Biom things 1"}: []localIDS{
			{Step: 750, Id: 24},
			{Step: 750, Id: 25},
		},
		{25, "Basic Grass Biom things 1"}: []localIDS{
			{Step: 750, Id: 25},
			{Step: 750, Id: 24},
		},
		{33, "Basic Grass Biom things 1"}: []localIDS{
			{Step: 750, Id: 33},
			{Step: 750, Id: 34},
		},
		{34, "Basic Grass Biom things 1"}: []localIDS{
			{Step: 750, Id: 34},
			{Step: 750, Id: 33},
		},

		// Water
		{0, "Water"}: []localIDS{
			water1,
			water2,
			water3,
			water4,
		},
		{1, "Water"}: []localIDS{
			water2,
			water3,
			water4,
			water1,
		},
		{2, "Water"}: []localIDS{
			water3,
			water4,
			water1,
			water2,
		},
		{3, "Water"}: []localIDS{
			water4,
			water1,
			water2,
			water3,
		},
	}
)

func (app *App) GetLocalID(localID int, textureName string) int {
	// Get localIDS
	animKey := AnimKey{localID: localID, textureName: textureName}
	local_IDS, ok := animation_convertor[animKey]
	if !ok {
		return localID
	}

	// Get animation duration
	animDuration := 0
	for _, new_local_id := range local_IDS {
		animDuration += new_local_id.Step
	}

	// Get cur_time
	elapsedMillis := time.Since(app.variables.StartTime).Milliseconds()
	cur_time := elapsedMillis % int64(animDuration)

	// Return new localID
	total_anim_time := 0
	for _, new_local_id := range local_IDS {
		total_anim_time += int(new_local_id.Step)
		if int64(total_anim_time) >= cur_time {
			return new_local_id.Id
		}
	}

	return localID
}
