package main

import (
	"fmt"
	"bufio"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	//fmt.Print("$ ")

	// Captures the user's command in the "command" variable
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println(command[:len(command)-1] + ": command not found")
}
