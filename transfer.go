package quark

func TransferAR(fn func()) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		fn()
		return Null, nil
	}
}

func TransferARI(fn func() int64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		return FromInterface(fn())
	}
}

func TransferAIRI(fn func(int64) int64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		a0, err := ToInt(args[0])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(a0))
	}
}

func TransferAIIRI(fn func(int64, int64) int64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		a0, err := ToInt(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := ToInt(args[1])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(a0, a1))
	}
}

func TransferAIIIRI(fn func(int64, int64, int64) int64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		a0, err := ToInt(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := ToInt(args[1])
		if err != nil {
			return nil, err
		}
		a2, err := ToInt(args[2])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(a0, a1, a2))
	}
}

func TransferARF(fn func() float64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		return FromInterface(fn())
	}
}

func TransferAFRF(fn func(float64) float64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		f0, err := ToFloat(args[0])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(f0))
	}
}

func TransferAFFRF(fn func(float64, float64) float64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		f0, err := ToFloat(args[0])
		if err != nil {
			return nil, err
		}
		f1, err := ToFloat(args[1])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(f0, f1))
	}
}

func TransferAFFFRF(fn func(float64, float64, float64) float64) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		f0, err := ToFloat(args[0])
		if err != nil {
			return nil, err
		}
		f1, err := ToFloat(args[1])
		if err != nil {
			return nil, err
		}
		f2, err := ToFloat(args[2])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(f0, f1, f2))
	}
}

func TransferARS(fn func() string) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		return FromInterface(fn())
	}
}

func TransferASRS(fn func(string) string) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		a0, err := ToString(args[0])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(a0))
	}
}

func TransferASSRS(fn func(string, string) string) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		a0, err := ToString(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := ToString(args[1])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(a0, a1))
	}
}

func TransferASSSRS(fn func(string, string, string) string) CallableFunction {
	return func(ctx *Context, args []Object) (Object, error) {
		a0, err := ToString(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := ToString(args[1])
		if err != nil {
			return nil, err
		}
		a2, err := ToString(args[2])
		if err != nil {
			return nil, err
		}
		return FromInterface(fn(a0, a1, a2))
	}
}
