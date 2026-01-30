package functional

func MakeNoArgs[ArgVal any, RetVal any](f func(ArgVal) RetVal, arg ArgVal) func() RetVal {
	return func() RetVal {
		return f(arg)
	}
}

func MakeSingleArg[Arg0Val any, Arg1Value, RetVal any](f func(Arg0Val, Arg1Value) RetVal, a0 Arg0Val) func(Arg1Value) RetVal {
	return func(arg1 Arg1Value) RetVal {
		return f(a0, arg1)
	}
}
