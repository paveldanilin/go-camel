package errmatcher

type Mode string

const (
	ModeEquals   Mode = "equals"
	ModeContains Mode = "contains"
	ModeRegex    Mode = "regex"
	ModeIs       Mode = "is"
)

// Definition represents error matcher definition.
//
//	Allowed matchers: equals, is, contains, regex.
type Definition struct {
	MatchMode Mode
	Target    string
}

func Is(target error) Definition {
	return Definition{
		MatchMode: ModeIs,
		Target:    target.Error(),
	}
}

func Equals(str string) Definition {
	return Definition{
		MatchMode: ModeEquals,
		Target:    str,
	}
}

func Any() Definition {
	return Definition{
		MatchMode: ModeEquals,
		Target:    "*",
	}
}

func Contains(str string) Definition {
	return Definition{
		MatchMode: ModeContains,
		Target:    str,
	}
}

func Matches(pattern string) Definition {
	return Definition{
		MatchMode: ModeRegex,
		Target:    pattern,
	}
}
