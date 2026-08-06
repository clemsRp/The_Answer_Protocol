package scenarios

import (
	"testing"
)

func TestScenariosCommands(t *testing.T) {
	for _, family := range OrderedScenarioFamilies {

		t.Run(family.FamilyName, func(t *testing.T) {

			for _, scenarioEntry := range family.Scenarios {
				t.Run(scenarioEntry.Name, func(t *testing.T) {
					RunScenario(t, scenarioEntry)
				})
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	t.Run(TakeItemRaceCondition.Name, func(t *testing.T) {
		RunConcurrentScenario(t, TakeItemRaceCondition)
	})
}
