package errs

type MatchMode string

const (
	MatchModeEquals   MatchMode = "equals"
	MatchModeContains MatchMode = "contains"
	MatchModeRegex    MatchMode = "regex"
	MatchModeIs       MatchMode = "is"
)

type Matcher struct {
	MatchMode MatchMode
	Target    string
}

func Is(target error) Matcher {
	return Matcher{
		MatchMode: MatchModeIs,
		Target:    target.Error(),
	}
}

func Equals(str string) Matcher {
	return Matcher{
		MatchMode: MatchModeEquals,
		Target:    str,
	}
}

func Any() Matcher {
	return Matcher{
		MatchMode: MatchModeEquals,
		Target:    "*",
	}
}

func Contains(str string) Matcher {
	return Matcher{
		MatchMode: MatchModeContains,
		Target:    str,
	}
}

func Matches(pattern string) Matcher {
	return Matcher{
		MatchMode: MatchModeRegex,
		Target:    pattern,
	}
}
