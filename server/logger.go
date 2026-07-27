package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (s *Server) writeLog(level string, message string, datas map[string]any) {
	entry := Log{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Datas:     datas,
	}

	jsonBytes, err := json.MarshalIndent(entry, "", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Critical error from logger: %v\n", err)
		return
	}

	if level == "ERROR" {
		fmt.Fprintln(os.Stderr, string(jsonBytes))
	} else {
		fmt.Fprintln(os.Stdout, string(jsonBytes))
	}
}

func (s *Server) LogInfo(msg string, datas map[string]any) {
	s.logs <- Log{Level: "INFO", Message: msg, Datas: datas}
}
func (s *Server) LogWarn(msg string, datas map[string]any) {
	s.logs <- Log{Level: "WARN", Message: msg, Datas: datas}
}
func (s *Server) LogError(msg string, datas map[string]any) {
	s.logs <- Log{Level: "ERROR", Message: msg, Datas: datas}
}
