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

const (
	// Combat Events
	EventCombatUpdate                    = "EVT COMBAT UPDATE"
	EventCombatTurn                      = "EVT COMBAT TURN player="
	EventCombatVictory                   = "EVT COMBAT VICTORY"
	EventCombatDefeat                    = "EVT COMBAT DEFEAT new_room="
	EventCombatStarted                   = "EVT COMBAT STARTED launcher="
	EventCombatAllyLeaveCombat           = "EVT COMBAT ALLY_LEAVE_COMBAT player="
	EventDistantGroupCombatStartedCombat = "EVT GROUP DISTANT_COMBAT_START launcher="

	// Room Events
	EventRoomPresenceEnter = "EVT ROOM PRESENCE ENTER"
	EventRoomPresenceLeave = "EVT ROOM PRESENCE LEAVE"

	// Group Events
	EventGroupInvite = "EVT GROUP INVITE"
	EventGroupJoin   = "EVT GROUP JOIN"
	EventGroupLeave  = "EVT GROUP LEAVE"

	// Chat Events
	EventGlobalChat = "EVT GLOBAL CHAT"
	EventRoomChat   = "EVT ROOM CHAT"
	EventGroupChat  = "EVT GROUP CHAT"

	// Stats Events
	EventStatsPlayers = "EVT STATS players="
)

type ServerEvent struct {
	Category string
	Type     string
	Data     string
}
