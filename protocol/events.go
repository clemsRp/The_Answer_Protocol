package protocol

const (
	CategoryRoom   = "ROOM"
	CategoryGlobal = "GLOBAL"
	CategoryGroup  = "GROUP"
	CategoryStats  = "STATS"
	CategoryCombat = "COMBAT"
)

const (
	TypePresenceEnter = "PRESENCE ENTER"
	TypePresenceLeave = "PRESENCE LEAVE"
	TypeChat          = "CHAT"
	TypeInvite        = "INVITE"
	TypeJoin          = "JOIN"
	TypeLeave         = "LEAVE"
	TypePlayers       = "PLAYERS"
	TypeAllyTurn      = "ALLY_TURN"
)

type ServerEvent struct {
	Category string
	Type     string
	Data     string
}
