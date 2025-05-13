package sesh

import "github.com/glycerine/zygomys/v9/zygo"

func NewInterpreter[T any](name string, ctx T) *Interpreter[T] {
	inter := &Interpreter[T]{
		config:   zygo.NewZlispConfig(name),
		env:      zygo.NewZlisp(),
		manPages: newManPages(),
		ctx:      ctx,
	}

	inter.env.StandardSetup()

	inter.AddBasicCallable("help", inter.manPages.handle, "Prints help for a command")
	return inter
}

type Interpreter[T any] struct {
	config   *zygo.ZlispConfig
	env      *zygo.Zlisp
	ctx      T
	manPages *manPages
}

func (i *Interpreter[T]) LoadText(text string) error {
	_, err := i.env.EvalString(text)
	return err
}

func (i *Interpreter[T]) AddBasicCallable(name string, callable func(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error), helpText string) {
	i.env.AddFunction(name, callable)
	i.addHelpText(name, helpText)
}

func (i *Interpreter[T]) AddCallable(name string, callable func(env *zygo.Zlisp, name string, args []zygo.Sexp, ctx T) (zygo.Sexp, error), helpText string) {
	newHandler := func(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
		return callable(env, name, args, i.ctx)
	}
	i.env.AddFunction(name, newHandler)
	i.addHelpText(name, helpText)
}

func (i *Interpreter[T]) addHelpText(name string, text string) {
	i.manPages.addPage(name, text)
}

func (i *Interpreter[T]) REPL() {
	zygo.Repl(i.env, i.config)
}
