package scenarios

import (
	"tap/server"
	"testing"
)

func GetJoinGroupAgainScenario(t *testing.T) []ScenariosCommandTest {
	t.Helper()

	return []ScenariosCommandTest{
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
}

func GetLeaveUnexistantGroupScenario(t *testing.T) []ScenariosCommandTest {
	t.Helper()

	return []ScenariosCommandTest{
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
}

func GetLeaveTwiceExistantGroupScenario(t *testing.T) []ScenariosCommandTest {
	t.Helper()

	return []ScenariosCommandTest{
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
}

func GetInviteUnexistantPersonScenario(t *testing.T) []ScenariosCommandTest {
	t.Helper()

	return []ScenariosCommandTest{
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
}

func GetUsernameAlreadyUsedScenario(t *testing.T) []ScenariosCommandTest {
	t.Helper()
	return []ScenariosCommandTest{
		{
			Name:    "Invalid connection: Name already used",
			Command: "CONNECT alice",
			ExpectedReplies: []Reply{
				{server.ErrNameInUse, "alice"},
			},
			ExpectsJSON: false,
		},
	}
}
