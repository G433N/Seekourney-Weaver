package indexing

import (
	"errors"
	"seekourney/utils"
	"strconv"
)

// IntoPort converts an integer to a port.
func IntoPort(integer uint) (utils.Port, bool) {

	if integer < uint(utils.MININDEXERPORT) ||
		integer > uint(utils.MAXINDEXERPORT) {
		return 0, false
	}

	return utils.Port(integer), true
}

// IsValidPort checks if port value is within designated range for indexer API.
func IsValidPort(port utils.Port) bool {
	_, ok := IntoPort(uint(port))
	return ok
}

func GetPort(args []string) (utils.Port, error) {

	if len(args) < 2 {
		return 0, errors.New("to few arguments")
	}

	num, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return 0, err
	}

	port, ok := IntoPort(uint(num))
	if !ok {
		return 0, errors.New("port out of range")
	}

	return utils.Port(port), nil
}
