package main

func inform_room(clients map[string]*Client, cli *Client, room string, msg string) {
	for ip := range clients {
		if clients[ip].datas.room == room && clients[ip].name != cli.name {
			clients[ip].ch <- Response{msg, Datas{}, Request{}}
		}
	}
}
