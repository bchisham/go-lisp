package builtins

import (
	"io"
	"os"

	"github.com/bchisham/go-lisp/scheme/internal/pkg/parser/values"
)

type EvaluatorCallback Expression

type Runtime struct {
	Out io.Writer
	Err io.Writer
	Env Environment
}

type configRuntime struct {
	out      io.Writer
	err      io.Writer
	env      Environment
	callback Expression
}

type OptionRuntime func(*configRuntime)

func WithEnv(env Environment) OptionRuntime {
	return func(c *configRuntime) {
		c.env = env
	}
}

func WithOut(out io.Writer) OptionRuntime {
	return func(c *configRuntime) {
		c.out = out
	}
}

func WithErr(err io.Writer) OptionRuntime {
	return func(c *configRuntime) {
		c.err = err
	}
}

func WithEvaluatorCallback(cb Expression) OptionRuntime {
	return func(c *configRuntime) {
		c.callback = cb
	}
}

func defaultConfig() configRuntime {
	return configRuntime{
		out: os.Stdout,
		err: os.Stderr,
		env: NewEnvironment(),
		callback: func(v values.Interface, runtime *Runtime) (values.Interface, error) {
			return v, nil
		},
	}
}

func NewRuntime(opt ...OptionRuntime) *Runtime {
	cfg := defaultConfig()
	for _, o := range opt {
		o(&cfg)
	}
	rt := &Runtime{
		Out: cfg.out,
		Err: cfg.err,
		Env: cfg.env,
	}
	rt.defaultEnvironment(cfg.callback)
	return rt
}
