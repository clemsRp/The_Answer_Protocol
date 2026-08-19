package scenarios

import (
	"bufio"
	"net"
	"strings"
	"tap/tests/utils"
	"testing"
)

func RunScenario(t *testing.T, scenarioEntry ScenarioEntry) {
	s, _ := utils.SetupTestServerEngine(t, "../../world.json")
	connections := make(map[string]net.Conn)
	readers := make(map[string]*bufio.Reader)

	for _, step := range scenarioEntry.Steps {
		if strings.HasPrefix(strings.ToLower(step.Command), "connect") {
			EstablishConnectionForScenario(t, s, step, connections, readers)
		}
		conn := connections[step.TestOnConnection]
		if step.Command != "" {
			if conn == nil {
				t.Fatalf("CRASH AVOIDED : scenario '%s' tried to send '%s' command pour '%s' user, You maybe forgot TestOnConnection in the Action or Connect In the scenario.",
					scenarioEntry.Name, step.Command, step.TestOnConnection)
			}
			SendScenarioCommand(t, conn, step.Command)
		}

		VerifyExpectedReplies(t, step, connections, readers)
	}
}
