package updater

func (up *Updater) UpdateGameView() {
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

	up.UpdatePlayer()

	up.app.Manager.Update("Game")
	if up.app.Variables.PanelsVariables.Chat.Open {
		up.UpdateChat()
	}
	up.app.Manager.Update("Inventory")
}
