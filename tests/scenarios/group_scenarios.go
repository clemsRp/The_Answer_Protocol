package scenarios

import (
	"tap/server"
)

var LeaveUnexistantGroupScenario = []ScenariosCommandTest{
	ConnectAlice,
	{
		Name:    "Leave unexistant group",
		Command: "GROUP LEAVE",
		ExpectedReplies: []Reply{
			{server.ErrNotInGroup, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var JoinGroupAgainScenario = []ScenariosCommandTest{
	ConnectAlice,
	ConnectBob,
	AliceCreatesGroup,
	AliceInvitesBobInGroup,
	BobJoinAliceGroup,
	{
		Name:    "Bob accepts Alice's invitation",
		Command: "GROUP JOIN alice",
		ExpectedReplies: []Reply{
			{server.ErrAlreadyInGroup, "bob"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "bob",
	},
}

var LeaveTwiceExistantGroupScenario = []ScenariosCommandTest{
	ConnectAlice,
	ConnectBob,
	AliceCreatesGroup,
	AliceInvitesBobInGroup,
	BobJoinAliceGroup,
	AliceLeavesGroup,
	{
		Name:    "Alice tries to leave group again",
		Command: "GROUP LEAVE",
		ExpectedReplies: []Reply{
			{server.ErrNotInGroup, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var InviteUnexistantPersonScenario = []ScenariosCommandTest{
	ConnectAlice,
	AliceCreatesGroup,
	{
		Name:    "Alice tries to invite Unknown in group",
		Command: "GROUP INVITE Unknown",
		ExpectedReplies: []Reply{
			{server.ErrUserDoesntExist, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}
