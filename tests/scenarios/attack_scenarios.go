package scenarios

import "tap/protocol"

var attackHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	aliceAttacksHostileNPC,
	aliceAttacksHostileNPCAgain,
	aliceDefeatsHostileNPC,
}

var attackNonHostileNPCScenario = []ScenariosCommandTest{
	connectAlice,
	aliceMovesToHealthAisle,
	{
		Name:    "Attack NON hostile NPC",
		Command: "ATTACK granny_jeanine",
		ExpectedReplies: []Reply{
			{protocol.ErrNpcNotHostile, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var attackGroupAndFleeScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	aliceStartsGroupCombat,
	aliceFleesCombat,
	bobAttacksAfterAliceFlees,
	bobFleesLastFromCombat,
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
		{
			Name:  "group combat survives an ally fleeing",
			Steps: attackGroupAndFleeScenario,
		},
	},
}
