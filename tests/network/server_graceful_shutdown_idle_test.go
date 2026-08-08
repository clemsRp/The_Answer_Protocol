package network

import (
	"runtime"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestServerGracefulShutdown_Idle(t *testing.T) {
	time.Sleep(50 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	s, e := utils.SetupTestServerEngine(t, "../../world.json")

	defer e.Stop()
	defer s.Stop()

	s.Stop()
	e.Stop()

	success := false
	for i := 0; i < 100; i++ {
		if runtime.NumGoroutine() <= baseline {
			success = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !success {
		current := runtime.NumGoroutine()
		t.Errorf("Go goroutines leaks detected: started with %d, ended with %d.", baseline, current)
	} else {
		t.Log("Success: server shutdown correctly without any leak.")
	}
}
