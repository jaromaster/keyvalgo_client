package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/term"
)

// getPassword allows user to enter password, then returns it
func getPassword() (string, error) {
	fmt.Print("Enter DB password: ")
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
	conn, err := tls.Dial("tcp", ":8000", conf)
	if err != nil {
		panic(fmt.Sprintf("could not connect to database at %s:%d", ip, port))
	}

	return conn
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

	// configure tls for secure connection
	conn := createConnection(ip, port)
	defer conn.Close()

	// let user enter commands

	conn_reader := bufio.NewReader(conn)

	t, _ := conn_reader.ReadString(' ')
	fmt.Println(t)
	conn.Write([]byte(password + "\n"))
	conn.Write([]byte("get hans\n"))
	t, _ = conn_reader.ReadString('\n')
	fmt.Println(t)
}
