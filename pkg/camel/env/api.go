package env

type Env interface {
	LookupVar(name string) (string, bool)
}
