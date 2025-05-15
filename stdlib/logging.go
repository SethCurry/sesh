package stdlib

import (
	"fmt"
	"os"

	"github.com/SethCurry/sesh"
	"github.com/glycerine/zygomys/v9/zygo"
	"github.com/rs/zerolog"
)

func NewLogging(logger *zerolog.Logger) *Logging {
	if logger == nil {
		defaultLogger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
		logger = &defaultLogger
	}
	return &Logging{
		logger: logger,
	}
}

type Logging struct {
	logger *zerolog.Logger
}

func (l *Logging) BasicCallables() []sesh.BaseCallable {
	return []sesh.BaseCallable{
		{
			Name:    "log_debug",
			Handler: l.Debug,
		},
		{
			Name:    "log_info",
			Handler: l.Info,
		},
		{
			Name:    "log_warn",
			Handler: l.Warn,
		},
		{
			Name:    "log_error",
			Handler: l.Error,
		},
	}
}

func (l *Logging) log(evt *zerolog.Event, args []zygo.Sexp) (zygo.Sexp, error) {
	if len(args) == 0 {
		return zygo.SexpNull, fmt.Errorf("expected at least one argument")
	}

	if len(args) > 2 {
		return zygo.SexpNull, fmt.Errorf("expected at most two arguments")
	}

	messageSexp, ok := args[0].(*zygo.SexpStr)
	if !ok {
		return zygo.SexpNull, fmt.Errorf("expected string message, got %T", args[0])
	}

	if len(args) == 2 {
		hashSexp, ok := args[1].(*zygo.SexpHash)
		if !ok {
			return zygo.SexpNull, fmt.Errorf("expected hash, got %T", args[1])
		}

		for _, pair := range hashSexp.Map {
			for _, v := range pair {
				keySexp, ok := v.Head.(*zygo.SexpSymbol)
				if !ok {
					return zygo.SexpNull, fmt.Errorf("expected string key, got %T", v.Head)
				}
				key := keySexp.Name()

				switch v.Tail.Type() {
				case zygo.GoStructRegistry.Registry["string"]:
					evt.Str(key, v.Tail.(*zygo.SexpStr).S)
				case zygo.GoStructRegistry.Registry["int64"]:
					evt.Int64(key, v.Tail.(*zygo.SexpInt).Val)
				case zygo.GoStructRegistry.Registry["float64"]:
					evt.Float64(key, v.Tail.(*zygo.SexpFloat).Val)
				default:
					return zygo.SexpNull, fmt.Errorf("unsupported type: %T", v)
				}
			}
		}
	}

	evt.Msg(messageSexp.S)

	return zygo.SexpNull, nil
}

func (l *Logging) Debug(_ *zygo.Zlisp, _ string, args []zygo.Sexp) (zygo.Sexp, error) {
	return l.log(l.logger.Debug(), args)
}

func (l *Logging) Info(_ *zygo.Zlisp, _ string, args []zygo.Sexp) (zygo.Sexp, error) {
	return l.log(l.logger.Info(), args)
}

func (l *Logging) Warn(_ *zygo.Zlisp, _ string, args []zygo.Sexp) (zygo.Sexp, error) {
	return l.log(l.logger.Warn(), args)
}

func (l *Logging) Error(_ *zygo.Zlisp, _ string, args []zygo.Sexp) (zygo.Sexp, error) {
	return l.log(l.logger.Error(), args)
}
