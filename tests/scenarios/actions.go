package scenarios

import (
	"tap/protocol"
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
		{protocol.EventStatsPlayers, "alice"},
		{protocol.EventRoomPresenceEnter, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}
var connectCarl = ScenariosCommandTest{
	Name:    "CONNECT user",
	Command: "CONNECT carl",
	ExpectedReplies: []Reply{
		{"OK connected", "carl"},
		{protocol.EventStatsPlayers, "carl"},
		{protocol.EventRoomPresenceEnter, "carl"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "carl",
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
		{protocol.EventGroupInvite, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var bobJoinAliceGroup = ScenariosCommandTest{
	Name:    "Bob accepts Alice's invitation",
	Command: "GROUP JOIN alice",
	ExpectedReplies: []Reply{
		{"OK group=", "bob"},
		{protocol.EventGroupJoin + " bob", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var aliceLeavesGroup = ScenariosCommandTest{
	Name:    "Alice leaves group",
	Command: "GROUP LEAVE",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{protocol.EventGroupLeave + " alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceLeavesGroupSolo = ScenariosCommandTest{
	Name:    "Alice leaves group SOLO",
	Command: "GROUP LEAVE",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceKicksBobFromGroup = ScenariosCommandTest{
	Name:    "Alice kicks Bob from group",
	Command: "GROUP KICK bob",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{protocol.EventGroupKicked, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}
var alicePromotesBob = ScenariosCommandTest{
	Name:    "alice promotes bob",
	Command: "GROUP PROMOTE bob",
	ExpectedReplies: []Reply{
		{protocol.GroupPromoteCmdResponse + "bob", "alice"},
		{protocol.EventGroupPromote, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var bobAcceptsPromotion = ScenariosCommandTest{
	Name:    "bob accepts promotion",
	Command: "GROUP ACCEPT",
	ExpectedReplies: []Reply{
		{protocol.GroupAcceptPromotionResponse, "bob"},
		{protocol.EventGroupNewLeader + "bob", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var aliceKicksHerself = ScenariosCommandTest{
	Name:    "Alice kicks herself from her own group",
	Command: "GROUP KICK alice",
	ExpectedReplies: []Reply{
		{"OK", "alice"}, 
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var bobJoinsGhostGroup = ScenariosCommandTest{
	Name:    "Bob tries to join a group that Alice just deleted by leaving",
	Command: "GROUP JOIN alice",
	ExpectedReplies: []Reply{
		{protocol.ErrGroupDoesntExist, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var bobLeavesGroup = ScenariosCommandTest{
	Name:    "Bob leaves group",
	Command: "GROUP LEAVE",
	ExpectedReplies: []Reply{
		{"OK", "bob"},
		{protocol.EventGroupLeave + " bob", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var aliceInvitesCarlInGroup = ScenariosCommandTest{
	Name:    "Alice invites Carl in group",
	Command: "GROUP INVITE carl",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{protocol.EventGroupInvite, "carl"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var carlJoinAliceGroup = ScenariosCommandTest{
	Name:    "Carl accepts Alice's invitation",
	Command: "GROUP JOIN alice",
	ExpectedReplies: []Reply{
		{"OK group=", "carl"},
		{protocol.EventGroupJoin + " carl", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "carl",
}

var alicePromotesCarlBlocked = ScenariosCommandTest{
	Name:    "Alice tries to promote Carl but is blocked by a previous unresolved promotion",
	Command: "GROUP PROMOTE carl",
	ExpectedReplies: []Reply{
		{protocol.ErrNoPermission, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

// MOVE NORMAL SCENARIOS
var aliceInvalidMovement = ScenariosCommandTest{
	Name:    "Alice invalid movement",
	Command: "MOVE nowhere",
	ExpectedReplies: []Reply{
		{protocol.ErrNoExit, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesEast = ScenariosCommandTest{
	Name:    "Alice valid movement east",
	Command: "MOVE east",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{protocol.EventRoomPresenceLeave + " alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesWest = ScenariosCommandTest{
	Name:    "Alice valid movement west",
	Command: "MOVE west",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{protocol.EventRoomPresenceLeave + " alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesNorth = ScenariosCommandTest{
	Name:    "Alice valid movement north",
	Command: "MOVE north",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{protocol.EventRoomPresenceLeave + " alice", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesToHealthAisle = ScenariosCommandTest{
	Name:    "Alice moves to the health aisle",
	Command: "MOVE north",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceMovesSouth = ScenariosCommandTest{
	Name:    "Alice valid movement south",
	Command: "MOVE south",
	ExpectedReplies: []Reply{
		{"OK room=", "alice"},
		{protocol.EventRoomPresenceLeave + " alice", "bob"},
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
		{protocol.EventGlobalChat + " alice Hello World!", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceChatRoom = ScenariosCommandTest{
	Name:    "Alice room chat",
	Command: "CHAT ROOM Hello Room!",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{protocol.EventRoomChat + " alice Hello Room!", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceChatGroup = ScenariosCommandTest{
	Name:    "Alice group chat",
	Command: "CHAT GROUP Hello Team!",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{protocol.EventGroupChat + " alice Hello Team!", "bob"},
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
		{protocol.EventRoomPresenceLeave + " alice", "bob"},
		{protocol.EventStatsPlayers, "bob"},
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
	Command: "TALK grandpa_gaston",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceAttacksHostileNPC = ScenariosCommandTest{

	Name:    "Attack hostile NPC",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatTurn + "Grandpa Gaston", "alice"},
		{protocol.EventCombatTurn + "alice", "alice"},
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceAttacksHostileNPCAgain = ScenariosCommandTest{
	Name:    "Alice attacks hostile NPC on her first turn",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatTurn + "Grandpa Gaston", "alice"},
		{protocol.EventCombatTurn + "alice", "alice"},
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceDefeatsHostileNPC = ScenariosCommandTest{
	Name:    "Alice defeats hostile NPC",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatVictory, "alice"},
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceStartsGroupCombat = ScenariosCommandTest{
	Name:    "Alice starts combat with Bob",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatStarted + "alice", "alice"},
		{protocol.EventCombatTurn + "Grandpa Gaston", "alice"},
		{protocol.EventCombatTurn + "alice", "alice"},
		{"OK", "alice"},
		{protocol.EventCombatStarted + "alice", "bob"},
		{protocol.EventCombatTurn + "Grandpa Gaston", "bob"},
		{protocol.EventCombatTurn + "alice", "bob"},
		{protocol.EventCombatUpdate, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}
var aliceStartsAnotherGroupCombat = ScenariosCommandTest{
	Name:    "Alice starts another combat while Bob is in combat in the same room",
	Command: "ATTACK samus",
	ExpectedReplies: []Reply{
		{protocol.EventCombatStarted + "alice", "alice"},
		{protocol.EventCombatTurn + "samus", "alice"},
		{protocol.EventCombatTurn + "alice", "alice"},
		{protocol.EventDistantGroupCombatStartedCombat + "alice", "bob"},
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceStartsGroupCombatAgainstKillerAndDiesToRespawnWithLessHp = ScenariosCommandTest{
	Name:    "Alice starts fatal combat vs Killer",
	Command: "ATTACK killer",
	ExpectedReplies: []Reply{
		{protocol.EventCombatStarted + "alice", "alice"},
		{protocol.EventCombatTurn + "Killer", "alice"},
		{protocol.EventCombatDefeat + protocol.RoomEntrance, "alice"},
		{"OK", "alice"},
		{protocol.EventCombatStarted + "alice", "bob"},
		{protocol.EventCombatTurn + "Killer", "bob"},
		{protocol.EventCombatPlayerDied + "alice", "bob"},
		{protocol.EventCombatTurn + "bob", "bob"},
		{protocol.EventCombatUpdate, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceStartsCombatAgainstKillerAndDiesToRespawnWithLessHp = ScenariosCommandTest{
	Name:    "Alice starts fatal combat vs Killer",
	Command: "ATTACK killer",
	ExpectedReplies: []Reply{
		{protocol.EventCombatStarted + "alice", "alice"},
		{protocol.EventCombatTurn + "Killer", "alice"},
		{protocol.EventCombatDefeat + protocol.RoomEntrance, "alice"},
		{"OK", "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceAttacksInGroupCombat = ScenariosCommandTest{
	Name:    "Alice attacks NPC while in groupped combat with bob",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatTurn + "bob", "alice"},
		{"OK", "alice"},
		{protocol.EventCombatTurn + "bob", "bob"},
		{protocol.EventCombatUpdate, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}
var bobAttacksInGroupCombat = ScenariosCommandTest{
	Name:    "Bob attacks NPC while in groupped combat with Alice",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatTurn + "Grandpa Gaston", "bob"},
		{protocol.EventCombatTurn + "alice", "bob"},
		{"OK", "bob"},
		{protocol.EventCombatTurn + "Grandpa Gaston", "alice"},
		{protocol.EventCombatTurn + "alice", "alice"},
		{protocol.EventCombatUpdate, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var aliceFleesCombat = ScenariosCommandTest{
	Name:    "Alice flees combat",
	Command: "FLEE",
	ExpectedReplies: []Reply{
		{"OK", "alice"},
		{protocol.EventCombatTurn + "bob", "bob"},
		{protocol.EventCombatAllyLeaveCombat + "alice", "bob"},
		{protocol.EventCombatUpdate, "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var bobAttacksAfterAliceFlees = ScenariosCommandTest{
	Name:    "Bob attacks after Alice flees",
	Command: "ATTACK grandpa_gaston",
	ExpectedReplies: []Reply{
		{protocol.EventCombatTurn + "Grandpa Gaston", "bob"},
		{protocol.EventCombatTurn + "bob", "bob"},
		{"OK", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var bobFleesLastFromCombat = ScenariosCommandTest{
	Name:    "Bob flees as the last combatant",
	Command: "FLEE",
	ExpectedReplies: []Reply{
		{"OK", "bob"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "bob",
}

var aliceMovesInCombatForbidden = ScenariosCommandTest{
	Name:    "Move while in combat",
	Command: "MOVE NORTH",
	ExpectedReplies: []Reply{
		{protocol.ErrForbiddenInCombat, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceTakesInCombatForbidden = ScenariosCommandTest{
	Name:    "Take while in combat",
	Command: "TAKE sword",
	ExpectedReplies: []Reply{
		{protocol.ErrForbiddenInCombat, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceDropsInCombatForbidden = ScenariosCommandTest{
	Name:    "Drop while in combat",
	Command: "DROP gold_denture",
	ExpectedReplies: []Reply{
		{protocol.ErrForbiddenInCombat, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceQuestInCombatForbidden = ScenariosCommandTest{
	Name:    "QUEST while in combat",
	Command: "QUEST quest_giver",
	ExpectedReplies: []Reply{
		{protocol.ErrForbiddenInCombat, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceTalksInCombatForbidden = ScenariosCommandTest{
	Name:    "Talk while in combat",
	Command: "TALK quest_giver",
	ExpectedReplies: []Reply{
		{protocol.ErrForbiddenInCombat, "alice"},
	},
	ExpectsJSON:      false,
	TestOnConnection: "alice",
}

var aliceForbiddenActionsInCombat = []ScenariosCommandTest{
	aliceMovesInCombatForbidden,
	aliceTakesInCombatForbidden,
	aliceDropsInCombatForbidden,
	aliceQuestInCombatForbidden,
	aliceTalksInCombatForbidden,
}
