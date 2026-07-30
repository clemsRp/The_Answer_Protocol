package panel

import (
	pr "tap/protocol"
)

const (
	COLOR_ERR = "red"
	COLOR_EVT = "orange"
	COLOR_OK  = "green"
)

func GetResponseColor(res pr.ServerResponse) string {
	color := "green"
	if IsErrorResponse(res) {
		color = "red"
	}
	if IsEventResponse(res) {
		color = "orange"
	}
	return color
}
