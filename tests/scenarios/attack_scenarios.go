package scenarios

import "tap/server"

var attackHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Attack hostile NPC",
		Command: "ATTACK Granny Jeanine",
		ExpectedReplies: []Reply{
			{"OK", "alice"},
		},
		ExpectsJSON:      true,
		TestOnConnection: "alice",
	},
}

var attackNonHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Attack NON hostile NPC",
		Command: "ATTACK Nonostil",
		ExpectedReplies: []Reply{
			{server.ErrNpcNotHostile, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var attackUnexistantNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Attack unexistant NPC",
		Command: "ATTACK osdojx",
		ExpectedReplies: []Reply{
			{server.ErrNpcNotFound, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var attackScenarioFamily = ScenarioFamily{
	FamilyName: "Attack scenario family",
	Scenarios: []ScenarioEntry{
		{
			Name:  "attack hostile NPC",
			Steps: attackHostileNPCScenario,
		},
		{
			Name:  "attack non hostile NPC",
			Steps: attackNonHostileNPCScenario,
		},
		{
			Name:  "attack unexistant NPC",
			Steps: attackUnexistantNPCScenario,
		},
	},
}
