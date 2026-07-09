package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type QuestComponent struct {
	Layout *tview.Flex
	List   *tview.List
}

func NewQuestComponent(app *tview.Application) *QuestComponent {
	src := &QuestComponent{}

	src.List = createListView(" Active Quest ", true, Default, Default, tcell.ColorOrange, Black, false, true)

	src.List.AddItem("⏳ No active quests...", "", 0, nil)

	src.Layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(src.List, 0, 1, false)

	return src
}

func (c *QuestComponent) ListenOutputs(app *tview.Application, questChan <-chan string, inputs chan<- string) {
	go func() {
		for msg := range questChan {
			app.QueueUpdateDraw(func() {
				c.List.Clear()

				if strings.HasPrefix(msg, "ERROR:") {
					c.List.AddItem("Failed to load quests.", "", 0, nil)
					return
				}

				quests := strings.Split(msg, ",")
				for _, quest := range quests {
					q := quest
					c.List.AddItem(q, "", 0, func() {
						inputs <- "TRACK_QUEST:" + q
					})
				}
			})
		}
	}()
}
