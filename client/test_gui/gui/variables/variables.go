package variables

type Variables struct {
	X            float32
	Y            float32
	Tileset_size int
	Current_room string
}

func GetVariables() *Variables {
	return &Variables{
		X:            float32(100),
		Y:            float32(100),
		Current_room: "true",
	}
}
