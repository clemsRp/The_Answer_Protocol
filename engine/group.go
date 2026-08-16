package engine

import (
	"errors"
	"fmt"
	"slices"
	pr "tap/protocol"

	"github.com/google/uuid"
)

func (e *Engine) remove_user_in_group(player *Player, group *Group) *Group {
	// Get user index inside group
	cli_index := -1
	for i, user := range group.players {
		if user.name == player.name {
			cli_index = i
			break
		}
	}

	// Remove user
	if cli_index != -1 {
		group.players = append(group.players[:cli_index], group.players[cli_index+1:]...)
	}

	return group
}

func (e *Engine) create_group(player *Player) (string, error) {
	// Create group ID
	group_id := uuid.New().String()

	// Check user already in a group
	if player.group != "" {
		return "", errors.New(pr.ErrAlreadyInGroup)
	}

	// Set group
	e.groups[group_id] = &Group{
		id:      group_id,
		leader:  player,
		players: []*Player{player},
	}
	player.group = group_id

	return "OK group=" + group_id, nil
}

func (e *Engine) invite_user_in_group(player *Player, user_name string) (string, error) {
	// Handle non existant groups
	group, ok := e.groups[player.group]
	if !ok {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Check player is group's leader
	if group.leader.name != player.name {
		return "", errors.New(pr.ErrNoPermission)
	}

	// Handle users already in group
	for _, user := range group.players {
		if user.name == user_name {
			return "", errors.New(pr.ErrAlreadyInGroup)
		}
	}

	// Get target user from e.players
	targetPlayer, exists := e.players[user_name]
	if !exists {
		return "", errors.New(pr.ErrUnknownUser)
	}

	// Check that invitations isn't already present
	for _, invite := range targetPlayer.invitations {
		if invite == player.group {
			return "", errors.New("ERR Invitation already send")
		}
	}

	// Add invitations to user
	targetPlayer.invitations = append(targetPlayer.invitations, player.group)

	e.inform_user(targetPlayer, "EVT GROUP INVITE "+player.name)

	return "OK", nil
}

func (e *Engine) kick_user_in_group(player *Player, user_name string) (string, error) {
	// Handle non existant groups
	_, ok := e.groups[player.group]
	if !ok {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Check player is group's leader
	if e.groups[player.group].players[0].name != player.name {
		return "", errors.New(pr.ErrNoPermission)
	}

	// Handle users not in group
	for _, user := range e.groups[player.group].players {
		if user.name == user_name {
			// Get target user from e.players
			targetPlayer, exists := e.players[user_name]
			if !exists {
				return "", errors.New(pr.ErrUnknownUser)
			}

			e.inform_group(player, player.group, "EVT GROUP KICK "+targetPlayer.name)

			// Remove user from his current group
			groupSlice := e.groups[player.group]
			groupSlice = e.remove_user_in_group(targetPlayer, groupSlice)
			e.groups[player.group] = groupSlice

			targetPlayer.group = ""

			return "OK", nil
		}
	}

	return "", errors.New(pr.ErrNotInGroup)
}

func (e *Engine) join_group(player *Player, leader_name string) (string, error) {
	group_name := ""

	// Handle users already in group
	for group_id, group := range e.groups {
		// Find leader group
		if group.leader.name == leader_name {
			group_name = group_id
		}

		for _, user := range group.players {
			if user.name == player.name {
				return "", errors.New(pr.ErrAlreadyInGroup)
			}
		}
	}

	// Handle non existant groups
	if group_name == "" {
		return "", errors.New(pr.ErrInternalServer)
	}
	group, ok := e.groups[group_name]
	if !ok {
		return "", errors.New(pr.ErrInternalServer)
	}

	// Check that user is invited
	invited := false
	for _, invite := range player.invitations {
		if invite == group_name {
			invited = true
			break
		}
	}
	if !invited {
		return "", errors.New(pr.ErrNotInvitedToGroup)
	}

	// Add user in group
	group.players = append(group.players, player)
	player.group = group_name
	e.inform_group(player, group.id, "EVT GROUP JOIN "+player.name)

	// Delete invitations
	invite_index := -1
	for i, invite := range player.invitations {
		if invite == group_name {
			invite_index = i
			break
		}
	}

	if invite_index != -1 {
		player.invitations = append(player.invitations[:invite_index], player.invitations[invite_index+1:]...)
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

	// Check leader
	was_leader := (e.groups[player.group].players[0].name == player.name)

	// Remove user from his current group
	group := e.groups[player.group]
	group = e.remove_user_in_group(player, group)
	e.groups[player.group] = group

	e.inform_group(player, group.id, "EVT GROUP LEAVE "+player.name)

	// Remove group if needed
	if len(group.players) == 0 {
		delete(e.groups, player.group)

	} else if was_leader {
		e.inform_group(player, player.group, "EVT new_leader="+e.groups[player.group].players[0].name)
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
	group := e.groups[player.group]

	// Handle invalid rights
	if group.leader.name != player.name {
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
	if !slices.Contains(group.players, targetPlayer) {
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

	oldLeader := e.groups[player.group].players[0]
	oldLeader.send_promotion = false

	// Update promotion and leadership
	player.promotion = false
	new_leader_index := GetElementIndex(e.groups[player.group].players, player)

	groupSlice := e.groups[player.group].players
	MoveElement(groupSlice, new_leader_index, 0)
	e.groups[player.group].players = groupSlice

	e.groups[player.group].leader = player
	e.inform_group(player, e.groups[player.group].id, "EVT GROUP PROMOTE ACCEPTED "+player.name)
	e.inform_group_invitations(player, e.groups[player.group].id, "EVT GROUP PROMOTE ACCEPTED "+player.name)

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
	_, exists := e.groups[player.group]
	if !exists {
		return "", errors.New(pr.ErrInternalServer)
	}

	// Change leader promotion status
	e.groups[player.group].players[0].send_promotion = false

	e.inform_group(player, player.group, "EVT GROUP PROMOTE DECLINED "+player.name)
	e.inform_group_invitations(player, player.group, "EVT GROUP PROMOTE DECLINED "+player.name)

	return "OK", nil
}

func (e *Engine) getPlayerGroup(player *Player) (*Group, bool) {
	currentGroup, isInGroup := e.groups[player.group]
	if !isInGroup {
		return nil, false
	}
	return currentGroup, true
}
