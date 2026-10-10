package cmd

import (
	"errors"
	"net"
	"os"
	"slices"
	"strconv"
)

func GetConnection(allowIP bool) (string, string, error) {
	// Default values
	ip := "localhost"
	port := "8080"

	var err error
	argv := os.Args
	argc := len(argv)

	if allowIP {
		switch argc {
		case 1:
			return ip, port, nil
		case 3, 5:
			num := (argc - 1) / 2
			err = ChangeConnectionDatas(argv, &ip, &port, num, allowIP)
			if err != nil {
				return "", "", err
			}
			return ip, port, nil
		default:
			return "", "", errors.New("Invalid format: [--ip <ip_address>] [--port <port>]")
		}
	} else {
		switch argc {
		case 1:
			return ip, port, nil
		case 3:
			err = ChangeConnectionDatas(argv, &ip, &port, 1, allowIP)
			if err != nil {
				return "", "", err
			}
			return ip, port, nil
		default:
			return "", "", errors.New("Invalid format: [--port <port>]")
		}
	}
}

func ChangeConnectionDatas(argv []string, ip, port *string, num int, allowIP bool) error {
	updated := make([]string, 0)
	start := 1

	for index := 0; index < num; index++ {
		key := argv[start+2*index]
		value := argv[start+2*index+1]

		switch key {
		case "--ip":
			if !allowIP {
				return errors.New("Invalid parameter: '--ip' is not allowed")
			}
			// Check ip isn't already defined
			if slices.Contains(updated, "--ip") {
				return errors.New("Invalid parameter: '--ip' already redefined")
			}

			// Check ip is valid
			parsedIP := net.ParseIP(value)
			if parsedIP == nil || parsedIP.To4() == nil {
				return errors.New("Invalid parameter: invalid IP format")
			}
			*ip = value

		case "--port":
			// Check port isn't already defined
			if slices.Contains(updated, "--port") {
				return errors.New("Invalid parameter: '--port' already redefined")
			}

			// Check port is valid
			portNum, err := strconv.ParseInt(value, 10, 32)
			if err != nil || portNum < 1 || portNum > 65535 {
				return errors.New("Invalid parameter: invalid port number")
			}
			*port = value

		default:
			if allowIP {
				return errors.New("Invalid parameter: need '--ip' or '--port'")
			} else {
				return errors.New("Invalid parameter: need '--port'")
			}
		}

		updated = append(updated, key)
	}

	return nil
}
