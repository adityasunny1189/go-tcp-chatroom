package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "6379"
	TYPE = "tcp"
)

type Client struct {
	Conn   net.Conn
	RoomId string
	Name   string
}

func main() {
	clients := make([]Client, 0)
	log.Println("Logs from your program will appear here!")
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	checkErr(err, "Failed to bind to port 6379 ")
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		checkErr(err, "Error accepting connection: ")
		go handleConnection(conn, &clients)
		defer conn.Close()
	}
}

func handleConnection(conn net.Conn, clients *[]Client) {
	reader := bufio.NewReader(conn)

	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	roomID, _ := reader.ReadString('\n')
	roomID = strings.TrimSpace(roomID)

	newClient := Client{
		Conn:   conn,
		Name:   name,
		RoomId: roomID,
	}

	log.Println("New Client: ", newClient)
	*clients = append(*clients, newClient)
	log.Println("Clients: ", clients)
	go handleMessages(newClient, clients, roomID)
}

func handleMessages(c Client, clients *[]Client, roomId string) {
	buf := bufio.NewReader(c.Conn)

	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			RemoveClient(c.Conn, clients)
			return
		}

		for _, client := range *clients {
			if client.RoomId == roomId && client.Conn != c.Conn {
				fmt.Fprintf(client.Conn, c.Name+": "+message)
			}
		}
	}
}

func RemoveClient(conn net.Conn, clients *[]Client) {
	for i, c := range *clients {
		if c.Conn == conn {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
			break
		}
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panic(msg, err)
	}
}
