package env

import "os"

type OsEnv struct {
}

func NewOsEnv() *OsEnv {
	return &OsEnv{}
}

func (env *OsEnv) LookupVar(name string) (string, bool) {
	return os.LookupEnv(name)
}
