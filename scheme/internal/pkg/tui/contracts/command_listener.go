package contracts

import "github.com/bchisham/go-lisp/scheme/internal/pkg/tui/events"

type CommandListener interface {
	OnCommandStarted(started events.CommandStarted)
	OnCommandFinished(finished events.CommandFinished)
}
