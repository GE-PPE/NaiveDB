package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

type entry struct {
	Key   string
	Value string
}

var tests = []entry{
	{"a", "atlas"},
	{"b", "banana"},
	{"c", "cachorro"},
}

func TestShowcase(t *testing.T) {
	conn, err := net.Dial("tcp", ":6442")
	defer conn.Close()
	if err != nil {
		t.FailNow()
	}
	reader := bufio.NewReader(conn)

	for _, pair := range tests {
		fmt.Fprintf(conn, "set %s %s", pair.Key, pair.Value)
		message, err := reader.ReadString('\n')
		if err != nil {
			t.FailNow()
		}
		message = strings.TrimSpace(message)
		if message != "OK" {
			t.Error(
				"Did not return OK on set operation",
			)
		}

		fmt.Fprintf(conn, "get %s", pair.Key)
		message, err = reader.ReadString('\n')
		if err != nil {
			t.FailNow()
		}
		message = strings.TrimSpace(message)
		if message != pair.Value {
			t.Error(
				"For", pair.Key,
				"expected", pair.Value,
				"got", message,
			)
		}

		fmt.Fprintf(conn, "delete %s", pair.Key)
		message, err = reader.ReadString('\n')
		if err != nil {
			t.FailNow()
		}
		message = strings.TrimSpace(message)
		if message != "OK" {
			t.Error(
				"Did not return OK on delete operation",
			)
		}
	}
}
