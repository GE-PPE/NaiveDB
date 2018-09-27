package main

import (
	"fmt"
	"io"
	"net"

	"../parser"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConn(conn net.Conn, m map[string]string) {

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

		line := string(buffer[:bytesRead])
		parseResult, err := parser.Parse(line)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}
		fmt.Println(parseResult)

		switch parseResult.Command {
		case parser.GetCommand:
			conn.Write([]byte(m[parseResult.Key] + "\n"))
			fmt.Println("received a get request")
		case parser.SetCommand:
			m[parseResult.Key] = parseResult.Value
			conn.Write([]byte("OK\n"))
		case parser.DeleteCommand:
			delete(m, parseResult.Key)
			conn.Write([]byte("OK\n"))
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
			go func() {
				handleConn(conn, m)
				conn.Close()
			}()
		}
	}
}
