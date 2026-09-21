package updater

func (up *Updater) UpdateGameView() {
	// Skip update if Talking
	if up.app.Variables.PanelsVariables.Talk.Talking != nil {
		up.UpdateTalk()
		return
	}

	// Init Buttons/Emotes
	if len(up.app.Manager.Buttons("Game")) == 0 {
		up.buildGameButtons()
	}
	if len(up.app.Manager.Emotes("Game")) == 0 {
		up.buildGameEmotes()
	}
	if len(up.app.Manager.Interactions("Game")) == 0 {
		up.buildGameInteractions()
	}
	if len(up.app.Manager.Buttons("Group")) == 0 {
		up.buildGroupButtons()
	}
	if len(up.app.Manager.Selects("Group")) == 0 {
		up.buildGroupSelects()
	}

	up.UpdatePlayer()

	up.app.Manager.Update("Game")
	up.app.Manager.Update("Items")
	up.app.Manager.Update("Npcs")
	if up.app.Variables.PanelsVariables.Chat.Open {
		up.UpdateChat()
	}
	up.app.Manager.Update("Inventory")
	if up.app.Variables.PanelsVariables.Group.Open {
		up.app.Manager.Update("Group")
	}
}
