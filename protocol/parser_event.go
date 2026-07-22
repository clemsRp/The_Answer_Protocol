package protocol

import "strings"

func ConvertStringToEvent(s string) (*ServerEvent, error) {
	s = strings.TrimSpace(s)

}
