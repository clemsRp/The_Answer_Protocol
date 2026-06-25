package main

import (
	"fmt"
	"time"
)

func get_timestamp() map[string]int {
	// Get timestamp value
	time := time.Now().Unix()
	timestamp := time - t_start

	// Convert it and stock it inside a map
	res := make(map[string]int)

	res["days"] = int(timestamp / 86400)
	res["hours"] = int((timestamp / 3600) % 24)
	res["min"] = int((timestamp / 60) % 60)
	res["sec"] = int(timestamp % 60)

	return res
}

func print_timestamp(timestamp map[string]int) {
	for unit, duration := range timestamp {
		if duration > 0 {
			fmt.Printf(" %d %s", duration, unit)
		}
	}
	fmt.Printf("\n")
}
