package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		switch command {
		case ExitCommand:
			os.Exit(0)
		case EchoCommand:
			fmt.Println(strings.Join(args, " "))
		case TypeCommand:
			handleTypeCommand(args)
		case PwdCommand:
			handlePwdCommand()
		case CdCommand:
			handleCdCommand(args)
		default:
			handleExternalCommand(command, args)
		}
	}
}
