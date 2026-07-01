package main

func is_inside(elements []string, value string) bool {
	for _, element := range elements {
		if element == value {
			return true
		}
	}

	return false
}

func inform_room(clients map[string]*Client, cli *Client, room string, msg string) {
	for ip := range clients {
		if clients[ip].datas.room == room && clients[ip].name != cli.name {
			clients[ip].ch <- Response{msg, Datas{}, Request{}}
		}
	}
}
