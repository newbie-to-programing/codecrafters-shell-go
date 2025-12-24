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

func processDoubleQuotes(content string) string {
	re := regexp.MustCompile(`\\(["\\])`)
	return re.ReplaceAllString(content, "$1")
}

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

		// The pattern identifies 4 types of "blobs" that form arguments:
		// 1. Double quotes: "(content)"
		// 2. Single quotes: '(content)'
		// 3. Escaped char: \. (outside quotes)
		// 4. Normal text:   [^\s"'\\]+
		pattern := `(?s)"([^"\\]*(?:\\.[^"\\]*)*)"|'([^']*)'|\\(.)|([^\s"'\\]+)`
		re := regexp.MustCompile(pattern)

		// Get indexes to handle concatenation (touching vs. separated by space)
		allMatches := re.FindAllStringSubmatchIndex(input, -1)
		if len(allMatches) == 0 {
			return
		}

		var args []string
		var currentArg strings.Builder
		var lastEnd int = -1

		for _, m := range allMatches {
			matchStart := m[0]
			matchEnd := m[1]

			// If there is a gap (whitespace) between this match and the last,
			// the previous argument is finished.
			if lastEnd != -1 && matchStart > lastEnd {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}

			// Check which capturing group matched and extract the literal content
			if m[2] != -1 {
				// Group 1: Double Quotes. Content is input[m[2]:m[3]]
				// Note: For now, we treat backslashes inside double quotes as literal
				// per your current stage requirements.
				currentArg.WriteString(processDoubleQuotes(input[m[2]:m[3]]))
			} else if m[4] != -1 {
				// Group 2: Single Quotes. Content is input[m[4]:m[5]]
				currentArg.WriteString(input[m[4]:m[5]])
			} else if m[6] != -1 {
				// Group 3: Backslash outside quotes. Content is input[m[6]:m[7]]
				// This handles: \ , \n, \\, \' etc.
				currentArg.WriteString(input[m[6]:m[7]])
			} else if m[8] != -1 {
				// Group 4: Normal unquoted text
				currentArg.WriteString(input[m[8]:m[9]])
			}

			lastEnd = matchEnd
		}

		// Always add the final building argument
		args = append(args, currentArg.String())

		command := args[0]
		otherArgs := args[1:]

		switch command {
		case ExitCommand:
			os.Exit(0)
		case EchoCommand:
			fmt.Println(strings.Join(otherArgs, " "))
		case TypeCommand:
			handleTypeCommand(otherArgs)
		case PwdCommand:
			handlePwdCommand()
		case CdCommand:
			handleCdCommand(otherArgs)
		default:
			handleExternalCommand(command, otherArgs)
		}
	}
}
