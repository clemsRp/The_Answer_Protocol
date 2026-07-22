package main

import (
	pr "tap/protocol"
)

func inform_room(clients map[string]*pr.Client, cli *pr.Client, room string, msg string) {
	for ip := range clients {
		if clients[ip].Datas.Connected && clients[ip].Datas.Room == room && clients[ip].Name != cli.Name {
			clients[ip].Ch <- pr.Response{Msg: msg}
		}
	}
}
