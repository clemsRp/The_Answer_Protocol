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
	groupMembers []string,
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertInteractions(npcs, players, npcData, npcDialogues, groupMembers, actionsChan)

	src := NewChoiceListComponent(app, popupGrid, "Interactions", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertInteractions(npcs, players []string, npcData map[string]pr.InspectNPCData, npcDialogues map[string]string, groupMembers []string, actionsChan chan<- Action) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(npcs) != 0 {
		res["NPCS"] = ConvertNpcsList(npcs, npcData, npcDialogues, actionsChan)
	}

	if len(players) != 0 {
		res["PLAYERS"] = ConvertPlayersList(players, groupMembers, actionsChan)
	}

	return res
}

// ConvertNpcsList builds the popup actions for each npc in the room. TALK and
// INSPECT are always available; ATTACK only appears once we know (via
// INSPECT) that the npc is hostile, and QUEST only if it has a quest to give.
func ConvertNpcsList(npcs []string, npcData map[string]pr.InspectNPCData, npcDialogues map[string]string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, npc := range npcs {
		actions := map[string]func(){
			pr.CmdTalk: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdTalk, npc),
				}
			},
			pr.CmdInspect: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s %s", pr.CmdInspect, pr.EntityTypeNpc, npc),
				}
			},
		}

		if data, ok := npcData[npc]; ok {
			if data.Hostile {
				actionName := pr.CmdAttack
				if data.InCombat {
					actionName = "JOIN COMBAT"
				}
				actions[actionName] = func() {
					actionsChan <- Action{
						Type:    ActionSendServer,
						Payload: fmt.Sprintf("%s %s", pr.CmdAttack, npc),
					}
				}
			}
			if data.QuestId != "" {
				actions[pr.CmdQuest] = func() {
					actionsChan <- Action{
						Type:    ActionSendServer,
						Payload: fmt.Sprintf("%s %s", pr.CmdQuest, npc),
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

		res[npc] = actions
	}

	return res
}

// ConvertPlayersList builds the popup actions for each player in the room.
// "JOIN COMBAT" is only shown for players who are in the same group as us
// (i.e. present in groupMembers). All players always have INSPECT.
func ConvertPlayersList(players []string, groupMembers []string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	groupSet := make(map[string]bool, len(groupMembers))
	for _, m := range groupMembers {
		groupSet[m] = true
	}

	for _, player := range players {
		p := player
		actions := map[string]func(){
			pr.CmdInspect: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s %s", pr.CmdInspect, pr.EntityTypePlayer, p),
				}
			},
		}
		// Only group members can be joined in combat
		if groupSet[p] {
			actions["JOIN COMBAT"] = func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdAttack, p),
				}
			}
		}
		res[p] = actions
	}

	return res
}
