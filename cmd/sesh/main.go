package main

import (
	"os"

	"github.com/SethCurry/sesh"
	"github.com/rs/zerolog"
)

type Context struct{}

func main() {
	ctx := Context{}
	shell := sesh.NewShell("sesh", ctx)

	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
	logging := sesh.NewLogging(&logger)
	shell.RegisterBasicModule(logging)

	shell.REPL()
}
