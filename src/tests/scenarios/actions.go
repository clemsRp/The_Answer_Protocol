package scenarios

import (
	"tap/src/protocol"
)

var connectAlice = ScenariosCommandTest{
	Name:    "CONNECT user",
	Command: "CONNECT alice",
	ExpectedReplies: []Reply{
		{"OK connected", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var connectBob = ScenariosCommandTest{
	Name:    "CONNECT user",
	Command: "CONNECT bob",
	ExpectedReplies: []Reply{
		{"OK connected", "bob"},
		{"EVT STATS players=", "alice"},
		{"EVT ROOM PRESENCE ENTER", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

// GROUP SCENARIOS

var aliceCreatesGroup = ScenariosCommandTest{
	Name:    "Alice creates a group",
	Command: "GROUP CREATE",
	ExpectedReplies: []Reply{
		{"OK group=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceInvitesBobInGroup = ScenariosCommandTest{
	Name:    "Alice invites Bob in group",
	Command: "GROUP INVITE bob",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{"EVT GROUP INVITE", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var bobJoinAliceGroup = ScenariosCommandTest{
	Name:    "Bob accepts Alice's invitation",
	Command: "GROUP JOIN alice",
	ExpectedReplies: []Reply{
		{"OK group=", "bob"},
		{"EVT GROUP JOIN bob", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}
var aliceLeavesGroup = ScenariosCommandTest{
	Name:    "Alice leaves group",
	Command: "GROUP LEAVE",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{"EVT GROUP LEAVE alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

// BASIC COMMANDS NORMAL SCENARIOS
var aliceTakesItem = ScenariosCommandTest{
	Name:    "Alice takes an item in map",
	Command: "TAKE pioche",
	ExpectedReplies: []Reply{
		{"OK taken=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceDropsItem = ScenariosCommandTest{
	Name:    "Alice drops an item in map",
	Command: "DROP pioche",
	ExpectedReplies: []Reply{
		{"OK dropped=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceTalksToNPC = ScenariosCommandTest{
	Name:    "Alice talks to NPC in entrance",
	Command: "TALK thomas_charbonneau",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceAttacksHostileNPC = ScenariosCommandTest{

	Name:    "Attack non-hostile NPC in entrance",
	Command: "ATTACK thomas_charbonneau",
	ExpectedReplies: []Reply{
		{protocol.ErrNpcNotHostile, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}
