package protocol

const (
	CmdConnect   = "CONNECT"
	CmdLook      = "LOOK"
	CmdMove      = "MOVE"
	CmdChat      = "CHAT"
	CmdTake      = "TAKE"
	CmdDrop      = "DROP"
	CmdInventory = "INVENTORY"
	CmdTalk      = "TALK"
	CmdAttack    = "ATTACK"
	CmdStatus    = "STATUS"
	CmdQuest     = "QUEST"
	CmdQuests    = "QUESTS"
	CmdWho       = "WHO"
	CmdUnGrouped = "UNGROUPED"
	CmdGrouped   = "GROUPED"
	CmdGroup     = "GROUP"
	CmdQuit      = "QUIT"

	GlobalChat = "GLOBAL"
	RoomChat   = "ROOM"
	GroupChat  = "GROUP"

	CreateGroup         = "CREATE"
	InviteGroup         = "INVITE"
	KickGroup           = "KICK"
	JoinGroup           = "JOIN"
	LeaveGroup          = "LEAVE"
	PromoteGroup        = "PROMOTE"
	AcceptPromoteGroup  = "ACCEPT"
	DeclinePromoteGroup = "DECLINE"

	South = "south"
	North = "north"
	East  = "east"
	West  = "west"
)
