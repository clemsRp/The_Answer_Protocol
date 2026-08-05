package network

import (
	"net"
	"runtime"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestServerGracefulShutdown_NoGoroutineLeaks(t *testing.T) {
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	s, e := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}

	conn.Close()

	s.Stop()
	e.Stop()

	time.Sleep(100 * time.Millisecond)

	current := runtime.NumGoroutine()

	if current > baseline {
		t.Errorf("Go goroutines leaks detected : we started with %d go routines (test environment goroutines), we end with %d. Server didn't end correctly.", baseline, current)
	} else {
		t.Log("Success : server shutdown correctly without any leak.")
	}
}
