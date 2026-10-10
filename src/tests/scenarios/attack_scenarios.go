package scenarios

import "tap/src/protocol"

var attackHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	aliceAttacksHostileNPC,
}

var attackNonHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Attack NON hostile NPC",
		Command: "ATTACK thomas_charbonneau",
		ExpectedReplies: []Reply{
			{protocol.ErrNpcNotHostile, "alice"},
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
			{protocol.ErrNpcNotFound, "alice"},
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
