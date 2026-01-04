package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleEchoCommand(args []string) string {
	return fmt.Sprintln(strings.Join(args, " "))
}

func handleTypeCommand(args []string) string {
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
	case ExitCommand, EchoCommand, TypeCommand, PwdCommand, CdCommand:
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

func handleExternalCommand(commands []Command) (string, string, error) {
	if len(commands) == 1 {
		c := commands[0]
		cmd := exec.Command(c.Path, c.Args...)

		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf

		err := cmd.Run()
		if err != nil {
			return fmt.Sprintf("%s: command not found\n", c.Path), outBuf.String(), errors.New(errBuf.String())
		}

		return outBuf.String(), outBuf.String(), nil
	}

	cmd1 := exec.Command(commands[0].Path, commands[0].Args...)
	cmd2 := exec.Command(commands[1].Path, commands[1].Args...)

	stdout, err := cmd1.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	cmd2.Stdin = stdout
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	// 3. Start the commands
	// Start the first command
	if err = cmd1.Start(); err != nil {
		return "", "", err
	}

	// Start the second command (which will read from the pipe)
	if err = cmd2.Start(); err != nil {
		return "", "", err
	}

	cmd1.Wait()
	cmd2.Wait()

	return "", "", nil
}

func handlePwdCommand() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("cannot print working directory: %v\n", err)
	}

	return fmt.Sprintf("%v\n", dir)
}

func handleCdCommand(args []string) string {
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
