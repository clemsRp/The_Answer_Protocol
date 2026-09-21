package updater

func (up *Updater) buildGroupSelects() {
	up.app.Manager.SetViewSelects("Group", up.app.GetGroupSelects())
}
