package networktests

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

type CommandTest struct {
	name          string
	command       string
	expectedReply string
	expectsJSON   bool
}

func GetMoveTestCommands(t *testing.T) []CommandTest {
	t.Helper()
	return []CommandTest{
		{
			name:          "Invalid movement",
			command:       "MOVE nowhere",
			expectedReply: "ERR 301 NO_EXIT",
			expectsJSON:   false,
		},
		{
			name:          "Valid movement",
			command:       "MOVE east",
			expectedReply: "OK room=",
			expectsJSON:   false,
		},
		{
			name:          "Valid movement",
			command:       "MOVE west",
			expectedReply: "OK room=",
			expectsJSON:   false,
		},
		{
			name:          "Valid movement",
			command:       "MOVE north",
			expectedReply: "OK room=",
			expectsJSON:   false,
		},
		{
			name:          "Valid movement",
			command:       "MOVE south",
			expectedReply: "OK room=",
			expectsJSON:   false,
		}}
}

func GetTestChatCommands(t *testing.T) []CommandTest {
	return []CommandTest{
		{
			name:          "Global chat command valid",
			command:       "CHAT GLOBAL",
			expectedReply: "OK",
			expectsJSON:   false,
		},
		{
			name:          "Room chat command valid",
			command:       "CHAT ROOM",
			expectedReply: "OK",
			expectsJSON:   false,
		},
		{
			name:          "Group chat command valid",
			command:       "CHAT GROUP",
			expectedReply: "OK",
			expectsJSON:   false,
		},
	}
}

func GetTestCommands(t *testing.T) []CommandTest {
	t.Helper()
	moveCommands := GetMoveTestCommands(t)
	chatCommands := GetTestChatCommands()
	return []CommandTest{
		{
			name:          "Valid connection",
			command:       "CONNECT alice",
			expectedReply: "OK connected",
			expectsJSON:   false,
		},
		{
			name:          "Invalid connection: name already used",
			command:       "CONNECT alice",
			expectedReply: "ERR 201 NAME_IN_USE",
			expectsJSON:   false,
		},
		{
			name:          "Look command valid",
			command:       "LOOK alice",
			expectedReply: "OK",
			expectsJSON:   true,
		},
		{
			name:          "Quitting the game",
			command:       "QUIT",
			expectedReply: "OK bye",
			expectsJSON:   false,
		},
	}
}

func TestTAPProtocolCommands(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	_, addr := setupTestServer(t)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Impossible to connect: %v", err)
	}
	defer conn.Close()

	bufio.NewReader(conn).ReadString('\n')

	tests := GetCommands(t)

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			response := sendCommand(t, conn, c.command)

			if c.expectsJSON {

				jsonPart := strings.TrimPrefix(response, c.expectedReply+" ")
				jsonPart = strings.TrimSpace(jsonPart)

				if !strings.HasPrefix(response, c.expectedReply) {
					t.Errorf("Command %q -> Expected prefix: %q, Received: %q", c.command, c.expectedReply, response)
				}

				if !json.Valid([]byte(jsonPart)) {
					t.Errorf("Command %q -> Response doesn't contain valid JSON. Received payload: %s", c.command, jsonPart)
				}

			} else {
				if !strings.HasPrefix(response, c.expectedReply) {
					t.Errorf("Command %q -> Expected: %q..., Received: %q", c.command, c.expectedReply, response)
				}
			}
		})
	}
}
