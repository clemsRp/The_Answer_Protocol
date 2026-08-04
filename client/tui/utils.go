package tui

import (
	"encoding/json"
	"fmt"
	"strings"
	pr "tap/protocol"
)

func findJsonStartIndex(line string) int {
	idx := -1
	for i := 0; i < len(line)-1; i++ {
		if line[i] == ' ' && (line[i+1] == '{' || line[i+1] == '[') {
			idx = i
			break
		}
	}
	return idx
}

func convertServerResponse(line string) pr.ServerResponse {

	line = strings.TrimSpace(line)
	jsonIndex := findJsonStartIndex(line)
	isChatEvent := strings.HasPrefix(line, "EVT ") && strings.Contains(line, " CHAT ")

	if jsonIndex != -1 && !isChatEvent {
		fmt.Println(jsonIndex)
		msg := strings.TrimSpace(line[:jsonIndex])
		datas := line[jsonIndex:]

		var data_json any
		err := json.Unmarshal([]byte(datas), &data_json)
		if err != nil {
			return pr.ServerResponse{Msg: line}
		}
		return pr.ServerResponse{Msg: msg, Datas: data_json}

	} else {
		return pr.ServerResponse{Msg: line}
	}
}
