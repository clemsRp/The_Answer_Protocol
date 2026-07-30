package scenarios

import "tap/protocol"

var talkToNpcScenario = []ScenariosCommandTest{
	connectAlice,
	aliceTalksToNPC,
}

var talkToUnexistantNPCScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Alice talks to unexistant NPC",
		Command: "TALK general",
		ExpectedReplies: []Reply{
			{protocol.ErrNpcNotFound, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var talkScenarioFamily = ScenarioFamily{
	FamilyName: "Talk scenarios family",
	Scenarios: []ScenarioEntry{
		{
			Name:  "Talk to npc successfully",
			Steps: talkToNpcScenario,
		},
		{
			Name:  "Try to talk to unexistant NPC",
			Steps: talkToUnexistantNPCScenario,
		},
	},
}
