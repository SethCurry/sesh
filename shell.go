package sesh

import (
	"os"
	"path/filepath"

	"github.com/glycerine/zygomys/v9/zygo"
)

func NewShell[T any](name string, ctx T) *Shell[T] {
	return &Shell[T]{
		Name:        name,
		interpreter: NewInterpreter(name, ctx),
	}
}

type Shell[T any] struct {
	Name        string
	interpreter *Interpreter[T]
}

func (s *Shell[T]) RegisterBasicCallable(callable BaseCallable) {
	s.interpreter.AddBasicCallable(callable.Name, callable.Handler, callable.Docs)
}

func (s *Shell[T]) RegisterCallable(callable Callable[T]) {
	s.interpreter.AddCallable(callable.Name, callable.Handler, callable.Docs)
}

func (s *Shell[T]) RegisterContextModule(module ContextModule[T]) {
	callables := module.Callables()
	for _, callable := range callables {
		s.RegisterCallable(callable)
	}
}

func (s *Shell[T]) RegisterBasicModule(module BasicModule) {
	callables := module.BasicCallables()
	for _, callable := range callables {
		s.RegisterBasicCallable(callable)
	}
}

func (s *Shell[T]) EvalScriptFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return s.interpreter.LoadText(string(content))
}

func (s *Shell[T]) EvalScriptDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err = s.EvalScriptDir(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		} else {
			err = s.EvalScriptFile(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Shell[T]) LoadDefaultScriptsDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return s.EvalScriptDir(filepath.Join(home, ".local", s.Name, "scripts.d"))
}

func (s *Shell[T]) REPL() {
	s.interpreter.REPL()
}

type BaseCallable struct {
	Name    string
	Handler func(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error)
	Docs    string
}

type Callable[T any] struct {
	Name    string
	Handler func(env *zygo.Zlisp, name string, args []zygo.Sexp, ctx T) (zygo.Sexp, error)
	Docs    string
}

type BasicModule interface {
	BasicCallables() []BaseCallable
}

type ContextModule[T any] interface {
	Callables() []Callable[T]
}
