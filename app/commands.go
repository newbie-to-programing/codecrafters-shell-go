package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ExecutePipeline(commands []Command) {
	var lastStdout io.Reader = os.Stdin

	for i, c := range commands {
		isLast := i == len(commands)-1

		var currentStdout io.Writer
		var currentStderr io.Writer = os.Stderr
		var pipeReadEnd io.ReadCloser
		var pipeWriteEnd io.WriteCloser

		if isLast {
			if c.RedirectOp != "" {
				flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				if strings.Contains(c.RedirectOp, ">>") {
					flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
				}

				file, err := os.OpenFile(c.OutputFile, flags, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "shell: %v\n", err)
					return
				}
				defer file.Close()

				if strings.HasPrefix(c.RedirectOp, "2") {
					currentStderr = file
					currentStdout = os.Stdout
				} else {
					currentStdout = file
					// Stderr remains os.Stderr
				}
			} else {
				currentStdout = os.Stdout
			}
		} else {
			r, w, _ := os.Pipe()
			currentStdout = w
			pipeWriteEnd = w
			pipeReadEnd = r
		}

		if isBuiltin(c.Path) {
			output := handleBuiltinCommand(c.Path, c.Args)
			fmt.Fprint(currentStdout, output)
			if pipeWriteEnd != nil {
				pipeWriteEnd.Close()
			}
		} else {
			cmd := exec.Command(c.Path, c.Args...)
			cmd.Stdin = lastStdout
			cmd.Stdout = currentStdout
			cmd.Stderr = currentStderr

			if isLast {
				cmd.Run()
			} else {
				cmd.Start()
				go func(cmd *exec.Cmd, w io.WriteCloser) {
					cmd.Wait()
					w.Close()
				}(cmd, pipeWriteEnd)
			}
		}

		if i > 0 {
			if closer, ok := lastStdout.(io.Closer); ok {
				closer.Close()
			}
		}

		lastStdout = pipeReadEnd
	}
}

func handleBuiltinCommand(command string, args []string) string {
	switch command {
	case ExitCommand:
		os.Exit(0)
	case EchoCommand:
		return HandleEchoCommand(args)
	case TypeCommand:
		return HandleTypeCommand(args)
	case PwdCommand:
		return HandlePwdCommand()
	case CdCommand:
		return HandleCdCommand(args)
	}

	return ""
}

func HandleEchoCommand(args []string) string {
	return fmt.Sprintln(strings.Join(args, " "))
}

func HandleTypeCommand(args []string) string {
	if len(args) == 0 {
		return ""
	}

	target := args[0]

	if isBuiltin(target) {
		return fmt.Sprintf("%s is a shell builtin\n", target)
	}

	if path, found := findInPath(target); found {
		return fmt.Sprintf("%s is %s\n", target, path)
	}

	return fmt.Sprintf("%s: not found\n", target)
}

func isBuiltin(name string) bool {
	switch name {
	case ExitCommand, EchoCommand, TypeCommand, PwdCommand, CdCommand, HistoryCommand:
		return true
	}
	return false
}

func findInPath(name string) (string, bool) {
	dirs := filepath.SplitList(os.Getenv("PATH"))

	for _, dir := range dirs {
		fullPath := filepath.Join(dir, name)
		info, err := os.Stat(fullPath)

		// Ensure it exists, isn't a folder, and is executable (0111)
		if err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0 {
			return fullPath, true
		}
	}
	return "", false
}

func HandleExternalCommand(command Command) (string, string, error) {
	cmd := exec.Command(command.Path, command.Args...)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("%s: command not found\n", command.Path), outBuf.String(), errors.New(errBuf.String())
	}

	return outBuf.String(), outBuf.String(), nil
}

func HandlePwdCommand() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("cannot print working directory: %v\n", err)
	}

	return fmt.Sprintf("%v\n", dir)
}

func HandleCdCommand(args []string) string {
	if len(args) == 0 {
		return ""
	}

	targetDir := args[0]

	if targetDir == "~" {
		home, _ := os.UserHomeDir()
		targetDir = home
	}

	err := os.Chdir(targetDir)
	if err != nil {
		return fmt.Sprintf("cd: %v: No such file or directory\n", targetDir)
	}

	return ""
}

func AddToHistoryCommands(historyCommands []string, commands []Command) []string {
	for _, command := range commands {
		historyCommands = append(historyCommands, formatHistoryCommand(command))
	}

	return historyCommands
}

func formatHistoryCommand(command Command) string {
	if len(command.Args) == 0 {
		return fmt.Sprintf("%v\n", command.Path)
	}
	return fmt.Sprintf("%v %v\n", command.Path, strings.Join(command.Args, " "))
}

func HandleHistoryCommand(args []string, history []string) []string {
	if len(args) > 2 {
		return history
	}

	// Use a switch for clearer intent
	var command string
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "-r":
		return readHistoryFromFile(args, history)
	case "-w":
		saveHistoryToFile(args, history)
	case "-a":
		appendHistoryToFile(args, history)
	default:
		printHistory(args, history)
	}

	return history
}

func readHistoryFromFile(args []string, history []string) []string {
	if len(args) < 2 {
		return history
	}

	filePath := args[1]
	file, err := os.Open(filePath)
	if err != nil {
		return history
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			history = append(history, fmt.Sprintf("%v\n", line))
		}
	}

	return history
}

func saveHistoryToFile(args []string, history []string) {
	if len(args) < 2 {
		return
	}

	filePath := args[1]
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	for _, command := range history {
		fmt.Fprint(file, command)
	}
}

func printHistory(args []string, history []string) {
	limit := len(history)

	if len(args) > 0 {
		if val, err := strconv.Atoi(args[0]); err == nil {
			limit = val
		}
	}

	startIndex := len(history) - limit
	if startIndex <= 0 {
		startIndex = 0
	}

	for i := startIndex; i < len(history); i++ {
		fmt.Printf("%v  %v", i+1, history[i])
	}
}

func appendHistoryToFile(args []string, history []string) {
	if len(args) < 2 {
		return
	}

	historyToBeAppended := make([]string, 0)

	prevSeen := -1
	currSeen := -1
	for i, command := range history {
		if strings.HasPrefix(command, "history -a") {
			prevSeen = currSeen
			currSeen = i
		}
	}

	for i := prevSeen + 1; i <= currSeen; i++ {
		historyToBeAppended = append(historyToBeAppended, history[i])
	}

	filePath := args[1]
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	for _, command := range historyToBeAppended {
		fmt.Fprint(file, command)
	}
}

func LoadHistoryOnStartup() []string {
	filePath := os.Getenv("HISTFILE")
	return HandleHistoryCommand([]string{"-r", filePath}, make([]string, 0))
}

func SaveHistoryOnExit(history []string) {
	filePath := os.Getenv("HISTFILE")
	HandleHistoryCommand([]string{"-w", filePath}, history)
}
