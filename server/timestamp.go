package server

import (
	"time"
)

func (s *Server) get_timestamp() map[string]int {
	// Get timestamp value
	time := time.Now().Unix()
	timestamp := time - s.t_start

	// Convert it and stock it inside a map
	res := make(map[string]int)

	res["days"] = int(timestamp / 86400)
	res["hours"] = int((timestamp / 3600) % 24)
	res["min"] = int((timestamp / 60) % 60)
	res["sec"] = int(timestamp % 60)

	return res
}
