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

var aliceKickBobFromGroupScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	aliceKicksBobFromGroup,
}

var aliceKickUnknownFromGroupScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	{
		Name:    "Alice kicks unknown from group",
		Command: "GROUP KICK unknown",
		ExpectedReplies: []Reply{
			{protocol.ErrUnknownUser, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var alicePromotesBobGroupScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	alicePromotesBob,
	bobAcceptsPromotion,
}

var aliceInvitesBobTwiceScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	{
		Name:    "Alice invites Bob in group AGAIN",
		Command: "GROUP INVITE bob",
		ExpectedReplies: []Reply{
			{protocol.ErrInvitationAlreadySent, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var bobJoinsGhostGroupScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	aliceLeavesGroupSolo,
	bobJoinsGhostGroup,
}

var leaderKicksHimselfScenario = []ScenariosCommandTest{
	connectAlice,
	aliceCreatesGroup,
	aliceKicksHerself,
}

var promoteBobThenCarlScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	connectCarl,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	aliceInvitesCarlInGroup,
	carlJoinAliceGroup,
	alicePromotesBob,
	bobLeavesBigGroup,
	alicePromotesCarl,
	carlAcceptsPromotion,
}

var promoteNonMemberScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	alicePromotesBobNotInGroup,
}

var declinePromotionScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	aliceInvitesBobInGroup,
	bobJoinAliceGroup,
	alicePromotesBob,
	bobDeclinesPromotion,
}

var joinWithoutInviteScenario = []ScenariosCommandTest{
	connectAlice,
	connectBob,
	aliceCreatesGroup,
	bobJoinsWithoutInvite,
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
		{
			Name:  "Alice kicks bob from group",
			Steps: aliceKickBobFromGroupScenario,
		},
		{
			Name:  "Alice tries to kick unknown from group",
			Steps: aliceKickUnknownFromGroupScenario,
		},
		{
			Name:  "Alice promotes bob",
			Steps: alicePromotesBobGroupScenario,
		},
		{
			Name:  "Alice invites bob twice in group",
			Steps: aliceInvitesBobTwiceScenario,
		},
		{
			Name:  "Bob tries to join a deleted group",
			Steps: bobJoinsGhostGroupScenario,
		},
		{
			Name:  "Leader kicks himself leaving a ghost group",
			Steps: leaderKicksHimselfScenario,
		},
		{
			Name:  "Promte bob then carl in group",
			Steps: promoteBobThenCarlScenario,
		},
		{
			Name:  "Alice tries to promote a non-member",
			Steps: promoteNonMemberScenario,
		},
		{
			Name:  "Bob declines promotion gracefully",
			Steps: declinePromotionScenario,
		},
		{
			Name:  "Carl tries to join without an invitation",
			Steps: joinWithoutInviteScenario,
		},
	},
}
