package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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

func handleExternalCommand(command1 string, args1 []string, command2 string, args2 []string) (string, string, error) {
	if command2 == "" {
		cmd := exec.Command(command1, args1...)

		var out bytes.Buffer
		cmd.Stdout = &out
		var errOut bytes.Buffer
		cmd.Stderr = &errOut

		err := cmd.Run()
		if err != nil {
			return fmt.Sprintf("%s: command not found\n", command1), out.String(), errors.New(errOut.String())
		}

		return out.String(), out.String(), nil
	} else {
		cmd1 := exec.Command(command1, args1...)
		cmd2 := exec.Command(command2, args2...)

		stdout, err := cmd1.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		cmd2.Stdin = stdout
		cmd2.Stdout = os.Stdout

		// 3. Start the commands
		// Start the first command
		if err = cmd1.Start(); err != nil {
			log.Fatal(err)
		}

		// Start the second command (which will read from the pipe)
		if err = cmd2.Start(); err != nil {
			log.Fatal(err)
		}

		// 4. Wait for the commands to complete
		// Wait for the first command to finish and close the stdout pipe
		if err = cmd1.Wait(); err != nil {
			log.Fatal(err)
		}

		// Wait for the second command to finish
		if err = cmd2.Wait(); err != nil {
			log.Fatal(err)
		}

		return "", "", nil
	}
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
