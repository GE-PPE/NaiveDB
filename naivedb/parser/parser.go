package parser

import (
	"fmt"
	"strings"
)

// Command - A NaiveDB command that involves the database.
type Command int

// MetaCommand - A NaiveDB command that involves the DB client.
type MetaCommand int
type parseToken string
type parseState int

const (
	ExitMetaCommand MetaCommand = 0

	GetCommand    Command = 0
	SetCommand    Command = 1
	DeleteCommand Command = 2

	getToken    parseToken = "get"
	setToken    parseToken = "set"
	deleteToken parseToken = "delete"
	exitToken   parseToken = ".exit"

	commandParseState parseState = 0
	keyParseState     parseState = 1
	valueParseState   parseState = 2
	errorParseState   parseState = 3

	invalidCommandString        string = "invalid command"
	invalidArgumentLengthString string = "invalid argument length"
)

var tokenToCommand = map[parseToken]Command{
	getToken:    GetCommand,
	setToken:    SetCommand,
	deleteToken: DeleteCommand,
}

var tokenToMetaCommand = map[parseToken]MetaCommand{
	exitToken: ExitMetaCommand,
}

type ParseResult struct {
	Command Command
	Key     string
	Value   string
}

type ParseMetaResult struct {
	MetaCommand MetaCommand
}

type parseError struct {
	argument string
	reason   string
}

func (p *parseError) Error() string {
	return fmt.Sprintf("error at argument %s: %s", p.argument, p.reason)
}

// IsPossibleMetaCommand - Decides whether the line could be a metacommand or not.
func IsPossibleMetaCommand(line string) bool {
	return strings.HasPrefix(line, ".")
}

func ParseMeta(line string) (*ParseMetaResult, error) {
	strippedLine := strings.TrimSpace(line)
	metaCommand, ok := tokenToMetaCommand[parseToken(strippedLine)]

	if !ok {
		return nil, &parseError{strippedLine, invalidCommandString}
	}

	return &ParseMetaResult{metaCommand}, nil
}

func Parse(line string) (*ParseResult, error) {
	splittedLine := strings.Fields(line)

	state := commandParseState
	result := ParseResult{}

	for _, word := range splittedLine {
		switch state {
		case commandParseState:
			if command, ok := tokenToCommand[parseToken(word)]; ok {
				result.Command = command
				state = keyParseState
			} else {
				return nil, &parseError{word, invalidCommandString}
			}
		case keyParseState:
			result.Key = word
			state = valueParseState
		case valueParseState:
			if result.Command != SetCommand {
				return nil, &parseError{word, invalidArgumentLengthString}
			}
			result.Value = word
			state = errorParseState
		default:
			return nil, &parseError{word, invalidArgumentLengthString}
		}
	}
	if state == keyParseState {
		return nil, &parseError{"", invalidArgumentLengthString}
	} else if result.Command == SetCommand && state == valueParseState {
		return nil, &parseError{"", invalidArgumentLengthString}
	}
	return &result, nil
}
