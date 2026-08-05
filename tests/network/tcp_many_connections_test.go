package network

import (
	"net"
	"sync"
	"tap/tests/utils"
	"testing"
	"time"
)

func TestConcurrentClientConnections(t *testing.T) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	addr := s.GetAddress()

	const numClients = 100
	var wg sync.WaitGroup

	// Channel to safely collect errors from goroutines
	errCh := make(chan error, numClients)

	// Launch all clients at the same time
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
			if err != nil {
				errCh <- err
				return
			}
			defer conn.Close()

			// Small delay to simulate an active user maintaining the connection
			time.Sleep(100 * time.Millisecond)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errCh)

	// Check if any goroutine reported an error
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Fatalf("Fail: %d out of %d clients failed to connect simultaneously. First error: %v", len(errors), numClients, errors[0])
	}

	t.Logf("Success: Server successfully handled %d concurrent connections without crashing.", numClients)
}
