package env

import "os"

type Os struct {
}

func NewOs() *Os {
	return &Os{}
}

func (env *Os) LookupVar(name string) (string, bool) {
	return os.LookupEnv(name)
}
