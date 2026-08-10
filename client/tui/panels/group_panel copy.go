package panel

/* import (
	"github.com/rivo/tview"
)

type GroupComponent struct {
	Layout *tview.Flex
	List   *tview.List
}

func NewGroupComponent(
	app *tview.Application,
	popup *tview.Grid,
	group string,
	invitations *[]string,
	inputs chan<- string,
	onOpenPopup func(popup *PopupComponent),
	onClosePopup func(),
) *ChoiceListComponent {
	options := ConvertGroup(group, invitations, inputs)

	src := NewChoiceListComponent(app, popup, "", options, onOpenPopup, onClosePopup, true)

	if group == "" {
		create_btn := tview.NewButton("Create").SetSelectedFunc(func() { inputs <- "GROUP CREATE" })
		src.Layout.AddItem(create_btn, 1, 0, false).SetBorder(true)
	}

	src.List.SetTitle(" Group ")

	return src
}

func ConvertGroup(group string, invitations *[]string, inputs chan<- string) OptionsMap {
	res := make(OptionsMap)
	invites := make(map[string]func())

	for _, invite := range *invitations {
		targetInvite := invite
		invites[targetInvite] = func() {
			slice := *invitations
			for i, inv := range *invitations {
				if inv == targetInvite {
					*invitations = append(slice[:i], slice[i+1:]...)
					break
				}
			}
			if group != "" {
				inputs <- "GROUP LEAVE"
			}
			inputs <- "GROUP JOIN " + targetInvite
		}
	}

	if len(invites) != 0 {
		if group != "" {
			res["LEAVE + JOIN"] = invites
		} else {
			res["JOIN"] = invites
		}
	}

	return res
}
*/
