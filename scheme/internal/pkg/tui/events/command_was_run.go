package events

import "time"

const (
	CommandStartedEventName  EventName = "CommandStarted"
	CommandFinishedEventName EventName = "CommandFinished"
)

type CommandStarted struct {
	Type
	Command string
}

type CommandFinished struct {
	Type
	Output string
	Err    error
}

func NewCommandStarted(command string) CommandStarted {
	return CommandStarted{
		Type: Type{
			At: time.Now(),
		},
		Command: command,
	}
}

func NewCommandFinished(output string, err error) CommandFinished {
	return CommandFinished{
		Type: Type{
			At: time.Now(),
		},
		Output: output,
		Err:    err,
	}
}
