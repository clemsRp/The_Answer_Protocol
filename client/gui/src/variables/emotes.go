package variables

import "time"

type EmotesPanel struct {
	Open           bool
	OpenTime       time.Time
	LastEmoteTime  time.Time
	LastEmoteIndex int
}
