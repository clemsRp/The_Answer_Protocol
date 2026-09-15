package updater

func (up *Updater) UpdateGameView() {
	// Init Buttons/Emotes
	// if len(up.app.Manager.Buttons("Game")) == 0 {
	// 	up.buildGameButtons()
	// }
	if len(up.app.Manager.Emotes("Game")) == 0 {
		up.buildGameEmotes()
	}

	up.UpdatePlayer()
}
