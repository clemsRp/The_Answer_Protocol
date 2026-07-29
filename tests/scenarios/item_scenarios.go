package scenarios

import "tap/protocol"

var itemTakeScenario = []ScenariosCommandTest{
	connectAlice,
	aliceTakesItem,
}

var itemNotFoundScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Tries to take same item",
		Command: "TAKE clement",
		ExpectedReplies: []Reply{
			{protocol.ErrItemNotFound, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var itemTakeTwiceScenario = []ScenariosCommandTest{
	connectAlice,
	aliceTakesItem,
	{
		Name:    "Tries to take same item",
		Command: "TAKE sword",
		ExpectedReplies: []Reply{
			{protocol.ErrItemNotFound, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var itemDropScenario = []ScenariosCommandTest{
	connectAlice,
	aliceTakesItem,
	aliceDropsItem,
}

var itemDroppedNotInInventoryScenario = []ScenariosCommandTest{
	connectAlice,
	{
		Name:    "Tries to drop item",
		Command: "DROP sword",
		ExpectedReplies: []Reply{
			{protocol.ErrItemNotInInventory, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var itemDropTwiceScenario = []ScenariosCommandTest{
	connectAlice,
	aliceTakesItem,
	aliceDropsItem,
	{
		Name:    "Tries to drop item",
		Command: "DROP sword",
		ExpectedReplies: []Reply{
			{protocol.ErrItemNotInInventory, "alice"},
		},
		ExpectsJSON:      false,
		TestOnConnection: "alice",
	},
}

var itemTakeAndCheckInventoryScenario = []ScenariosCommandTest{
	connectAlice,
	aliceTakesItem,
	{
		Name:    "Alice checks inventory if taken item is in there",
		Command: "INVENTORY",
		ExpectedReplies: []Reply{
			{"OK", "alice"},
		},
		ExpectsJSON: true,
	},
}

var itemScenarioFamily = ScenarioFamily{
	FamilyName: "Item scenarios family",
	Scenarios: []ScenarioEntry{
		{
			Name:  "Take item successfully",
			Steps: itemTakeScenario,
		},
		{
			Name:  "Try to take a non-existent item",
			Steps: itemNotFoundScenario,
		},
		{
			Name:  "Try to take the same item twice",
			Steps: itemTakeTwiceScenario,
		},
		{
			Name:  "Drop item successfully",
			Steps: itemDropScenario,
		},
		{
			Name:  "Try to drop an item not in inventory",
			Steps: itemDroppedNotInInventoryScenario,
		},
		{
			Name:  "Try to drop the same item twice",
			Steps: itemDropTwiceScenario,
		},
		{
			Name:  "Try to take an item and check inventory",
			Steps: itemTakeAndCheckInventoryScenario,
		},
	},
}
