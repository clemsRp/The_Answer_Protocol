package main

import (
	"errors"
)

func remove_user_in_group(cli *Client, group []*Client) []*Client {
	// Get user index inside group
	cli_index := -1
	for i, user := range group {
		if user.name == cli.name {
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

func create_group(cli *Client, group_name string) (string, error) {
	// Handle existing group
	_, ok := groups[group_name]
	if ok {
		return "", errors.New("ERR Group already exist")
	}

	// Check user already in a group
	for _, group := range groups {
		for _, user := range group {
			if user.name == cli.name {
				return "", errors.New("ERR user already inside a group")
			}
		}
	}

	// Set group
	groups[group_name] = []*Client{cli}
	cli.datas.group = group_name

	return "OK group=" + group_name, nil
}

func invite_user_in_group(clients map[string]*Client, cli *Client, user_name string) (string, error) {
	// Handle non existant groups
	_, ok := groups[cli.datas.group]
	if !ok {
		return "", errors.New("ERR Group doesn't exist yet")
	}

	// Check cli is group's leader
	if groups[cli.datas.group][0].name != cli.name {
		return "", errors.New("ERR User isn't group's leader")
	}

	// Handle users already in group
	for _, user := range groups[cli.datas.group] {
		if user.name == user_name {
			return "", errors.New("ERR player already in group")
		}
	}

	// Get user
	var new_cli *Client
	for ip := range clients {
		if clients[ip].name == user_name {
			new_cli = clients[ip]

			// Check that invitation isn't already present
			for _, invite := range new_cli.datas.invitation {
				if invite == cli.datas.group {
					return "", errors.New("ERR Invitation already send")
				}
			}

			// Add invitation to user
			new_cli.datas.invitation = append(new_cli.datas.invitation, cli.datas.group)

			return "OK", nil
		}
	}

	return "", errors.New("ERR new user not find")
}

func join_group(cli *Client, group_name string) (string, error) {
	// Handle non existant groups
	_, ok := groups[group_name]
	if !ok {
		return "", errors.New("ERR Group doesn't exist yet")
	}

	// Handle users already in group
	for _, group := range groups {
		for _, user := range group {
			if user.name == cli.name {
				return "", errors.New("ERR player already in group")
			}
		}
	}

	// Check that user is invited
	invited := false
	for _, invite := range cli.datas.invitation {
		if invite == group_name {
			invited = true
			break
		}
	}
	if !invited {
		return "", errors.New("ERR User isn't invited by this group")
	}

	// Add user in group
	groups[group_name] = append(groups[group_name], cli)
	cli.datas.group = group_name

	// Delete invitation
	invite_index := -1
	for i, invite := range cli.datas.invitation {
		if invite == group_name {
			invite_index = i
			break
		}
	}

	if invite_index != -1 {
		cli.datas.invitation = append(cli.datas.invitation[:invite_index], cli.datas.invitation[invite_index+1:]...)
	}

	return "OK group=" + group_name, nil
}

func leave_group(cli *Client) (string, error) {
	// Check user is inside the group
	if cli.datas.group == "" {
		return "", errors.New("ERR User isn't inside a group")
	}

	// Remove user from his current group
	groupSlice := groups[cli.datas.group]
	groupSlice = remove_user_in_group(cli, groupSlice)
	groups[cli.datas.group] = groupSlice

	// Remove group if needed
	if len(groupSlice) == 0 {
		delete(groups, cli.datas.group)
	}

	// Re initialize his group value
	cli.datas.group = ""

	return "OK", nil
}
