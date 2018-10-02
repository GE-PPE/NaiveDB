package main

import (
	"fmt"
	"io"
	"net"

	"../database"
	"../parser"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConn(conn net.Conn, db *database.NaiveDB) {

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
			value, _ := db.Get(parseResult.Key)
			conn.Write([]byte(value + "\n"))
			fmt.Println("received a get request")
		case parser.SetCommand:
			db.Set(parseResult.Key, parseResult.Value)
			conn.Write([]byte("OK\n"))
		case parser.DeleteCommand:
			db.Delete(parseResult.Key)
			conn.Write([]byte("OK\n"))
		}
	}
}

func main() {
	db := database.New("/tmp/testdatabase.db")
	listener, err := net.Listen("tcp", ":6442")

	checkError(err)
	fmt.Println("Server Listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
		} else {
			go func() {
				handleConn(conn, db)
				conn.Close()
			}()
		}
	}
}
