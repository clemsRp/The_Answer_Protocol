package protocol

// func IsValidResponse(s string) error {
// 	s = strings.TrimSpace(s)

// 	if strings.HasPrefix(s, "OK") {
// 		return checkOkResponse(s)
// 	}

// 	if strings.HasPrefix(s, "ERR") {
// 		return checkErrResponse(s)
// 	}

// 	return fmt.Errorf("invalid response, neither OK nor ERR: %q", s)
// }

// func checkOkResponse(s string) error {

// 	payload := strings.TrimSpace(strings.TrimPrefix(s, "OK"))

// 	if len(payload) > 0 {
// 		if !json.Valid([]byte(payload)) {
// 			return fmt.Errorf("invalid JSON payload: %q", payload)
// 		}
// 	}

// 	return nil
// }

// func checkErrResponse(s string) error {
// 	payload := strings.TrimSpace(strings.TrimPrefix(s, "ERR"))
// 	parts := strings.SplitN(payload, " ", 2)

// 	if len(parts) < 2 {
// 		return fmt.Errorf("malformed ERR response, missing code or message: %q", s)
// 	}

// 	_, err := strconv.Atoi(parts[0])
// 	if err != nil {
// 		return fmt.Errorf("error code is not a valid integer: %q", parts[0])
// 	}

// 	msg := strings.TrimSpace(parts[1])
// 	if msg == "" {
// 		return fmt.Errorf("error message is empty")
// 	}

// 	return nil
// }

// func (response *ServerResponse) ConvertToString() string {
// 	if response.Success {
// 		if response.Data != "" {
// 			return fmt.Sprintf("OK %s\n", response.Data)
// 		}
// 		return "OK\n"
// 	}
// 	return fmt.Sprintf("ERR %d %s\n", response.Code, response.Message)
// }
