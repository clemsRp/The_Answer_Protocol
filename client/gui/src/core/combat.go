package core

func (app *App) IsPlayerTurn() bool {
	cs := app.Variables.PanelsVariables.CombatState
	return cs != nil && cs.CurrentTurn != "" && cs.CurrentTurn == app.Variables.Player.Pseudo
}
