package scenarios

type Reply struct {
	Msg  string
	User string
}

type ScenariosCommandTest struct {
	Name             string
	Command          string
	ExpectedReplies  []Reply
	ExpectsJSON      bool
	TestOnConnection string
}

type ScenarioEntry struct {
	Name  string
	Steps []ScenariosCommandTest
}

type ScenarioFamily struct {
	FamilyName string
	Scenarios  []ScenarioEntry
}

var OrderedScenarioFamilies = []ScenarioFamily{
	{
		FamilyName: "Connection_Family",
		Scenarios: []ScenarioEntry{
			{"Connect on same user again", UsernameAlreadyUsedScenario},
		},
	},
	{
		FamilyName: "Group_Family",
		Scenarios: []ScenarioEntry{
			{"Join Group Again", JoinGroupAgainScenario},
			{"Leave Twice Existant Group", LeaveTwiceExistantGroupScenario},
			{"Leave Unexistant Group", LeaveUnexistantGroupScenario},
			{"Invite Unexistant Person", InviteUnexistantPersonScenario},
		},
	},
}
