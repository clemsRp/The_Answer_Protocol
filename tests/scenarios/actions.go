package scenarios

import (
	"tap/server"
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

// MOVE NORMAL SCENARIOS
var aliceInvalidMovement = ScenariosCommandTest{
	Name:    "Alice invalid movement",
	Command: "MOVE nowhere",
	ExpectedReplies: []Reply{
		{server.ErrNoExit, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesEast = ScenariosCommandTest{
	Name:    "Alice valid movement east",
	Command: "MOVE east",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{"EVT ROOM PRESENCE LEAVE alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesWest = ScenariosCommandTest{
	Name:    "Alice valid movement west",
	Command: "MOVE west",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{"EVT ROOM PRESENCE LEAVE alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesNorth = ScenariosCommandTest{
	Name:    "Alice valid movement north",
	Command: "MOVE north",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{"EVT ROOM PRESENCE LEAVE alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesSouth = ScenariosCommandTest{
	Name:    "Alice valid movement south",
	Command: "MOVE south",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{"EVT ROOM PRESENCE LEAVE alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

// CHAT NORMAL SCENARIOS
var aliceChatGlobal = ScenariosCommandTest{
	Name:    "Alice global chat",
	Command: "CHAT GLOBAL Hello World!",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{"EVT GLOBAL CHAT alice Hello World!", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceChatRoom = ScenariosCommandTest{
	Name:    "Alice room chat",
	Command: "CHAT ROOM Hello Room!",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{"EVT ROOM CHAT alice Hello Room!", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceChatGroup = ScenariosCommandTest{
	Name:    "Alice group chat",
	Command: "CHAT GROUP Hello Team!",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{"EVT GROUP CHAT alice Hello Team!", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

// BASIC COMMANDS NORMAL SCENARIOS

var aliceLooks = ScenariosCommandTest{
	Name:    "Alice look command",
	Command: "LOOK",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
	},
	ExpectsJSON:      true,
	TestOnConnection: "alice",
}

var aliceWho = ScenariosCommandTest{
	Name:    "Alice who command",
	Command: "WHO",
	ExpectedReplies: []Reply{
		{"OK players=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceQuits = ScenariosCommandTest{
	Name:    "Alice quits the game",
	Command: "QUIT",
	ExpectedReplies: []Reply{
		{"OK bye", "alice"},
		{"EVT ROOM PRESENCE LEAVE alice", "bob"},
		{"EVT STATS players=", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceTakesItem = ScenariosCommandTest{
	Name:    "Alice takes an item in map",
	Command: "TAKE sword",
	ExpectedReplies: []Reply{
		{"OK taken=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceDropsItem = ScenariosCommandTest{
	Name:    "Alice drops an item in map",
	Command: "DROP sword",
	ExpectedReplies: []Reply{
		{"OK dropped=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceChecksInventory = ScenariosCommandTest{
	Name:    "Alice checks inventory",
	Command: "INVENTORY",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
	},
	ExpectsJSON:      true,
	TestOnConnection: "alice",
}

var aliceTalksToNPC = ScenariosCommandTest{
	Name:    "Alice talks to NPC in entrance",
	Command: "TALK granny jeanine",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}
