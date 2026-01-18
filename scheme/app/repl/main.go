package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/lexer"
	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser"
)

func main() {

	var (
		srcPath    string
		dstPath    string
		prompt     string
		debugLevel int
	)
	flag.StringVar(&srcPath, "src", "", "source file")
	flag.StringVar(&dstPath, "dst", "", "destination file")
	flag.StringVar(&prompt, "prompt", "go-scheme> ", "prompt")
	flag.IntVar(&debugLevel, "debug", int(parser.Info), "debug level")

	flag.Parse()

	var (
		in  = os.Stdin
		out = os.Stdout
		err error
	)

	if srcPath != "" {
		in, err = os.Open(srcPath)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()
	}
	if dstPath != "" {
		out, err = os.Create(dstPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	p := parser.New(
		context.Background(),
		lexer.New(in),
		parser.WithPrompt(prompt),
		parser.WithVerbose(parser.VerboseLevel(debugLevel)))

	p.Repl()

}
