package panel

import (
	"fmt"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

func NewInteractionComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	npcs,
	players []string,
	npcData map[string]pr.InspectNPCData,
	npcDialogues map[string]string,
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertInteractions(npcs, players, npcData, npcDialogues, actionsChan)

	src := NewChoiceListComponent(app, popupGrid, "Interactions", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertInteractions(npcs, players []string, npcData map[string]pr.InspectNPCData, npcDialogues map[string]string, actionsChan chan<- Action) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(npcs) != 0 {
		res["NPCS"] = ConvertNpcsList(npcs, npcData, npcDialogues, actionsChan)
	}

	if len(players) != 0 {
		res["PLAYERS"] = ConvertPlayersList(players, actionsChan)
	}

	return res
}

// ConvertNpcsList builds the popup actions for each npc in the room. TALK and
// INSPECT are always available; ATTACK only appears once we know (via
// INSPECT) that the npc is hostile, and QUEST only if it has a quest to give.
func ConvertNpcsList(npcs []string, npcData map[string]pr.InspectNPCData, npcDialogues map[string]string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, npc := range npcs {
		n := npc

		// Setup the text with dialogue in gray as secondary text.
		displayName := n
		if dialogue, exists := npcDialogues[n]; exists && dialogue != "" {
			displayName = n + "\n[gray]" + dialogue + "[-]"
		}

		actions := map[string]func(){
			pr.CmdTalk: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdTalk, n),
				}
			},
			pr.CmdInspect: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s %s", pr.CmdInspect, pr.EntityTypeNpc, n),
				}
			},
		}

		if data, ok := npcData[n]; ok {
			if data.Hostile {
				actionName := pr.CmdAttack
				if data.InCombat {
					actionName = "JOIN COMBAT"
				}
				actions[actionName] = func() {
					actionsChan <- Action{
						Type:    ActionSendServer,
						Payload: fmt.Sprintf("%s %s", pr.CmdAttack, n),
					}
				}
			}
			if data.QuestId != "" {
				actions[pr.CmdQuest] = func() {
					actionsChan <- Action{
						Type:    ActionSendServer,
						Payload: fmt.Sprintf("%s %s", pr.CmdQuest, n),
					}
				}
				actions["COMPLETE QUEST"] = func() {
					actionsChan <- Action{
						Type:    ActionSendServer,
						Payload: fmt.Sprintf("%s %s", pr.CmdCompleteQuest, data.QuestId),
					}
				}
			}
		}

		res[displayName] = actions
	}

	return res
}

func ConvertPlayersList(players []string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, player := range players {
		p := player
		res[p] = map[string]func(){
			"JOIN COMBAT": func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdAttack, p),
				}
			},
			pr.CmdInspect: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s %s", pr.CmdInspect, pr.EntityTypePlayer, p),
				}
			},
		}
	}

	return res
}
