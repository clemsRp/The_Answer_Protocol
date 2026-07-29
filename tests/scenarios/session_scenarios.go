package scenarios

import "tap/server"

var usernameAlreadyUsedScenario = []ScenariosCommandTest{
	connectAlice,
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

var sessionScenarioFamily = ScenarioFamily{
	FamilyName: "Session family",
	Scenarios: []ScenarioEntry{
		{
			Name:  "Try to connect with already used name",
			Steps: usernameAlreadyUsedScenario,
		},
	},
}
