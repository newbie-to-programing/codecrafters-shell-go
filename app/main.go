package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type CommandResult struct {
	Output string
	Stdout string
	Stderr error
}

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	completer := NewUnifiedCompleter([]string{"echo", "exit"})

	l, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
		Listener:     &MyListener{},
	})
	completer.SetInstance(l)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		fmt.Print("$ ")

		//input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		//if err != nil {
		//	continue
		//}
		input, err := l.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				// User pressed Ctrl+C
				continue
			} else if errors.Is(err, io.EOF) {
				// User pressed Ctrl+D
				break
			}
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		args := parseArguments(input)
		if len(args) == 0 {
			continue
		}

		var command1 string
		var otherArgs1 []string
		var command2 string
		var otherArgs2 []string
		var cmdArgs1 []string
		var cmdArgs2 []string
		var redirectOp string
		var outFile string
		cmdArgs1, redirectOp, outFile = extractRedirection(args)
		if redirectOp == "" {
			cmdArgs1, cmdArgs2 = extractPipeline(args)
		}

		command1 = cmdArgs1[0]
		otherArgs1 = cmdArgs1[1:]
		if len(cmdArgs2) > 0 {
			command2 = cmdArgs2[0]
			otherArgs2 = cmdArgs2[1:]
		}

		var res CommandResult
		switch command1 {
		case ExitCommand:
			os.Exit(0)
		case EchoCommand:
			ret := handleEchoCommand(otherArgs1)
			res.Output = ret
			res.Stdout = ret
		case TypeCommand:
			ret := handleTypeCommand(otherArgs1)
			res.Output = ret
			res.Stdout = ret
		case PwdCommand:
			ret := handlePwdCommand()
			res.Output = ret
			res.Stdout = ret
		case CdCommand:
			ret := handleCdCommand(otherArgs1)
			res.Output = ret
			res.Stdout = ret
		default:
			res.Output, res.Stdout, res.Stderr = handleExternalCommand(command1, otherArgs1, command2, otherArgs2)
		}

		handleOutput(res, redirectOp, outFile)
	}
}

func handleOutput(result CommandResult, redirectOp, filename string) {
	// If no redirection, print to standard output
	if redirectOp == "" {
		fmt.Print(result.Output)
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
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
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
