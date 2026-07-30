package scenarios

import (
	"strings"
	"testing"
)

var TakeItemRaceCondition = ConcurrentScenario{
	Name: "Concurrent TAKE on same item",

	SetupSteps: []ScenariosCommandTest{
		connectAlice,
		connectBob,
	},

	ConcurrentCmds: map[string]string{
		"alice": "TAKE sword",
		"bob":   "TAKE sword",
	},

	ValidationFunc: func(t *testing.T, results map[string]string) {
		aliceWon := strings.HasPrefix(results["alice"], "OK taken=")
		bobWon := strings.HasPrefix(results["bob"], "OK taken=")

		if aliceWon && bobWon {
			t.Fatalf("DEADLOCK/DUPE : Alice and Bob got the item both  !\nAlice: %sBob: %s", results["alice"], results["bob"])
		}

		if !aliceWon && !bobWon {
			t.Fatalf("BUG : Noone got the item !\nAlice: %sBob: %s", results["alice"], results["bob"])
		}

		t.Logf("Success : Concurrency managed correctly. Alice won: %v, Bob won: %v", aliceWon, bobWon)
	},
}
