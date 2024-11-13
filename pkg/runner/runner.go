package runner

import (
	"app/pkg/runner/lexer"
	"fmt"
	"os"
	"os/exec"
)

func Run(command string) error {
	commandPart := lexer.Lex(command)
	if len(commandPart) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(commandPart[0], commandPart[1:]...)
	cmd.Env = os.Environ()
	return cmd.Run()
}
