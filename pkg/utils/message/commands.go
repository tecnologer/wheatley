package message

import (
	"fmt"
	"regexp"
	"strings"
)

var separatorRegEx = regexp.MustCompile(`^([^:=]*)([=:](.*))?`)

// ExtractCommand extracts the command name and its arguments from a message.
// A command is a message that starts with a `/`.
// The command name is the first part of the message, and the arguments are the rest.
//
// Examples:
//   - ExtractCommand("hello") => "", nil
//   - ExtractCommand("/start") => "start", nil
//   - ExtractCommand("/echo hello world") => "echo", []string{"hello", "world"}
func ExtractCommand(message string) (string, []string) {
	if !IsCommand(message) {
		return "", nil
	}

	for i := 1; i < len(message); i++ {
		if message[i] == ' ' {
			return message[1:i], splitArgs(message[i+1:])
		}
	}

	return message[1:], nil
}

func ExtractCommandNamedBot(message string, botName string) (string, []string) {
	cmd, args := ExtractCommand(message)
	if cmd == "" {
		return "", nil
	}

	return strings.ReplaceAll(cmd, "@"+botName, ""), args
}

// IsCommand checks if a message is a command.
// A command is a message that starts with a `/`.
func IsCommand(message string) bool {
	if message == "" {
		return false
	}

	return message[0] == '/'
}

func splitArgs(message string) []string {
	var (
		currentArg string
		args       []string
	)

	for i := 0; i < len(message); i++ {
		if message[i] == ' ' {
			args = append(args, currentArg)
			currentArg = ""

			continue
		}

		currentArg += string(message[i])
	}

	if currentArg != "" {
		args = append(args, currentArg)
	}

	return args
}

// ExtractValueFromArg extracts the argument name and its value from a string.
// The argument name is the first part of the string, and the value is the second part.
//   - The separator can be either `:`, `=`, or `:=`.
//   - If the separator is not found, the whole string is considered the argument name.
//   - If the separator is found, the argument name is the first part, and the value is the second part.
//
// Examples:
//   - ExtractValueFromArg("key value") => "key value", ""
//   - ExtractValueFromArg("key:value") => "key", "value"
//   - ExtractValueFromArg("key=value") => "key", "value"
//   - ExtractValueFromArg("key:=value") => "key", "value"
//   - ExtractValueFromArg("keyvalue") => "keyvalue", ""
func ExtractValueFromArg(arg string) (string, string) {
	if arg == "" {
		return "", ""
	}

	if isURL(arg) {
		return arg, ""
	}

	arg = strings.TrimSpace(arg)

	chunks := separatorRegEx.FindStringSubmatch(arg)

	return strings.TrimSpace(chunks[1]), strings.TrimSpace(chunks[3])
}

// ArgsToMap converts a slice of arguments to a map.
// The order slice defines the order of the arguments in case the argument name is not provided.
//
// Examples:
//   - ArgsToMap([]string{"key value"}, []string{"key"}) => map[string]string{"key": "key value"}
//   - ArgsToMap([]string{"key:value"}, []string{"key"}) => map[string]string{"key": "value"}
//   - ArgsToMap([]string{"key=value"}, []string{"key"}) => map[string]string{"key": "value"}
//   - ArgsToMap([]string{"value"}, []string{"key"}) => map[string]string{"key": "value"}
func ArgsToMap(args []string, order []string) (map[string]string, error) {
	if len(args) != len(order) {
		return nil, fmt.Errorf("args and order slices must have the same length")
	}

	argsMap := make(map[string]string)

	for i, arg := range args {
		argName, argValue := ExtractValueFromArg(arg)
		if argValue == "" {
			argValue = argName
			argName = order[i]
		}

		argsMap[argName] = argValue
	}

	return argsMap, nil
}

func isURL(arg string) bool {
	// A very basic URL check, can be improved with a proper regex if needed.
	return strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://")
}
