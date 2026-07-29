package scenarios

import "tap/server"

var UsernameAlreadyUsedScenario = []ScenariosCommandTest{
	ConnectAlice,
	{
		Name:    "Invalid connection: Name already used",
		Command: "CONNECT alice",
		ExpectedReplies: []Reply{
			{server.ErrNameInUse, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var NotConnectedValidCommandScenario = []ScenariosCommandTest{
	{
		Name:    "Not connected, Valid command",
		Command: "LOOK",
		ExpectedReplies: []Reply{
			{server.ErrNotConnected, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}
