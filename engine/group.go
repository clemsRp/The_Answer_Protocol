package engine

import (
	"errors"
	"fmt"
	"slices"
	pr "tap/protocol"

	"github.com/google/uuid"
)

func (e *Engine) remove_user_in_group(player *Player, group []*Player) []*Player {
	// Get user index inside group
	cli_index := -1
	for i, user := range group {
		if user.name == player.name {
			cli_index = i
			break
		}
	}

	// Remove user
	if cli_index != -1 {
		group = append(group[:cli_index], group[cli_index+1:]...)
	}

	return group
}

func (e *Engine) create_group(player *Player) (string, error) {
	// Create group ID
	var group_id string

	for {
		group_id = uuid.New().String()
		if _, exists := e.groups[group_id]; !exists {
			break
		}
	}

	// Check user already in a group
	if player.group != "" {
		return "", errors.New(pr.ErrAlreadyInGroup)
	}

	// Set group
	e.groups[group_id] = []*Player{player}
	player.group = group_id

	return "OK group=" + group_id, nil
}

func (e *Engine) invite_user_in_group(player *Player, user_name string) (string, error) {
	// Handle non existant groups
	_, ok := e.groups[player.group]
	if !ok {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Check player is group's leader
	if e.groups[player.group][0].name != player.name {
		return "", errors.New(pr.ErrNoPermission)
	}

	// Handle users already in group
	for _, user := range e.groups[player.group] {
		if user.name == user_name {
			return "", errors.New(pr.ErrAlreadyInGroup)
		}
	}

	// Get target user from e.players (au lieu de e.clients)
	targetPlayer, exists := e.players[user_name]
	if !exists {
		return "", errors.New(pr.ErrUnknownUser)
	}

	// Check that invitation isn't already present
	for _, invite := range targetPlayer.invitation {
		if invite == player.group {
			return "", errors.New("ERR Invitation already send")
		}
	}

	// Add invitation to user
	targetPlayer.invitation = append(targetPlayer.invitation, player.group)

	// Utilisation de inform_user pour notifier proprement via ServerOutput
	e.inform_user(targetPlayer, "EVT GROUP INVITE "+player.name)

	return "OK", nil
}

func (e *Engine) join_group(player *Player, leader_name string) (string, error) {
	group_name := ""

	// Handle users already in group
	for group_id, group := range e.groups {
		// Find leader group
		if group[0].name == leader_name {
			group_name = group_id
		}

		for _, user := range group {
			if user.name == player.name {
				return "", errors.New(pr.ErrAlreadyInGroup)
			}
		}
	}

	// Handle non existant groups
	if group_name == "" {
		return "", errors.New(pr.ErrInternalServer)
	}
	_, ok := e.groups[group_name]
	if !ok {
		return "", errors.New(pr.ErrInternalServer)
	}

	// Check that user is invited
	invited := false
	for _, invite := range player.invitation {
		if invite == group_name {
			invited = true
			break
		}
	}
	if !invited {
		return "", errors.New(pr.ErrNotInvitedToGroup)
	}

	// Add user in group
	e.groups[group_name] = append(e.groups[group_name], player)
	player.group = group_name
	e.inform_group(player, group_name, "EVT GROUP JOIN "+player.name)

	// Delete invitation
	invite_index := -1
	for i, invite := range player.invitation {
		if invite == group_name {
			invite_index = i
			break
		}
	}

	if invite_index != -1 {
		player.invitation = append(player.invitation[:invite_index], player.invitation[invite_index+1:]...)
	}

	return "OK group=" + group_name, nil
}

func (e *Engine) leave_group(player *Player) (string, error) {
	// Check user is inside the group
	if player.group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Delete promotion
	player.promotion = false

	// Remove user from his current group
	groupSlice := e.groups[player.group]
	groupSlice = e.remove_user_in_group(player, groupSlice)
	e.groups[player.group] = groupSlice

	e.inform_group(player, player.group, "EVT GROUP LEAVE "+player.name)
	e.inform_group(player, player.group, "EVT new_leader="+e.groups[player.group][0].name)

	// Remove group if needed
	if len(groupSlice) == 0 {
		delete(e.groups, player.group)
	}

	// Re initialize his group value
	player.group = ""

	return "OK", nil
}

func (e *Engine) promote_user(player *Player, new_leader string) (string, error) {
	// Handle invalid command
	if player.group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Find the group
	group_users := e.groups[player.group]

	// Handle invalid rights
	if group_users[0].name != player.name {
		return "", errors.New(pr.ErrNoPermission)

	} else if player.name == new_leader {
		return "", errors.New(pr.ErrAlreadyLeader)

	} else if player.send_promotion {
		return "", errors.New(pr.ErrNoPermission)
	}

	// Find new_leader in e.players
	targetPlayer, exists := e.players[new_leader]
	if !exists {
		return "", fmt.Errorf(pr.ErrInternalServer)
	}

	// Check new_leader is inside group
	if !slices.Contains(group_users, targetPlayer) {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Send promotion
	player.send_promotion = true
	targetPlayer.promotion = true
	e.inform_user(targetPlayer, "EVT GROUP PROMOTE "+new_leader)

	return "OK pending_leader=" + player.name, nil
}

func (e *Engine) accept_promotion(player *Player) (string, error) {
	// Check user have a group promotion
	if player.group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	} else if !player.promotion {
		return "", errors.New(pr.ErrNotPromoted)
	}

	// Update promotion and leadership
	player.promotion = false
	new_leader_index := GetElementIndex(e.groups[player.group], player)
	MoveElement(e.groups[player.group], new_leader_index, 0)

	e.inform_group(player, player.group, "EVT GROUP PROMOTE ACCEPTED "+player.name)
	e.inform_group_invitations(player, player.group, "EVT GROUP PROMOTE ACCEPTED "+player.name)

	return "OK new_leader=" + player.name, nil
}

func (e *Engine) decline_promotion(player *Player) (string, error) {
	// Check user have a group promotion
	if player.group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	} else if !player.promotion {
		return "", errors.New(pr.ErrNotPromoted)
	}

	// Update promotion
	player.promotion = false

	e.inform_group(player, player.group, "EVT GROUP PROMOTE DECLINED "+player.name)
	e.inform_group_invitations(player, player.group, "EVT GROUP PROMOTE DECLINED "+player.name)

	return "OK", nil
}
