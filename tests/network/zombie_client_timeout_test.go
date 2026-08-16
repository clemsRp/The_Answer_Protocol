package network

import (
	"runtime"
	"tap/engine"
	pr "tap/protocol"
	"tap/server"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestZombieClientTimeout(t *testing.T) {
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	exchanger := pr.Exchanger{ServerInput: make(chan pr.ServerRequest, 100),
		ServerOutput: make(chan pr.EngineResponse, 100),
		JoinChan:     make(chan string, 10),
		LeaveChan:    make(chan string, 10)}

	world, err := engine.Get_map("../../world.json")
	if err != nil {
		t.Fatalf("ERROR parsing: %v", err)
	}

	s, err := server.NewServer("localhost:0", exchanger)
	if err != nil {
		t.Fatalf("Server couldn't be created %v", err)
	}

	e := engine.NewEngine(world, exchanger)

	s.IdleTimeout = 1 * time.Second

	go e.Start()
	go s.Start()

	t.Cleanup(func() {
		s.Stop()
		e.Stop()
	})

	conn, _ := utils.ConnectAndGreet(t, s.GetAddress())

	time.Sleep(2 * s.IdleTimeout)

	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err == nil {
		t.Fatalf("Fail: zombie connection still active. server didn't close it with timeout.")
	}

	s.Stop()
	e.Stop()
	time.Sleep(100 * time.Millisecond)

	current := runtime.NumGoroutine()
	if current > baseline {
		t.Errorf("Fail:  Zombie connection, input.scan() is blocked and goroutines are not properly stopped.")
	} else {
		t.Log("Success: server detected inactivity and cut the connection properly.")
	}
}
