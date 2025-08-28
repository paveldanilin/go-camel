package camel

import (
	"errors"
	"regexp"
	"strings"
)

func ErrorEquals(str string) func(error) bool {

	return func(err error) bool {
		if err == nil {
			return false
		}

		return err.Error() == str
	}
}

func ErrorContains(substr string) func(error) bool {

	substrLower := strings.ToLower(substr)

	return func(err error) bool {
		if err == nil {
			return false
		}

		errorMsg := strings.ToLower(err.Error())

		return strings.Contains(errorMsg, substrLower)
	}
}

func ErrorIs(target string) func(error) bool {

	return func(err error) bool {
		return errors.Is(err, errors.New(target))
	}
}

func ErrorMatches(pattern string) func(error) bool {

	return func(err error) bool {
		if err == nil {
			return false
		}

		return regexp.MustCompile(pattern).MatchString(err.Error())
	}
}

func ErrorAny() func(error) bool {

	return func(err error) bool {
		return err != nil
	}
}
