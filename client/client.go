package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "6379"
	TYPE = "tcp"
)

type Client struct {
	Name string
	UserName string
	Password string
}

func main() {
	fmt.Println("!!Welcome to super secure terminal chat app!!")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEnter your name: ")
	clientName, err := reader.ReadString('\n')
	checkErr(err, "Error reading name")
	clientName = strings.TrimSpace(clientName)

	fmt.Print("\nEnter room Id: ")
	roomID, err := reader.ReadString('\n')
	checkErr(err, "Error reading roomid")
	roomID = strings.TrimSpace(roomID)

	conn, err := net.Dial(TYPE, HOST+":"+PORT)
	checkErr(err, "Unable to connect")

	fmt.Fprintln(conn, clientName)
	fmt.Fprintln(conn, roomID)

	go handleMessages(conn)

	for {
		msgReader := bufio.NewReader(os.Stdin)
		msg, _ := msgReader.ReadString('\n')
		if msg == "quit" {
			fmt.Println("Terminating session")
			conn.Close()
			os.Exit(0)
		}
		fmt.Fprintf(conn, msg)
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panic(msg, err)
	}
}

func handleMessages(conn net.Conn) {
	buf := bufio.NewReader(conn)

	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed")
			return
		}
		fmt.Print(message)
	}
}
