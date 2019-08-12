package runner

import "flag"

type Runner interface {
	Run() error
}

func parseConfig() (config, error) {
	flag.Parse()
	return config{
		drivers: flag.Args(),
	}, nil
}

type config struct {
	drivers []string
}

type Continue struct {
	cnf config
}

type Help struct {
	err error
}

func (h *Help) Run() error {
	flag.Usage()
	return h.err
}
