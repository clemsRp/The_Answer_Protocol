package engine

import (
	"slices"
	pr "tap/protocol"
)

func (e *Engine) inform_user(player *Player, msg string) {
	e.exchanger.ServerOutput <- pr.EngineResponse{Id: player.id, Msg: msg}
}

func (e *Engine) inform_room(player *Player, room *Room, msg string) {
	for pseudo, p := range e.players {
		if p.room == room && pseudo != player.name {
			e.exchanger.ServerOutput <- pr.EngineResponse{Id: p.id, Msg: msg}
		}
	}
}

func (e *Engine) inform_group(player *Player, group string, msg string) {
	for pseudo, p := range e.players {
		if p.group == group && pseudo != player.name {
			e.exchanger.ServerOutput <- pr.EngineResponse{Id: p.id, Msg: msg}
		}
	}
}

func (e *Engine) inform_group_invitations(player *Player, group string, msg string) {
	for pseudo, p := range e.players {
		if slices.Contains(player.invitations, group) && pseudo != player.name {
			e.exchanger.ServerOutput <- pr.EngineResponse{Id: p.id, Msg: msg}
		}
	}
}

func (e *Engine) inform_all(player *Player, msg string) {
	for pseudo, p := range e.players {
		if pseudo != player.name {
			e.exchanger.ServerOutput <- pr.EngineResponse{Id: p.id, Msg: msg}
		}
	}
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
