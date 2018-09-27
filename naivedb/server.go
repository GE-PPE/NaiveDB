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

func handleConn(conn net.Conn, m map[string]string) {
	defer conn.Close()

	buffer := make([]byte, 4096)

	for {
		for i := range buffer {
			buffer[i] = 0
		}

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
		if len(splittedStr) < 2 {
			fmt.Println("Malformed")
			continue
		}
		key := splittedStr[1]

		if splittedStr[0] == "get" {
			value := m[key]
			conn.Write([]byte(value))
			fmt.Println(m[key])
			fmt.Println("received a get request")
		} else if splittedStr[0] == "set" {
			if len(splittedStr) < 3 {
				fmt.Println("Malformed")
			}
			value := splittedStr[2]
			m[key] = value
			conn.Write([]byte("OK\n"))
			fmt.Println("received a set request")
		} else if splittedStr[0] == "delete" {
			m[key] = ""
			conn.Write([]byte("OK\n"))
			fmt.Println("received a delete request")
		} else {
			fmt.Println("received anything")
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
		} else {
			go handleConn(conn, m)
		}
	}
}
