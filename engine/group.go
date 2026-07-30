package engine

import (
	"errors"
	"fmt"
	"slices"
	pr "tap/protocol"

	"github.com/google/uuid"
)

func (e *Engine) remove_user_in_group(cli *pr.Client, group []*pr.Client) []*pr.Client {
	// Get user index inside group
	cli_index := -1
	for i, user := range group {
		if user.Name == cli.Name {
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

func (e *Engine) create_group(cli *pr.Client) (string, error) {
	// Create group ID
	var group_id string

	for {
		group_id = uuid.New().String()
		if _, exists := e.groups[group_id]; !exists {
			break
		}
	}

	// Check user already in a group
	for _, group := range e.groups {
		for _, user := range group {
			if user.Name == cli.Name {
				return "", errors.New(pr.ErrAlreadyInGroup)
			}
		}
	}

	// Set group
	e.groups[group_id] = []*pr.Client{cli}
	cli.Datas.Group = group_id

	return "OK group=" + group_id, nil
}

func (e *Engine) invite_user_in_group(cli *pr.Client, user_name string) (string, error) {
	// Handle non existant e.groups
	_, ok := e.groups[cli.Datas.Group]
	if !ok {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Check cli is group's leader
	if e.groups[cli.Datas.Group][0].Name != cli.Name {
		return "", errors.New(pr.ErrNoPermission)
	}

	// Handle users already in group
	for _, user := range e.groups[cli.Datas.Group] {
		if user.Name == user_name {
			return "", errors.New(pr.ErrAlreadyInGroup)
		}
	}

	// Get user
	var new_cli *pr.Client
	for ip := range e.clients {
		if e.clients[ip].Name == user_name {
			new_cli = e.clients[ip]

			// Check that invitation isn't already present
			for _, invite := range new_cli.Datas.Invitation {
				if invite == cli.Datas.Group {
					return "", errors.New("ERR Invitation already send")
				}
			}

			// Add invitation to user
			new_cli.Datas.Invitation = append(new_cli.Datas.Invitation, cli.Datas.Group)
			e.clients[ip].Ch <- pr.ServerResponse{Msg: "EVT GROUP INVITE " + cli.Name}

			return "OK", nil
		}
	}

	return "", errors.New(pr.ErrUnknownUser)
}

func (e *Engine) join_group(cli *pr.Client, leader_name string) (string, error) {
	group_name := ""

	// Handle users already in group
	for group_id, group := range e.groups {
		// Find leader group
		if group[0].Name == leader_name {
			group_name = group_id
		}

		for _, user := range group {
			if user.Name == cli.Name {
				return "", errors.New(pr.ErrAlreadyInGroup)
			}
		}
	}

	// Handle non existant e.groups
	if group_name == "" {
		return "", errors.New(pr.ErrInternalServer)
	}
	_, ok := e.groups[group_name]
	if !ok {
		return "", errors.New(pr.ErrInternalServer)
	}

	// Check that user is invited
	invited := false
	for _, invite := range cli.Datas.Invitation {
		if invite == group_name {
			invited = true
			break
		}
	}
	if !invited {
		return "", errors.New(pr.ErrNotInvitedToGroup)
	}

	// Add user in group
	e.groups[group_name] = append(e.groups[group_name], cli)
	cli.Datas.Group = group_name
	e.inform_group(cli, group_name, "EVT GROUP JOIN "+cli.Name)

	// Delete invitation
	invite_index := -1
	for i, invite := range cli.Datas.Invitation {
		if invite == group_name {
			invite_index = i
			break
		}
	}

	if invite_index != -1 {
		cli.Datas.Invitation = append(cli.Datas.Invitation[:invite_index], cli.Datas.Invitation[invite_index+1:]...)
	}

	return "OK group=" + group_name, nil
}

func (e *Engine) leave_group(cli *pr.Client) (string, error) {
	// Check user is inside the group
	if cli.Datas.Group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Delete promotion
	cli.Datas.Promotion = false

	// Remove user from his current group
	groupSlice := e.groups[cli.Datas.Group]
	groupSlice = e.remove_user_in_group(cli, groupSlice)
	e.groups[cli.Datas.Group] = groupSlice

	e.inform_group(cli, cli.Datas.Group, "EVT GROUP LEAVE "+cli.Name)

	// Remove group if needed
	if len(groupSlice) == 0 {
		delete(e.groups, cli.Datas.Group)
	}

	// Re initialize his group value
	cli.Datas.Group = ""

	return "OK", nil
}

func (e *Engine) promote_user(cli *pr.Client, new_leader string) (string, error) {
	// Handle invalid command
	if cli.Datas.Group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	}

	// Find the group
	group_users := e.groups[cli.Datas.Group]

	// Handle invalid rights
	if group_users[0].Name != cli.Name {
		return "", errors.New(pr.ErrNoPermission)

	} else if cli.Name == new_leader {
		return "", errors.New(pr.ErrAlreadyLeader)
	}

	// Find new_leader
	for _, client := range e.clients {
		if client.Name == new_leader {

			// Check new_leader is inside group
			if !slices.Contains(group_users, client) {
				return "", errors.New(pr.ErrNotInGroup)
			}

			// Send promotion
			client.Datas.Promotion = true
			e.inform_user(new_leader, "EVT GROUP PROMOTE "+new_leader)

			return "OK pending_leader=" + cli.Name, nil
		}
	}

	// Handle invalid new_leader
	return "", fmt.Errorf(pr.ErrInternalServer)
}

func (e *Engine) accept_promotion(cli *pr.Client) (string, error) {
	// Check user have a group promotion
	if cli.Datas.Group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	} else if !cli.Datas.Promotion {
		return "", errors.New(pr.ErrNotPromoted)
	}

	// Update promotion and leadership
	cli.Datas.Promotion = false
	new_leader_index := GetElementIndex(e.groups[cli.Datas.Group], cli)
	MoveElement(e.groups[cli.Datas.Group], new_leader_index, 0)

	e.inform_group(cli, cli.Datas.Group, "EVT GROUP PROMOTE ACCEPTED "+cli.Name)
	e.inform_group_invitations(cli, cli.Datas.Group, "EVT GROUP PROMOTE ACCEPTED "+cli.Name)

	return "OK new_leader=" + cli.Name, nil
}

func (e *Engine) decline_promotion(cli *pr.Client) (string, error) {
	// Check user have a group promotion
	if cli.Datas.Group == "" {
		return "", errors.New(pr.ErrNotInGroup)
	} else if !cli.Datas.Promotion {
		return "", errors.New(pr.ErrNotPromoted)
	}

	// Update promotion
	cli.Datas.Promotion = false

	e.inform_group(cli, cli.Datas.Group, "EVT GROUP PROMOTE DECLINED "+cli.Name)
	e.inform_group_invitations(cli, cli.Datas.Group, "EVT GROUP PROMOTE DECLINED "+cli.Name)

	return "OK", nil
}
