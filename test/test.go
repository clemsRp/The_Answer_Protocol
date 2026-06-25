package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. Timestamp en secondes (Unix Epoch)
	sec := time.Now().Unix()
	fmt.Printf("Secondes : %d\n", sec)

	// 2. Timestamp en millisecondes
	milli := time.Now().UnixMilli()
	fmt.Printf("Millisecondes : %d\n", milli)

	// 3. Timestamp en microsecondes
	micro := time.Now().UnixMicro()
	fmt.Printf("Microsecondes : %d\n", micro)

	// 4. Timestamp en nanosecondes
	nano := time.Now().UnixNano()
	fmt.Printf("Nanosecondes : %d\n", nano)
}