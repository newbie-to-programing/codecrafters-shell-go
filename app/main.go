package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

		// 1. Replace every occurrence of two single quotes with an empty string
		replaced := strings.ReplaceAll(input, "''", "")

		// 2. This regex finds either text inside single quotes OR non-space characters
		re := regexp.MustCompile(`'[^']*'|[^\s]+`)
		parts := re.FindAllString(replaced, -1)

		cleaned := make([]string, 0, len(parts))
		for _, p := range parts {
			clean := strings.ReplaceAll(p, "'", "")
			cleaned = append(cleaned, clean)
		}

		command := cleaned[0]
		args := cleaned[1:]

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
