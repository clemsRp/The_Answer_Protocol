package networktests

import (
	"bufio"
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
}

func GetCommands() []CommandTest {

	return []CommandTest{
		{
			name:          "Connexion valide",
			command:       "CONNECT alice",
			expectedReply: "OK connected",
		},
		{
			name:          "Connexion avec un nom déjà pris",
			command:       "CONNECT alice",
			expectedReply: "ERR 201 NAME_IN_USE",
		},
		{
			name:          "Mouvement invalide",
			command:       "MOVE nowhere",
			expectedReply: "ERR 301 NO_EXIT",
		},
	}

}

func TestTAPProtocolCommands(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	_, addr := setupTestServer(t)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Impossible de se connecter: %v", err)
	}
	defer conn.Close()

	bufio.NewReader(conn).ReadString('\n')

	tests := GetCommands()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			response := sendCommand(t, conn, tc.command)

			if !strings.HasPrefix(response, tc.expectedReply) {
				t.Errorf("Commande %q -> Attendu: %q..., Reçu: %q", tc.command, tc.expectedReply, response)
			}
		})
	}
}
