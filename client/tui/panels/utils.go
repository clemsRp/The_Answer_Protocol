package panel

import (
	"strings"
	pr "tap/protocol"
)

func IsErrorResponse(res pr.ServerResponse) bool {
	if strings.HasPrefix(res.Msg, "ERR") {
		return true
	}
	return false
}

func IsEventResponse(res pr.ServerResponse) bool {
	if strings.HasPrefix(res.Msg, "EVT") {
		return true
	}
	return false
}

func IsOKResponse(res pr.ServerResponse) bool {
	if strings.HasPrefix(res.Msg, "OK") {
		return true
	}
	return false
}
