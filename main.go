package main

import (
	"app/cmd/dump"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error : ", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	actions := map[string]func([]string) error{
		dump.Action: dump.Run,
	}

	listOfAction := make([]string, 0, len(actions))
	for a := range actions {
		listOfAction = append(listOfAction, a)
	}

	invalidAction := fmt.Errorf("invalid action, valid are [%s]", strings.Join(listOfAction, " | "))

	if len(args) == 0 {
		return invalidAction
	}

	run, ok := actions[args[0]]
	if !ok {
		return invalidAction
	}
	return run(args[1:])
}
