package main

import (
	"flag"
	"log"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	var (
		prompt     string
		debugLevel int
	)

	flag.StringVar(&prompt, "prompt", "go-scheme> ", "prompt")
	flag.IntVar(&debugLevel, "debug", int(parser.Info), "debug level")

	flag.Parse()

	p := tea.NewProgram(tui.InitialModel("go-scheme> "), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
