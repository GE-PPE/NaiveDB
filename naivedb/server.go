package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConn(conn net.Conn, m *map[string]string) {
	defer conn.Close()

	buffer := make([]byte, 4096)

	for {
		bytesRead, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("RECEIVED EOF, ignored %d bytes\n", bytesRead)
			}
			return
		}
		fmt.Println(string(buffer))

		str := string(buffer)
		splittedStr := strings.Split(str, " ")

		if splittedStr[0] == "get" {
			fmt.Println("received a get request")
		} else if splittedStr[0] == "set" {
			fmt.Println("received a set request")
		} else if splittedStr[0] == "delete" {
			fmt.Println("received a delete request")
		} else {
			fmt.Println("received anything")
		}

		for i := range buffer {
			buffer[i] = 0
		}
	}
}

func main() {
	m := make(map[string]string)
	listener, err := net.Listen("tcp", ":6442")

	checkError(err)
	fmt.Println("Server Listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
		}

		go handleConn(conn, &m)
	}
}
