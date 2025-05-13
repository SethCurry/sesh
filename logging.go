package sesh

import (
	"fmt"

	"github.com/glycerine/zygomys/v9/zygo"
	"github.com/rs/zerolog"
)

func NewLogging(logger *zerolog.Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

type Logging struct {
	logger *zerolog.Logger
}

func (l *Logging) BasicCallables() []BaseCallable {
	return []BaseCallable{
		{
			Name:    "log_debug",
			Handler: l.Debug,
		},
	}
}

func (l *Logging) Register(env *zygo.Zlisp) {
	env.AddFunction("log_debug", l.Debug)
	env.AddFunction("log_info", l.Info)
	env.AddFunction("log_warn", l.Warn)
	env.AddFunction("log_error", l.Error)
}

func (l *Logging) log(evt *zerolog.Event, args []zygo.Sexp) (zygo.Sexp, error) {
	var keySexp *zygo.SexpStr
	key := ""
	ok := false
	for i, v := range args {
		if i == len(args)-1 {
			keySexp, ok = v.(*zygo.SexpStr)
			if !ok {
				return zygo.SexpNull, fmt.Errorf("expected string message, got %T", v)
			}
			evt.Msg(keySexp.S)
			break
		}

		if i%2 == 0 {
			keySexp, ok = v.(*zygo.SexpStr)
			if !ok {
				return zygo.SexpNull, fmt.Errorf("expected string key, got %T", v)
			}
			key = keySexp.S
		} else {
			switch v.Type() {
			case zygo.GoStructRegistry.Registry["string"]:
				evt.Str(key, v.(*zygo.SexpStr).S)
			case zygo.GoStructRegistry.Registry["int64"]:
				evt.Int64(key, v.(*zygo.SexpInt).Val)
			case zygo.GoStructRegistry.Registry["float64"]:
				evt.Float64(key, v.(*zygo.SexpFloat).Val)
			default:
				return zygo.SexpNull, fmt.Errorf("unsupported type: %T", v)
			}
		}
	}

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
