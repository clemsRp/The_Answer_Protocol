package panel

import (
	"fmt"
	"strings"
	pr "tap/protocol"

	"github.com/rivo/tview"
)

const npcDialoguePrefix = "   L "

type AvailableQuestData struct {
	Id     string
	Status string
}

var quest_datas map[string]bool

func NewInteractionComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	npcs,
	players []string,
	npcData map[string]pr.InspectNPCData,
	npcDialogues map[string]string,
	groupMembers []string,
	panel_width int,
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
	quests *[]pr.TrackedQuestData,
) *ChoiceListComponent {
	// Update quest ids
	if quests != nil {
		quest_datas = make(map[string]bool)
		for _, quest := range *quests {
			quest_datas[quest.Id] = quest.Status == "completed"
		}
	}

	options := ConvertInteractions(npcs, players, npcData, npcDialogues, groupMembers, actionsChan)

	src := NewChoiceListComponent(app, popupGrid, "Interactions", options, onOpenPopup, onClosePopup, false)

	attachNpcDialogues(src.List, npcs, npcDialogues, panel_width)

	return src
}

func dialogueWrapWidth(panel_width int) int {
	const margin = 4

	width := panel_width - len(npcDialoguePrefix) - margin
	if width < 20 {
		width = 20
	}
	return width
}

func wrapDialogue(text string, width int) []string {
	if width <= 0 {
		width = 60
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	lines := make([]string, 0, 4)
	current := words[0]
	for _, w := range words[1:] {
		if len(current)+1+len(w) > width {
			lines = append(lines, current)
			current = w
			continue
		}
		current += " " + w
	}
	lines = append(lines, current)

	return lines
}

func attachNpcDialogues(list *tview.List, npcs []string, npcDialogues map[string]string, panel_width int) {
	if list == nil || len(npcDialogues) == 0 {
		return
	}

	list.SetSecondaryTextColor(AppTheme.TextSecondary)

	wrapWidth := dialogueWrapWidth(panel_width)
	indent := strings.Repeat(" ", len(npcDialoguePrefix))

	for _, npc := range npcs {
		dialogue := strings.TrimSpace(npcDialogues[npc])
		if dialogue == "" {
			continue
		}

		idx, ok := findNpcItemIndex(list, npc)
		if !ok {
			continue
		}

		lines := wrapDialogue(dialogue, wrapWidth)
		if len(lines) == 0 {
			continue
		}

		mainText, _ := list.GetItemText(idx)
		list.SetItemText(idx, mainText, npcDialoguePrefix+lines[0])

		for j := 1; j < len(lines); j++ {
			list.InsertItem(idx+j, "", indent+lines[j], 0, nil)
		}
	}
}

func findNpcItemIndex(list *tview.List, npc string) (int, bool) {
	for i := 0; i < list.GetItemCount(); i++ {
		mainText, _ := list.GetItemText(i)
		if strings.TrimSpace(mainText) == npc {
			return i, true
		}
	}
	return -1, false
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
				completed, available := quest_datas[data.QuestId]

				if !available {
					actions[pr.CmdQuest] = func() {
						actionsChan <- Action{
							Type:    ActionSendServer,
							Payload: fmt.Sprintf("%s %s", pr.CmdQuest, npc),
						}
					}

				} else if !completed {
					actions["COMPLETE QUEST"] = func() {
						actionsChan <- Action{
							Type:    ActionSendServer,
							Payload: fmt.Sprintf("%s %s", pr.CmdCompleteQuest, data.QuestId),
						}
					}

				}
			}
		}

		res[npc] = actions
	}

	return res
}

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
