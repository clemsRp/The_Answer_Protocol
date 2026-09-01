package protocol

const (
	CategoryRoom   = "ROOM"
	CategoryGlobal = "GLOBAL"
	CategoryGroup  = "GROUP"
	CategoryStats  = "STATS"
	CategoryCombat = "COMBAT"
	CategoryItem   = "ITEM"
)

const (
	TypePresenceEnter        = "PRESENCE ENTER"
	TypePresenceLeave        = "PRESENCE LEAVE"
	TypeChat                 = "CHAT"
	TypeInvite               = "INVITE"
	TypeJoin                 = "JOIN"
	TypeKick                 = "KICK"
	TypeLeave                = "LEAVE"
	TypePlayers              = "PLAYERS"
	TypeAllyTurn             = "ALLY_TURN"
	TypeItemDropped          = "ITEM DROPPED"
	TypeItemTook             = "ITEM TOOK"
	TypeGroupPromote         = "GROUP PROMOTE"
	TypeGroupPromoteAccepted = "GROUP PROMOTE ACCEPTED"
	TypeGroupPromoteDeclined = "GROUP PROMOTE DECLINED"
	TypeStats                = "STATS"
)

const (
	PrefixEvtNewLeader = "EVT new_leader="
)

type ServerEvent struct {
	Category string
	Type     string
	Data     string
}
