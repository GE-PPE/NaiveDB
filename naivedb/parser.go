package main

import (
	"fmt"
	"strings"
)

type Command int
type parseToken string
type parseState int

const (
	GetCommand    Command = 0
	SetCommand    Command = 1
	DeleteCommand Command = 2

	getToken    parseToken = "get"
	setToken    parseToken = "set"
	deleteToken parseToken = "delete"

	commandParseState parseState = 0
	keyParseState     parseState = 1
	valueparseState   parseState = 2
	errorParseState   parseState = 3
)

var tokenToCommand = map[parseToken]Command{
	getToken:    GetCommand,
	setToken:    SetCommand,
	deleteToken: DeleteCommand,
}

type ParseResult struct {
	command Command
	key     string
	value   string
}

type parseError struct {
	argument string
	reason   string
}

func (p *parseError) Error() string {
	return fmt.Sprintf("error at argument %s: %s", p.argument, p.reason)
}

func Parse(line string) (*ParseResult, error) {
	splittedLine := strings.Fields(line)

	state := commandParseState
	result := ParseResult{}
	// wordIsFinished := false

	for _, word := range splittedLine {
		switch state {
		case commandParseState:
			if command, ok := tokenToCommand[parseToken(word)]; ok {
				result.command = command
				state = keyParseState
			} else {
				return nil, &parseError{word, "invalid command"}
			}
		case keyParseState:
			result.key = word
			state = valueparseState
		case valueparseState:
			result.value = word
			state = errorParseState
		default:
			return nil, &parseError{word, "invalid argument length"}
		}
	}
	return &result, nil
}
