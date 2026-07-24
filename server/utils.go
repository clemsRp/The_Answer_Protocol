package main

import (
	pr "tap/protocol"
)

func inform_user(clients map[string]*pr.Client, user string, msg string) {
	for ip := range clients {
		if clients[ip].Datas.Connected && clients[ip].Name == user {
			clients[ip].Ch <- pr.Response{Msg: msg}
		}
	}
}

func inform_room(clients map[string]*pr.Client, cli *pr.Client, room string, msg string) {
	for ip := range clients {
		if clients[ip].Datas.Connected && clients[ip].Datas.Room == room && clients[ip].Name != cli.Name {
			clients[ip].Ch <- pr.Response{Msg: msg}
		}
	}
}

func inform_group(clients map[string]*pr.Client, cli *pr.Client, group string, msg string) {
	for ip := range clients {
		if clients[ip].Datas.Connected && clients[ip].Datas.Group == group && clients[ip].Name != cli.Name {
			clients[ip].Ch <- pr.Response{Msg: msg}
		}
	}
}

func inform_all(clients map[string]*pr.Client, cli *pr.Client, msg string) {
	for ip := range clients {
		if clients[ip].Datas.Connected && clients[ip].Name != cli.Name {
			clients[ip].Ch <- pr.Response{Msg: msg}
		}
	}
}

func get_nb_connected_players(clients map[string]*pr.Client) int {
	res := 0
	for _, cli := range clients {
		if cli.Datas.Connected {
			res++
		}
	}

	return res
}
