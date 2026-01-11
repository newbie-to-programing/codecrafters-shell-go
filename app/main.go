package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	historyCommands = make([]string, 0)
)

func main() {
	l := InitializeReadline()
	defer l.Close()

	historyCommands = LoadHistoryOnStartup()

	for {
		fmt.Print("$ ")

		input, err := l.Readline()
		if err != nil {
			action := HandleReadlineError(err)
			if action == StopLoop {
				break
			}
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		commands := ParseInput(input)
		if len(commands) <= 0 {
			continue
		}

		historyCommands = AddToHistoryCommands(historyCommands, commands)

		if len(commands) > 1 {
			ExecutePipeline(commands)
			continue
		}

		c := commands[0]
		command := c.Path
		otherArgs := c.Args

		var res CommandResult
		switch command {
		case ExitCommand:
			SaveHistoryOnExit(historyCommands)
			os.Exit(0)
		case EchoCommand:
			ret := HandleEchoCommand(otherArgs)
			res.Output = ret
			res.Stdout = ret
		case TypeCommand:
			ret := HandleTypeCommand(otherArgs)
			res.Output = ret
			res.Stdout = ret
		case PwdCommand:
			ret := HandlePwdCommand()
			res.Output = ret
			res.Stdout = ret
		case CdCommand:
			ret := HandleCdCommand(otherArgs)
			res.Output = ret
			res.Stdout = ret
		case HistoryCommand:
			historyCommands = HandleHistoryCommand(otherArgs, historyCommands)
		default:
			res.Output, res.Stdout, res.Stderr = HandleExternalCommand(c)
		}

		handleOutput(res, c.RedirectOp, c.OutputFile)
	}
}

func handleOutput(result CommandResult, redirectOp, filename string) {
	// If no redirection, print to standard output
	if redirectOp == "" {
		if result.Output != "" {
			fmt.Print(result.Output)
		}
		return
	}

	var toBeWrittenToTerminal string
	var toBeWrittenToFile string
	if redirectOp == "2>" || redirectOp == "2>>" {
		toBeWrittenToTerminal = result.Stdout
		if result.Stderr != nil {
			toBeWrittenToFile = result.Stderr.Error()
		}
	} else {
		toBeWrittenToFile = result.Stdout
		if result.Stderr != nil {
			toBeWrittenToTerminal = result.Stderr.Error()
		}
	}

	if toBeWrittenToTerminal != "" {
		fmt.Print(toBeWrittenToTerminal)
	}

	// Handle Redirection
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if redirectOp == ">>" || redirectOp == "1>>" || redirectOp == "2>>" {
		flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shell: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprint(file, toBeWrittenToFile)
}
