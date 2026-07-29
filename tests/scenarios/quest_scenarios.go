package scenarios

import "tap/server"

var questRequestNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "quest hostile NPC",
		Command: "quest Granny Jeanine",
		ExpectedReplies: []Reply{
			{"OK", "alice"},
		},
		ExpectsJSON:      true,
		TestOnConnection: "alice",
	},
}

var questNonHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "quest NON hostile NPC",
		Command: "quest Nonostil",
		ExpectedReplies: []Reply{
			{server.ErrNpcNotHostile, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var questUnexistantNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "quest unexistant NPC",
		Command: "quest osdojx",
		ExpectedReplies: []Reply{
			{server.ErrNpcNotFound, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var questScenarioFamily = ScenarioFamily{
	FamilyName: "quest scenario family",
	Scenarios: []ScenarioEntry{
		{
			Name:  "quest hostile NPC",
			Steps: questHostileNPCScenario,
		},
		{
			Name:  "quest non hostile NPC",
			Steps: questNonHostileNPCScenario,
		},
		{
			Name:  "quest unexistant NPC",
			Steps: questUnexistantNPCScenario,
		},
	},
}
