package scenarios

import "testing"

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

type ConcurrentScenario struct {
	Name           string
	SetupSteps     []ScenariosCommandTest
	ConcurrentCmds map[string]string
	ValidationFunc func(t *testing.T, results map[string]string)
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
	sessionScenarioFamily,
	groupScenarioFamily,
	itemScenarioFamily,
	talkScenarioFamily,
	attackScenarioFamily,
}
