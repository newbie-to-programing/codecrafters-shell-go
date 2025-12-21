package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for true {
		// TODO: Uncomment the code below to pass the first stage
		fmt.Print("$ ")

		// Captures the user's command in the "command" variable
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		if command[:len(command)-1] == "exit" {
			break
		} else if command[0:4] == "echo" {
			fmt.Print(command[5:])
		} else if command[0:4] == "type" {
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
		} else {
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
