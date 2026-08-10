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

// type Datas struct {
// 	Room          string
// 	Status        string
// 	Inventory     []string
// 	Invitation    []string
// 	Group         string
// 	Promotion     bool
// 	Hp            int
// 	Max_hp        int
// 	Connected     bool
// 	Last_cmd_time time.Time
// 	Spam_warning  int
// 	Quests        []*TrackedQuestData
// }

// server to engine
type ServerRequest struct {
	Ip  string
	Msg string
}

// engine to server
type EngineResponse struct {
	Ip    string
	Msg   string
	Datas any
	Err   error
}

// server to client
type ServerResponse struct {
	Msg   string `json:"msg"`
	Datas any    `json:"datas,omitempty"`
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

	CreateGroup         = "CREATE"
	InviteGroup         = "INVITE"
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

type Exchanger struct {
	ServerInput  chan ServerRequest
	ServerOutput chan EngineResponse
	JoinChan     chan string
	LeaveChan    chan string
}

type PlayerInfo struct {
	Name      string   `json:"name"`
	Location  string   `json:"location"`
	Hp        int      `json:"hp"`
	HpMax     int      `json:"hp_max"`
	Room      string   `json:"room"`
	Group     string   `json:"group"`
	Inventory []string `json:"inventory"`
}
