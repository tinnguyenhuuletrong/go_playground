package module_opt_play

type opt struct {
	address string
	port    uint
	ssl     bool
}

type moduleX struct {
	opt
}

func makeDefaultOpt() *opt {
	return &opt{
		address: "localhost",
		port:    8080,
		ssl:     false,
	}
}

type OptModify func(*opt) *opt

func WithAddress(address string) OptModify {
	return func(o *opt) *opt {
		o.address = address
		return o
	}
}

func WithPort(port uint) OptModify {
	return func(o *opt) *opt {
		o.port = port
		return o
	}
}

func WithSSL(ssl bool) OptModify {
	return func(o *opt) *opt {
		o.ssl = ssl
		return o
	}
}

func NewModuleX(opts ...OptModify) *moduleX {
	opt := makeDefaultOpt()
	for _, v := range opts {
		opt = v(opt)
	}
	return &moduleX{
		*opt,
	}
}
