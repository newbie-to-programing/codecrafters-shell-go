package main

import (
	"fmt"
	"bufio"
	"os"
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
			builtin := command[5:len(command)-1]
			if builtin == "echo" || builtin == "type" || builtin == "exit" {
				fmt.Printf("%v is a shell builtin" + builtin)
			} else {
				fmt.Printf("%v: not found" + builtin)
			}
		} else {
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
