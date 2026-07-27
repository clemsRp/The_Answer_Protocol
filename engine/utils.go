package engine

import (
	"slices"
	pr "tap/protocol"
)

func (e *Engine) inform_user(user string, msg string) {
	for ip := range e.clients {
		if e.clients[ip].Datas.Connected && e.clients[ip].Name == user {
			e.clients[ip].Ch <- pr.ServerResponse{Msg: msg}
		}
	}
}

func (e *Engine) inform_room(cli *pr.Client, room string, msg string) {
	for ip := range e.clients {
		if e.clients[ip].Datas.Connected && e.clients[ip].Datas.Room == room && e.clients[ip].Name != cli.Name {
			e.clients[ip].Ch <- pr.ServerResponse{Msg: msg}
		}
	}
}

func (e *Engine) inform_group(cli *pr.Client, group string, msg string) {
	for ip := range e.clients {
		if e.clients[ip].Datas.Connected && e.clients[ip].Datas.Group == group && e.clients[ip].Name != cli.Name {
			e.clients[ip].Ch <- pr.ServerResponse{Msg: msg}
		}
	}
}

func (e *Engine) inform_group_invitations(cli *pr.Client, group string, msg string) {
	for ip := range e.clients {
		if e.clients[ip].Datas.Connected && slices.Contains(e.clients[ip].Datas.Invitation, group) && e.clients[ip].Name != cli.Name {
			e.clients[ip].Ch <- pr.ServerResponse{Msg: msg}
		}
	}
}

func (e *Engine) inform_all(cli *pr.Client, msg string) {
	for ip := range e.clients {
		if e.clients[ip].Datas.Connected && e.clients[ip].Name != cli.Name {
			e.clients[ip].Ch <- pr.ServerResponse{Msg: msg}
		}
	}
}

func (e *Engine) get_nb_connected_players() int {
	res := 0
	for _, cli := range e.clients {
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
