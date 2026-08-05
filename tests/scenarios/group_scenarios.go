package scenarios

import (
	"tap/protocol"
)

var leaveUnexistantGroupScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Leave unexistant group",
		Command: "GROUP LEAVE",
		ExpectedReplies: []Reply{
			{protocol.ErrNotInGroup, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var joinGroupAgainScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	{
		Name:    "Bob accepts Alice's invitation",
		Command: "GROUP JOIN alice",
		ExpectedReplies: []Reply{
			{protocol.ErrAlreadyInGroup, "bob"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "bob",
	},
}

var leaveTwiceExistantGroupScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	aliceLeavesGroup,
	{
		Name:    "Alice tries to leave group again",
		Command: "GROUP LEAVE",
		ExpectedReplies: []Reply{
			{protocol.ErrNotInGroup, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var inviteUnexistantPersonScenario = []ScenariosCommandTest{
	connectAlice,
	aliceCreatesGroup,
	{
		Name:    "Alice tries to invite Unknown in group",
		Command: "GROUP INVITE Unknown",
		ExpectedReplies: []Reply{
			{protocol.ErrUnknownUser, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var groupScenarioFamily = ScenarioFamily{
	FamilyName: "Group scenarios family",
	Scenarios: []ScenarioEntry{
		{
			Name:  "Leave unexistant group",
			Steps: leaveUnexistantGroupScenario,
		},
		{
			Name:  "Join group again (already in group)",
			Steps: joinGroupAgainScenario,
		},
		{
			Name:  "Leave twice existant group",
			Steps: leaveTwiceExistantGroupScenario,
		},
		{
			Name:  "Invite unexistant person",
			Steps: inviteUnexistantPersonScenario,
		},
	},
}
