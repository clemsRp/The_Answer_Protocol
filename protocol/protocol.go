package protocol

const (
	CategoryRoom   = "ROOM"
	CategoryGlobal = "GLOBAL"
	CategoryGroup  = "GROUP"
	CategoryStats  = "STATS"
)

const (
	TypePresenceEnter = "PRESENCE ENTER"
	TypePresenceLeave = "PRESENCE LEAVE"
	TypeChat          = "CHAT"
	TypeInvite        = "INVITE"
	TypeJoin          = "JOIN"
	TypeLeave         = "LEAVE"
	TypePlayers       = "PLAYERS"
)

type ServerEvent struct {
	Category string
	Type     string
	Data     string
}

type ServerResponse struct {
	Success bool
	Code    int
	Message string
	Data    string
}
