package protocol

/*
import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ConvertStringToResponse(s string) (*ServerResponse, error) {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "OK") {
		return parseOkResponse(s)
	}

	if strings.HasPrefix(s, "ERR") {
		return parseErrResponse(s)
	}

	return nil, fmt.Errorf("invalid response, neither OK nor ERR: %q", s)
}

func parseOkResponse(s string) (*ServerResponse, error) {
	resp := &ServerResponse{Success: true}

	payload := strings.TrimSpace(strings.TrimPrefix(s, "OK"))

	if len(payload) > 0 {
		if !json.Valid([]byte(payload)) {
			return nil, fmt.Errorf("invalid JSON payload: %q", payload)
		}
		resp.Data = payload
	}

	return resp, nil
}

func parseErrResponse(s string) (*ServerResponse, error) {
	payload := strings.TrimSpace(strings.TrimPrefix(s, "ERR"))
	parts := strings.SplitN(payload, " ", 2)

	if len(parts) < 2 {
		return nil, fmt.Errorf("malformed ERR response, missing code or message: %q", s)
	}

	code, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error code is not a valid integer: %q", parts[0])
	}

	msg := strings.TrimSpace(parts[1])
	if msg == "" {
		return nil, fmt.Errorf("error message is empty")
	}

	resp := &ServerResponse{
		Success: false,
		Code:    code,
		Message: msg,
	}

	return resp, nil
}

func (response *ServerResponse) ConvertToString() string {
	if response.Success {
		if response.Data != "" {
			return fmt.Sprintf("OK %s\n", response.Data)
		}
		return "OK\n"
	}
	return fmt.Sprintf("ERR %d %s\n", response.Code, response.Message)
}
*/
