package main

import (
	"fmt"
	"io"
	"net"

	"./parser"
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

<<<<<<< HEAD:naivedb/server.go
		line := string(buffer[:bytesRead])
		parseResult, err := parser.Parse(line)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
=======
		str := string(buffer[:bytesRead])
		splittedStr := strings.Fields(str)
		if len(splittedStr) < 2 {
			fmt.Println("Malformed")
>>>>>>> 4dd1bb8d0f85d1cfacc67c0a61970ee23c59fd9a:naivedb/src/server.go
			continue
		}
		fmt.Println(parseResult)

<<<<<<< HEAD:naivedb/server.go
		switch parseResult.Command {
		case parser.GetCommand:
			conn.Write([]byte(m[parseResult.Key] + "\n"))
=======
		if splittedStr[0] == "get" {
			value := m[key]
			conn.Write([]byte(value + "\n"))
			fmt.Println(m[key])
>>>>>>> 4dd1bb8d0f85d1cfacc67c0a61970ee23c59fd9a:naivedb/src/server.go
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
