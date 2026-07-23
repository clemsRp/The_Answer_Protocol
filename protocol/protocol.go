package protocol

import (
	"net"
	"time"
)

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

// type ServerResponse struct {
// 	Success bool
// 	Code    int
// 	Message string
// 	Data    string
// }

type Datas struct {
	Room          string
	Status        string
	Inventory     []string
	Invitation    []string
	Group         string
	Hp            int
	Max_hp        int
	Connected     bool
	Last_cmd_time time.Time
	Spam_warning  int
}

type Client struct {
	Conn  net.Conn      `json:"-"`
	Ch    chan Response `json:"-"`
	Ip    string
	Name  string
	Datas Datas
}

type Request struct {
	Cli *Client `json:"-"`
	Msg string
}

type Response struct {
	Msg   string
	Datas any
	Req   Request
}

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
	CmdGroup     = "GROUP"
	CmdQuit      = "QUIT"

	GlobalChat = "GLOBAL"
	RoomChat   = "ROOM"
	GroupChat  = "GROUP"

	CreateGroup = "CREATE"
	InviteGroup = "INVITE"
	JoinGroup   = "JOIN"
	LeaveGroup  = "LEAVE"

	South = "south"
	North = "north"
	East  = "east"
	West  = "west"
)
