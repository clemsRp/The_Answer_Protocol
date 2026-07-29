package tests

import (
	sc "tap/tests/scenarios"
	"tap/tests/utils"
	"testing"
)

func TestScenariosCommands(t *testing.T) {
	for _, family := range sc.OrderedScenarioFamilies {

		t.Run(family.FamilyName, func(t *testing.T) {

			for _, scenarioEntry := range family.Scenarios {
				utils.RunScenario(t, scenarioEntry)
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	t.Run(sc.TakeItemRaceCondition.Name, func(t *testing.T) {
		utils.RunConcurrentScenario(t, sc.TakeItemRaceCondition)
	})
}
