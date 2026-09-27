package testTUI

import (
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"tap/client/tui"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestTUINormalDisconnection(t *testing.T) {
	s, e := utils.SetupTestServerEngine(t, "../../world.json")

	conn, err := net.Dial("tcp", s.GetAddress())
	if err != nil {
		t.Fatalf("Connection error: %s", err)
	}

	baseline := runtime.NumGoroutine()

	cli := tui.NewTuiClient(conn)

	go cli.Start()
	time.Sleep(200 * time.Millisecond)
	cli.Stop()
	s.Stop()
	e.Stop()

	success := false
	for i := 0; i < 2000; i++ {
		if runtime.NumGoroutine() <= baseline {
			success = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !success {
		current := runtime.NumGoroutine()
		t.Errorf("Go goroutines leaks detected: started with %d, ended with %d.", baseline, current)
		pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	} else {
		t.Log("Success: tui shutdown correctly without any leak.")
	}
}
