package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func writeLog(level string, message string, datas map[string]any) {
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

func LogInfo(msg string, datas map[string]any) {
	logs <- Log{Level: "INFO", Message: msg, Datas: datas}
}
func LogWarn(msg string, datas map[string]any) {
	logs <- Log{Level: "WARN", Message: msg, Datas: datas}
}
func LogError(msg string, datas map[string]any) {
	logs <- Log{Level: "ERROR", Message: msg, Datas: datas}
}
