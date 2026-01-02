package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UnifiedCompleter struct {
	builtins []string
}

func NewUnifiedCompleter(builtins []string) *UnifiedCompleter {
	return &UnifiedCompleter{builtins: builtins}
}

func (u *UnifiedCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	typedSoFar := string(line[:pos])
	if typedSoFar == "" {
		return nil, 0
	}

	var suggestions [][]rune
	seen := make(map[string]bool) // To prevent duplicates (e.g., if 'echo' is also in /bin)

	// 1. Check Builtin Commands (echo, exit)
	for _, cmd := range u.builtins {
		if strings.HasPrefix(cmd, typedSoFar) {
			suffix := cmd[len(typedSoFar):]
			suggestions = append(suggestions, []rune(suffix+" "))
			seen[cmd] = true
		}
	}

	// 2. Check $PATH for External Executables
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			name := file.Name()
			// Only suggest if it matches prefix AND we haven't suggested it as a builtin
			if strings.HasPrefix(name, typedSoFar) && !seen[name] {
				suffix := name[len(typedSoFar):]
				suggestions = append(suggestions, []rune(suffix+" "))
				seen[name] = true
			}
		}
	}

	return suggestions, len(typedSoFar)
}

// 1. Define a struct for your listener
type MyListener struct{}

// 2. Implement the OnChange method required by the Listener interface
func (l *MyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// 9 is the ASCII code for Tab
	if key == 9 {
		currentLine := string(line)
		// Logic to detect if the user pressed tab on a non-command
		if !strings.HasPrefix("echo", currentLine) && !strings.HasPrefix("exit", currentLine) {
			fmt.Print("\a") // Play terminal bell
		}
	}
	return nil, 0, false
}
