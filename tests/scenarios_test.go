package tests

import (
	"bufio"
	"net"
	"strings"
	sc "tap/tests/scenarios"
	"tap/tests/utils"
	"testing"
)

func TestScenariosCommands(t *testing.T) {
	for _, family := range sc.OrderedScenarioFamilies {

		t.Run(family.FamilyName, func(t *testing.T) {

			for _, scenarioEntry := range family.Scenarios {
				s := utils.SetupTestServerEngine(t)
				connections := make(map[string]net.Conn)
				readers := make(map[string]*bufio.Reader)

				defer func() {
					for _, conn := range connections {
						conn.Close()
					}
					s.Stop()
				}()

				for _, step := range scenarioEntry.Steps {
					if strings.HasPrefix(step.Command, "connect") {
						utils.EstablishConnection(t, s, step, connections, readers)
					}

					if step.Command != "" {
						utils.SendScenarioCommand(t, connections[step.TestOnConnection], step.Command)
					}

					utils.VerifyExpectedReplies(t, step, connections, readers)
				}

			}
		})
	}
}
