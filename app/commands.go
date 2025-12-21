package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func handleTypeCommand(args []string) {
	if len(args) == 0 {
		return
	}

	target := args[0]

	if isBuiltin(target) {
		fmt.Printf("%s is a shell builtin\n", target)
		return
	}

	if path, found := findInPath(target); found {
		fmt.Printf("%s is %s\n", target, path)
		return
	}

	fmt.Printf("%s: not found\n", target)
}

func isBuiltin(name string) bool {
	switch name {
	case ExitCommand, EchoCommand, TypeCommand, PwdCommand:
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

func handleExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s: command not found\n", command)
	}
}

func handlePwdCommand() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("cannot print working directory: %v\n", err)
		return
	}

	fmt.Printf("%v\n", dir)
}
