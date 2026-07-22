package network

import (
	"fmt"
	"strings"
)

func (s *Server) GameLoop() {

	for msg := range s.messageChan {
		payload := strings.TrimSpace(string(msg.payload))
		// parts := strings.Split(payload, " ")
		fmt.Printf("received message from connection (%s): %s\n", msg.from, payload)
		// if strings.HasPrefix(payload, "CONNECT") && len(parts) == 2 {

		// }
	}
}
