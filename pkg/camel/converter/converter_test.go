package converter

import (
	"testing"
)

func TestStringToDateTime(t *testing.T) {
	c := StringToDateTime()

	time, err := c.Convert("2010-05-05 05:05:05", nil)

	if err != nil {
		t.Fatalf("TestStringToDateTime() = %s", err)
	}

	if time.Year() != 2010 {
		t.Fatalf("TestStringToDateTime() = %d; want = %d", time.Year(), 2010)
	}

	if time.Month() != 5 {
		t.Fatalf("TestStringToDateTime() = %d; want = %d", time.Month(), 5)
	}

	if time.Day() != 5 {
		t.Fatalf("TestStringToDateTime() = %d; want = %d", time.Day(), 5)
	}
}
