package utils

import (
	"encoding/json"
	"strings"
	"testing"
)

func AssertResponse(t *testing.T, command, wanted, got string, expectsJSON bool) {
	t.Helper()

	if expectsJSON {
		jsonPart := strings.TrimPrefix(got, wanted+" ")
		jsonPart = strings.TrimSpace(jsonPart)

		if !strings.HasPrefix(got, wanted) {
			t.Error(FormatMismatch(command, wanted, got))
		} else if !json.Valid([]byte(jsonPart)) {
			t.Error(FormatInvalidJSON(jsonPart))
		}
	} else {
		if !strings.HasPrefix(got, wanted) {
			t.Error(FormatMismatch(command, wanted, got))
		}
	}
}
