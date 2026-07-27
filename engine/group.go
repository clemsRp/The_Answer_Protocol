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
	group_id := uuid.New().String()

	// Handle existing group
	_, ok := e.groups[group_id]
	if ok {
		return "", errors.New("ERR Group already exist")
	}

	// Check user already in a group
	for _, group := range e.groups {
		for _, user := range group {
			if user.Name == cli.Name {
				return "", errors.New("ERR user already inside a group")
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
		return "", errors.New("ERR 401 NOT_IN_GROUP")
	}

	// Check cli is group's leader
	if e.groups[cli.Datas.Group][0].Name != cli.Name {
		return "", errors.New("ERR User isn't group's leader")
	}

	// Handle users already in group
	for _, user := range e.groups[cli.Datas.Group] {
		if user.Name == user_name {
			return "", errors.New("ERR 402 ALREADY_IN_GROUP")
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

	return "", errors.New("ERR new user not find")
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
				return "", errors.New("ERR player already in group")
			}
		}
	}

	// Handle non existant e.groups
	if group_name == "" {
		return "", errors.New("ERR Invalid leader name")
	}
	_, ok := e.groups[group_name]
	if !ok {
		return "", errors.New("ERR Group doesn't exist yet")
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
		return "", errors.New("ERR User isn't invited by this group")
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
		return "", errors.New("ERR 401 NOT_IN_GROUP")
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
		return "", errors.New("ERR 401 NOT_IN_GROUP")
	}

	// Find the group
	group_users := e.groups[cli.Datas.Group]

	// Handle invalid rights
	if group_users[0].Name != cli.Name {
		return "", errors.New("ERR You are not the leader of the group")

	} else if cli.Name == new_leader {
		return "", errors.New("ERR You are already the leader of the group")
	}

	// Find new_leader
	for _, client := range e.clients {
		if client.Name == new_leader {

			// Check new_leader is inside group
			if !slices.Contains(group_users, client) {
				return "", errors.New("ERR 401 NOT_IN_GROUP")
			}

			// Send promotion
			client.Datas.Promotion = true
			e.inform_user(new_leader, "EVT GROUP PROMOTE "+new_leader)

			return "OK pending_leader=" + cli.Name, nil
		}
	}

	// Handle invalid new_leader
	return "", fmt.Errorf("ERR 404 USER_NOT_FOUND")
}

func (e *Engine) accept_promotion(cli *pr.Client) (string, error) {
	// Check user have a group promotion
	if cli.Datas.Group == "" {
		return "", errors.New("ERR 401 NOT_IN_GROUP")
	} else if !cli.Datas.Promotion {
		return "", errors.New("ERR You don't have promotion")
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
		return "", errors.New("ERR 401 NOT_IN_GROUP")
	} else if !cli.Datas.Promotion {
		return "", errors.New("ERR You don't have promotion")
	}

	// Update promotion
	cli.Datas.Promotion = false

	e.inform_group(cli, cli.Datas.Group, "EVT GROUP PROMOTE DECLINED "+cli.Name)
	e.inform_group_invitations(cli, cli.Datas.Group, "EVT GROUP PROMOTE DECLINED "+cli.Name)

	return "OK", nil
}
