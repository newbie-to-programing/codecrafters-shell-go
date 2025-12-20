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
		} else {
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
