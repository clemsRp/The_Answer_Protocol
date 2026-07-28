package tests

// import (
// 	"tap/server"
// 	"testing"
// )

// type CommandTest struct {
// 	Name          string
// 	Command       string
// 	expectedReply string
// 	ExpectsJSON   bool
// }

// func GetMoveTestCommands(t *testing.T) []CommandTest {
// 	t.Helper()
// 	return []CommandTest{
// 		{
// 			Name:          "Invalid movement",
// 			Command:       "MOVE nowhere",
// 			expectedReply: server.ErrNoExit,
// 			ExpectsJSON:   false,
// 		},
// 		{
// 			Name:          "Valid movement",
// 			Command:       "MOVE east",
// 			expectedReply: "OK room=",
// 			ExpectsJSON:   false,
// 		},
// 		{
// 			Name:          "Valid movement",
// 			Command:       "MOVE west",
// 			expectedReply: "OK room=",
// 			ExpectsJSON:   false,
// 		},
// 		{
// 			Name:          "Valid movement",
// 			Command:       "MOVE north",
// 			expectedReply: "OK room=",
// 			ExpectsJSON:   false,
// 		},
// 		{
// 			Name:          "Valid movement",
// 			Command:       "MOVE south",
// 			expectedReply: "OK room=",
// 			ExpectsJSON:   false,
// 		}}
// }

// func GetChatTestCommands(t *testing.T) []CommandTest {
// 	return []CommandTest{
// 		{
// 			Name:          "Global chat Command valid",
// 			Command:       "CHAT GLOBAL",
// 			expectedReply: "OK",
// 			ExpectsJSON:   false,
// 		},
// 		{
// 			Name:          "Room chat Command valid",
// 			Command:       "CHAT ROOM",
// 			expectedReply: "OK",
// 			ExpectsJSON:   false,
// 		},
// 		{
// 			Name:          "Group chat Command valid",
// 			Command:       "CHAT GROUP",
// 			expectedReply: "OK",
// 			ExpectsJSON:   false,
// 		},
// 	}
// }

// func GetGroupTestCommands(t *testing.T) []CommandTest {
// 	return []CommandTest{
// 		{
// 			Name:          "Create group Command valid",
// 			Command:       "GROUP CREATE",
// 			expectedReply: "OK group=",
// 			ExpectsJSON:   false,
// 		},
// 	}
// }

// func GetTestCommands(t *testing.T) []CommandTest {
// 	t.Helper()

// 	moveCommands := GetMoveTestCommands(t)
// 	chatCommands := GetChatTestCommands(t)

// 	baseCommands := []CommandTest{
// 		{
// 			Name:          "Valid connection",
// 			Command:       "CONNECT alice",
// 			expectedReply: "OK connected",
// 			ExpectsJSON:   false,
// 		},

// 		{
// 			Name:          "Look Command valid",
// 			Command:       "LOOK",
// 			expectedReply: "OK",
// 			ExpectsJSON:   true,
// 		},
// 		{
// 			Name:          "Who Command valid",
// 			Command:       "WHO",
// 			expectedReply: "OK players=",
// 			ExpectsJSON:   false,
// 		},
// 	}

// 	allCommands := append(baseCommands, moveCommands...)
// 	allCommands = append(allCommands, chatCommands...)
// 	allCommands = append(allCommands, CommandTest{
// 		Name:          "Quitting the game",
// 		Command:       "QUIT",
// 		expectedReply: "OK bye",
// 		ExpectsJSON:   false,
// 	})
// 	return allCommands
// }

// func TestTAPProtocolCommands(t *testing.T) {
// 	log.SetOutput(io.Discard)
// 	defer log.SetOutput(os.Stderr)
// 	s := utils.SetupTestServerEngine(t)

// 	conn, err := net.DialTimeout("tcp", s.GetAddress(), 2*time.Second)

// 	if err != nil {
// 		t.Fatalf("Impossible to connect: %v", err)
// 	}
// 	defer conn.Close()

// 	bufio.NewReader(conn).ReadString('\n')

// 	command_tests := GetTestCommands(t)

// 	for _, c := range command_tests {
// 		t.Run(c.Name, func(t *testing.T) {
// 			response := utils.SendCommand(t, conn, c.Command)

// 			if c.ExpectsJSON {
// 				jsonPart := strings.TrimPrefix(response, c.expectedReply+" ")
// 				jsonPart = strings.TrimSpace(jsonPart)

// 				if !strings.HasPrefix(response, c.expectedReply) {
// 					// Utilisation du helper classique
// 					t.Error(utils.FormatMismatch(c.Command, c.expectedReply, response))
// 				} else if !json.Valid([]byte(jsonPart)) {
// 					// Utilisation du helper JSON
// 					t.Error(utils.FormatInvalidJSON(jsonPart))
// 				}
// 			} else {
// 				if !strings.HasPrefix(response, c.expectedReply) {
// 					// Utilisation du helper classique
// 					t.Error(utils.FormatMismatch(c.Command, c.expectedReply, response))
// 				}
// 			}
// 		})
// 	}
// }
