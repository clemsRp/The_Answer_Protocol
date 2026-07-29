package protocol

// Session creation, enter in game errors
const (
	ErrNameInUse    = "ERR 201 NAME_IN_USE"
	ErrNotConnected = "ERR 200 NOT_CONNECTED"
	ErrInvalidName  = "ERR 202 INVALID_NAME"
)

// Movement errors
const (
	ErrNoExit = "ERR 301 NO_EXIT"
)

// 400 codes are GAME LOGIC ERRORS
// Group errors 401, 402
const (
	ErrInvalidCommand  = "ERR 400 INVALID_COMMAND"
	ErrNotInGroup      = "ERR 401 NOT_IN_GROUP"
	ErrAlreadyInGroup  = "ERR 402 ALREADY_IN_GROUP"
	ErrUserDoesntExist = "ERR 403 UNKNOWN_USER"
)

// Items & Inventory errors
const (
	ErrItemNotFound       = "ERR 404 ITEM_NOT_FOUND"
	ErrItemNotInInventory = "ERR 404 ITEM_NOT_IN_INVENTORY"
)

// NPCs & Quests
const (
	ErrNpcNotFound      = "ERR 404 NPC_NOT_FOUND"
	ErrNpcNotHostile    = "ERR 405 NPC_NOT_HOSTILE"
	ErrNoQuestAvailable = "ERR 406 NO_QUEST_AVAILABLE"
)

// Internal server
const (
	ErrInternalServer = "ERR 500 INTERNAL_SERVER_ERROR"
)

// Authentication & Connection errors
const (
	ErrConnectionFailed = "ERR 900 CONNECTION_FAILED"
	ErrSendFailed       = "ERR 901 SEND_FAILED"
)

// MORE errors (NOT IN RFC protocol) 500+ codes
