package try

import (
	"errors"
	"regexp"
	"strings"
)

type ErrorMatcher func(error) bool

func ErrorEquals(str string) ErrorMatcher {
	return func(err error) bool {
		if err == nil {
			return false
		}

		return err.Error() == str
	}
}

func ErrorContains(str string) ErrorMatcher {
	substrLower := strings.ToLower(str)

	return func(err error) bool {
		if err == nil {
			return false
		}

		errorMsg := strings.ToLower(err.Error())

		return strings.Contains(errorMsg, substrLower)
	}
}

func ErrorIs(target string) ErrorMatcher {
	return func(err error) bool {
		return errors.Is(err, errors.New(target))
	}
}

func ErrorMatches(pattern string) ErrorMatcher {
	errRegex := regexp.MustCompile(pattern)

	return func(err error) bool {
		if err == nil {
			return false
		}
		return errRegex.MatchString(err.Error())
	}
}

func AnyError() ErrorMatcher {
	return func(err error) bool {
		return err != nil
	}
}
