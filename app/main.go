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

		var replaced string
		var re *regexp.Regexp
		if strings.Contains(input, "\"") {
			replaced = strings.ReplaceAll(input, "\"\"", "") // double quotes
			re = regexp.MustCompile(`"[^"]*"|[^\s]+`)
		} else {
			replaced = strings.ReplaceAll(input, "''", "") // single quotes
			re = regexp.MustCompile(`'[^']*'|[^\s]+`)
		}

		parts := re.FindAllString(replaced, -1)

		cleaned := make([]string, 0, len(parts))
		for _, p := range parts {
			fmt.Printf("part: %q\n", p)
			clean := strings.ReplaceAll(p, "'", "")
			clean = strings.ReplaceAll(p, "\"", "")
			cleaned = append(cleaned, clean)
			fmt.Printf("clean: %q\n", clean)
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
