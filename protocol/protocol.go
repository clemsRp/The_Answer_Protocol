package protocol

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



type Exchanger struct {
	ServerInput  chan ServerRequest
	ServerOutput chan EngineResponse
	JoinChan     chan string
	LeaveChan    chan string
}

