package sesh

import (
	"fmt"

	"github.com/glycerine/zygomys/v9/zygo"
)

type ManPageNotFoundError struct {
	Command string
}

func (e *ManPageNotFoundError) Error() string {
	return fmt.Sprintf("no help available for %s", e.Command)
}

func newManPages() *manPages {
	return &manPages{
		pages: make(map[string]string),
	}
}

type manPages struct {
	pages map[string]string
}

func (m *manPages) addPage(name string, text string) {
	m.pages[name] = text
}

func (m *manPages) handle(env *zygo.Zlisp, _ string, args []zygo.Sexp) (zygo.Sexp, error) {
	if len(args) != 1 {
		return zygo.SexpNull, fmt.Errorf("usage: help <command>")
	}

	sexp, ok := args[0].(*zygo.SexpStr)
	if !ok {
		return zygo.SexpNull, fmt.Errorf("usage: help <command>")
	}

	command := sexp.S
	page, ok := m.pages[command]
	if !ok {
		return nil, fmt.Errorf("no help available for %s", command)
	}

	return &zygo.SexpStr{S: page}, nil
}
