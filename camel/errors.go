package camel

import (
	"errors"
	"regexp"
	"strings"
)

func ErrorEquals(err error, str string) bool {

	if err == nil {
		return false
	}

	return err.Error() == str
}

func PredicateErrorEquals(str string) func(error) bool {

	return func(err error) bool {
		return ErrorEquals(err, str)
	}
}

func ErrorContains(err error, substr string) bool {

	if err == nil {
		return false
	}

	errorMsg := strings.ToLower(err.Error())
	substrLower := strings.ToLower(substr)

	return strings.Contains(errorMsg, substrLower)
}

func PredicateErrorContains(substr string) func(error) bool {

	return func(err error) bool {
		return ErrorContains(err, substr)
	}
}

func ErrorIs(err error, target string) bool {

	return errors.Is(err, errors.New(target))
}

func PredicateErrorIs(target string) func(error) bool {

	return func(err error) bool {
		return ErrorIs(err, target)
	}
}

func ErrorMatches(err error, pattern string) bool {

	if err == nil {
		return false
	}

	return regexp.MustCompile(pattern).MatchString(err.Error())
}

func PredicateErrorMatches(pattern string) func(error) bool {

	return func(err error) bool {
		return ErrorMatches(err, pattern)
	}
}
