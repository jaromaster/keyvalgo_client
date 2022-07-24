package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// getPassword allows user to enter password, then returns it
func getPassword() (string, error) {
	fmt.Print("Enter password: ")
	password_bytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(password_bytes), nil
}

// createConnection creates tls connection to database server
func createConnection(ip string, port int) *tls.Conn {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), conf)
	if err != nil {
		panic(fmt.Sprintf("could not connect to database at %s:%d", ip, port))
	}

	return conn
}

// printHelp shows help for users
func printHelp() {
	fmt.Println("COMMANDS:")
	fmt.Println("set key:value (add new key-value pair)")
	fmt.Println("get key (get value of key)")
	fmt.Println("delete key (delete key-value pair)")
	fmt.Println("import (load data.csv)")
	fmt.Println("export (persist to data.csv)")
	fmt.Println("exit (close connection and exit)")
}

// mainLoop handles user commands and sends them to database server
func mainLoop(password string, ip string, port int) {

	for {
		// configure tls for secure connection
		conn := createConnection(ip, port)
		conn_reader := bufio.NewReader(conn)
		cmd_reader := bufio.NewReader(os.Stdin)

		// auth
		conn.Write([]byte(password + "\n"))                // send password
		auth_response, err := conn_reader.ReadString('\n') // get response
		if err != nil {
			fmt.Println(err)
		}
		if strings.HasPrefix(auth_response, "Auth failed") {
			fmt.Println()
			fmt.Println(auth_response)
			return
		}

		// command
		fmt.Print("> ")
		command, err := cmd_reader.ReadString('\n') // read command
		if err != nil {
			fmt.Println(err)
		}
		command = strings.TrimSpace(command)

		if command == "help" {
			conn.Close()
			printHelp()
			continue
		}

		conn.Write([]byte(command + "\n"))                // send command
		response_msg, err := conn_reader.ReadString('\n') // read response
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(strings.TrimSpace(response_msg))

		// exit if command is exit
		if command == "exit" {
			return
		}
	}
}

func main() {
	// get ip and password
	if len(os.Args) != 3 {
		panic("Usage: kvg-client IP PORT")
	}

	// get ip and port from args
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("error parsing port from args")
	}

	// read password
	password, err := getPassword()
	if err != nil {
		panic(err)
	}
	fmt.Println()

	// user input
	mainLoop(strings.TrimSpace(password), ip, port)
}
