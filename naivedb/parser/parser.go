package parser

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
	valueParseState   parseState = 2
	errorParseState   parseState = 3
)

var tokenToCommand = map[parseToken]Command{
	getToken:    GetCommand,
	setToken:    SetCommand,
	deleteToken: DeleteCommand,
}

type ParseResult struct {
	Command Command
	Key     string
	Value   string
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
				result.Command = command
				state = keyParseState
			} else {
				return nil, &parseError{word, "invalid command"}
			}
		case keyParseState:
			result.Key = word
			state = valueParseState
		case valueParseState:
			if result.Command != SetCommand {
				return nil, &parseError{word, "invalid argument length"}
			}
			result.Value = word
			state = errorParseState
		default:
			return nil, &parseError{word, "invalid argument length"}
		}
	}
	if state == keyParseState {
		return nil, &parseError{"", "invalid argument length"}
	} else if result.Command == SetCommand && state == valueParseState {
		return nil, &parseError{"", "invalid argument length"}
	}
	return &result, nil
}
