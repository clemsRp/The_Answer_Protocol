package panel

import (
	"strings"
	pr "tap/src/protocol"
)

const (
	ActionSendServer      = "SEND_SERVER"
	ActionNavigate        = "NAVIGATE"
	ActionQuit            = "QUIT"
	ActionOpenPopUp       = "OPEN_POPUP"
	ActionClosePopUp      = "CLOSE_POPUP"
	ActionOpenGroupInfo   = "OPEN_GROUP_INFO"
	ActionOpenGroupInvite = "OPEN_GROUP_INVITE"
	ActionOpenCombat      = "OPEN_COMBAT"
	ActionCloseCombat     = "CLOSE_COMBAT"
)

type Action struct {
	Type    string
	Payload any
}

func IsErrorResponse(res pr.ServerResponse) bool {
	return strings.HasPrefix(res.Msg, "ERR")
}

func IsEventResponse(res pr.ServerResponse) bool {
	return strings.HasPrefix(res.Msg, "EVT")
}

func IsOKResponse(res pr.ServerResponse) bool {
	return strings.HasPrefix(res.Msg, "OK")
}
