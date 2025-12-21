package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")

		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		commandName := command[0:4]
		switch commandName {
		case ExitCommand:
			break
		case EchoCommand:
			processEchoCommand(command)
		case TypeCommand:
			processTypeCommand(command)
		default:
			commandWithoutNextLine := strings.TrimSpace(command)
			arguments := strings.Split(commandWithoutNextLine, " ")
			executable := arguments[0]
			actualArguments := arguments[1:]

			pathList := os.Getenv("PATH")

			dirs := filepath.SplitList(pathList)

			hasFound := false
			for _, dir := range dirs {
				entries, err := os.ReadDir(dir)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					// Handle other real errors (like permission denied)
					log.Printf("Error reading %s: %v", dir, err)
					continue
				}

				for _, entry := range entries {
					if entry.IsDir() {
						continue
					}

					if entry.Name() != executable {
						if entry.Name() == "hello" {
							fmt.Println(executable == entry.Name())
							fmt.Printf("entry: %v, exec: %q", entry.Name(), executable)
						}
						continue
					}

					info, err := entry.Info()
					if err != nil {
						continue
					}

					mode := info.Mode()
					if mode.Perm()&0111 != 0 {
						cmd := exec.Command(filepath.Join(dir, entry.Name()), arguments...)
						// Redirect the command's output directly to the current terminal
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						// Also connect Stdin if the program needs user input
						cmd.Stdin = os.Stdin

						err = cmd.Run()
						if err != nil {
							fmt.Printf("%s: command not found\n", executable)
						}
						hasFound = true
					}
				}

				if hasFound {
					break
				}
			}
		}
	}
}

func processEchoCommand(command string) {
	fmt.Print(command[5:])
}

func processTypeCommand(command string) {
	builtin := command[5 : len(command)-1]
	if builtin == "echo" || builtin == "type" || builtin == "exit" {
		fmt.Printf("%v is a shell builtin\n", builtin)
	} else {
		pathList := os.Getenv("PATH")

		dirs := filepath.SplitList(pathList)

		hasFound := false
		for _, dir := range dirs {
			entries, err := os.ReadDir(dir)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				// Handle other real errors (like permission denied)
				log.Printf("Error reading %s: %v", dir, err)
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				if entry.Name() != builtin {
					continue
				}

				info, err := entry.Info()
				if err != nil {
					continue
				}

				mode := info.Mode()
				if mode.Perm()&0111 != 0 {
					fmt.Printf("%s is %s\n", builtin, filepath.Join(dir, entry.Name()))
					hasFound = true
				}
			}

			if hasFound {
				break
			}
		}

		if !hasFound {
			fmt.Printf("%s: not found\n", builtin)
		}
	}
}
