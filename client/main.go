package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	app := NewMyApp()

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("Execution error: %v", err))
	}

	// Print server responses in user output
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			res := scanner.Text()
			fmt.Println(res)

			if res == "OK bye" {
				conn.Close()
				os.Exit(0)
			}
		}
	}()

	// Send user commands to server input
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Fprintf(conn, input.Text()+"\n")
	}
}
