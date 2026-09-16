package variables

import "time"

type Chat struct {
	Msg        string
	Pseudo     string
	EmoteIndex int
	Time       time.Time
}

type ChatPanel struct {
	PanelOpen    bool
	CurrentScope string
	Msg          string
	LastMsgText  string
	LastMsgScope string
	ScopeChats   map[string][]Chat
}
