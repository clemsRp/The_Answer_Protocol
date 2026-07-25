package main

import (
	"slices"
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

func inform_group_invitations(clients map[string]*pr.Client, cli *pr.Client, group string, msg string) {
	for ip := range clients {
		if clients[ip].Datas.Connected && slices.Contains(clients[ip].Datas.Invitation, group) && clients[ip].Name != cli.Name {
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

func GetElementIndex[T comparable](slice []T, element_need T) int {
	for index, element_possible := range slice {
		if element_possible == element_need {
			return index
		}
	}

	return -1
}

func MoveElement[T any](slice []T, from int, to int) []T {
	n := len(slice)

	// Check slice limits
	if from < 0 || from >= n || to < 0 || to >= n || from == to {
		return slice
	}

	elem := slice[from]

	// Delete element from his initial position
	slice = append(slice[:from], slice[from+1:]...)

	// Add element to his final position
	slice = append(slice[:to], append([]T{elem}, slice[to:]...)...)

	return slice
}
