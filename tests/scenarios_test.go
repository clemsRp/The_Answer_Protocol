package tests

import (
	"bufio"
	"net"
	"strings"
	sc "tap/tests/scenarios"
	"tap/tests/utils"
	"testing"
	"time"
)

var scenariosGroups = []struct {
	GroupName string
	Function  func(t *testing.T) []sc.ScenariosCommandTest
}{
	{"Join Group Again", sc.GetJoinGroupAgainScenario},
	{"Leave Twice Existant Group", sc.GetLeaveTwiceExistantGroupScenario},
	{"Leave Unexistant Group", sc.GetLeaveUnexistantGroupScenario},
	{"Invite Unexistant Person", sc.GetInviteUnexistantPersonScenario},
}

func TestScenariosCommands(t *testing.T) {
	var res string
	var err error

	for _, group := range scenariosGroups {
		s := utils.SetupTestServerEngine(t)

		connections := make(map[string]net.Conn)
		readers := make(map[string]*bufio.Reader)
		scenarios := group.Function(t)
		t.Run(group.GroupName, func(t *testing.T) {
			for _, scenario := range scenarios {

				// Handle Start of connection
				if strings.HasPrefix(scenario.Command, "CONNECT") {
					conn, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)
					if err != nil {
						t.Fatalf(scenario.TestOnConnection+" failed to connect: %v", err)
					}
					connections[scenario.TestOnConnection] = conn

					readers[scenario.TestOnConnection] = bufio.NewReader(conn)
					conn.SetReadDeadline(time.Now().Add(1 * time.Second))
					res, err = readers[scenario.TestOnConnection].ReadString('\n')
					if err != nil {
						t.Fatalf("Read error: %v", err)
					}
					if !strings.HasPrefix(res, "OK hello proto=") {
						t.Error(utils.FormatMismatch(scenario.Command, "OK hello proto=", res))
					}
				}

				// Send Command if necessary and get server response
				if scenario.Command != "" {
					utils.SendScenarioCommand(t, connections[scenario.TestOnConnection], scenario.Command)
				}

				// Compare responses
				for _, reply := range scenario.ExpectedReplies {
					// Get response
					connections[reply.User].SetReadDeadline(time.Now().Add(1 * time.Second))
					res, err = readers[reply.User].ReadString('\n')
					if err != nil {
						t.Fatalf("Read error: %v", err)
					}

					if !strings.HasPrefix(res, reply.Msg) {
						t.Error(utils.FormatMismatch(scenario.Command, reply.Msg, res))
					}
				}
			}
		})
		// Close connections
		for _, conn := range connections {
			conn.Close()
		}
	}
}
