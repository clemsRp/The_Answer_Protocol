package scenarios

import "tap/protocol"

var usernameAlreadyUsedScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Invalid connection: Name already used",
		Command: "CONNECT alice",
		ExpectedReplies: []Reply{
			{protocol.ErrNameInUse, "alice"},
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
