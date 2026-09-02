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
	actionsChan chan<- Action,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertInteractions(npcs, players, actionsChan)

	src := NewChoiceListComponent(app, popupGrid, "Interactions", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertInteractions(npcs, players []string, actionsChan chan<- Action) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(npcs) != 0 {
		res["NPCS"] = ConvertNpcsList(npcs, actionsChan)
	}

	if len(players) != 0 {
		res["PLAYERS"] = ConvertPlayersList(players, actionsChan)
	}

	return res
}

func ConvertNpcsList(npcs []string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, npc := range npcs {
		n := npc
		res[n] = map[string]func(){
			pr.CmdTalk: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdTalk, n),
				}
			},
			pr.CmdAttack: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdAttack, n),
				}
			},
			pr.CmdQuest: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdQuest, n),
				}
			},
		}
	}

	return res
}

func ConvertPlayersList(players []string, actionsChan chan<- Action) OptionsMap {
	res := make(OptionsMap)

	for _, player := range players {
		p := player
		res[p] = map[string]func(){
			pr.CmdAttack: func() {
				actionsChan <- Action{
					Type:    ActionSendServer,
					Payload: fmt.Sprintf("%s %s", pr.CmdAttack, p),
				}
			},
		}
	}

	return res
}
