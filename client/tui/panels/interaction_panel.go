package panel

import (
	"fmt"

	"github.com/rivo/tview"
)

func NewInteractionComponent(
	app *tview.Application,
	popupGrid *tview.Grid,
	npcs,
	players []string,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertInteractions(npcs, players, inputs)

	src := NewChoiceListComponent(app, popupGrid, "Interactions", options, onOpenPopup, onClosePopup, false)

	return src
}

func ConvertInteractions(npcs, players []string, inputs chan<- string) map[string]OptionsMap {
	res := make(map[string]OptionsMap)

	if len(npcs) != 0 {
		res["NPCS"] = ConvertInteractionsList(npcs, []string{"TALK", "ATTACK"}, inputs)
	}

	if len(players) != 0 {
		res["PLAYERS"] = ConvertInteractionsList(players, []string{"ATTACK"}, inputs)
	}

	return res
}

func ConvertInteractionsList(people, cmds []string, inputs chan<- string) OptionsMap {
	res := make(OptionsMap)

	for _, person := range people {
		for _, cmd := range cmds {

			res[person] = map[string]func(){
				cmd: func() {
					go func() {
						inputs <- fmt.Sprintf("%s %s", cmd, person)
					}()
				},
			}
		}
	}

	return res
}
