package scenarios

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"tap/tests/utils"
	"testing"
	"time"
)

func RunConcurrentScenario(t *testing.T, scenario ConcurrentScenario) {
	s, _ := utils.SetupTestServerEngine(t, "../world.json")
	connections := make(map[string]net.Conn)
	readers := make(map[string]*bufio.Reader)

	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
		s.Stop()
	}()

	for _, step := range scenario.SetupSteps {
		if strings.HasPrefix(strings.ToLower(step.Command), "connect") {
			EstablishConnectionForScenario(t, s, step, connections, readers)
		}

		if step.Command != "" {
			SendScenarioCommand(t, connections[step.TestOnConnection], step.Command)
		}

		VerifyExpectedReplies(t, step, connections, readers)
	}

	var wg sync.WaitGroup
	startPistol := make(chan struct{})

	var mu sync.Mutex
	results := make(map[string]string)

	for user, cmd := range scenario.ConcurrentCmds {
		wg.Add(1)

		go func(u string, c string) {
			defer wg.Done()

			conn := connections[u]
			reader := readers[u]

			<-startPistol

			fmt.Fprintf(conn, "%s\n", c)

			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			res, _ := reader.ReadString('\n')

			mu.Lock()
			results[u] = res
			mu.Unlock()
		}(user, cmd)
	}

	close(startPistol)
	wg.Wait()

	if scenario.ValidationFunc != nil {
		scenario.ValidationFunc(t, results)
	}
}
