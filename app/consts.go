package main

const (
	ExitCommand    = "exit"
	TypeCommand    = "type"
	EchoCommand    = "echo"
	PwdCommand     = "pwd"
	CdCommand      = "cd"
	HistoryCommand = "history"
)

type LoopAction int

const (
	ContinueLoop LoopAction = iota
	StopLoop
)
