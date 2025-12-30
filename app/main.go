package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type CommandResult struct {
	Output string
	Err    error
}

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		args := parseArguments(input)
		if len(args) == 0 {
			continue
		}

		cmdArgs, redirectOp, outFile := extractRedirection(args)

		var res CommandResult
		command := cmdArgs[0]
		otherArgs := cmdArgs[1:]

		switch command {
		case ExitCommand:
			os.Exit(0)
		case EchoCommand:
			res.Output = handleEchoCommand(otherArgs)
		case TypeCommand:
			res.Output, res.Err = handleTypeCommand(otherArgs)
		case PwdCommand:
			res.Output, res.Err = handlePwdCommand()
		case CdCommand:
			res.Output, res.Err = handleCdCommand(otherArgs)
		default:
			res.Output, res.Err = handleExternalCommand(command, otherArgs)
		}

		handleOutput(res, redirectOp, outFile)
	}
}

func parseArguments(input string) []string {
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
		return nil
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

	return args
}

func processDoubleQuotes(content string) string {
	re := regexp.MustCompile(`\\(["\\])`)
	return re.ReplaceAllString(content, "$1")
}

func extractRedirection(args []string) (cmdArgs []string, redirectOp string, outputFile string) {
	for i, arg := range args {
		// Check for redirection operators
		if arg == ">" || arg == "1>" || arg == ">>" || arg == "2>" {
			redirectOp = arg

			// Extract command arguments (everything before the operator)
			cmdArgs = args[:i]

			// Safely extract filename (the argument immediately following the operator)
			if i+1 < len(args) {
				outputFile = args[i+1]
			}

			return cmdArgs, redirectOp, outputFile
		}
	}

	// If no operator was found, all args are command args
	return args, "", ""
}

func handleOutput(result CommandResult, redirectOp, filename string) {
	// If no redirection, print to standard output
	if redirectOp == "" {
		fmt.Print(result.Output)
		return
	}

	var toBeWrittenToTerminal string
	var toBeWrittenToFile string

	if redirectOp == "2>" {
		toBeWrittenToFile = result.Err.Error()
		toBeWrittenToTerminal = result.Output
	} else {
		toBeWrittenToFile = result.Output
		toBeWrittenToTerminal = result.Err.Error()
	}

	if toBeWrittenToTerminal != "" {
		fmt.Print(toBeWrittenToTerminal)
	}

	// Handle Redirection
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	// if redirectOp == ">>" { flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND }

	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shell: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprint(file, toBeWrittenToFile)
}
