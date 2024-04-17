package pages

import (
	"fmt"
	"strconv"
	"strings"
)

type errorHandler func(str string) error

func withinLen(min, max int, name string) errorHandler {
    return func(str string) error {
        if len(str) < min && len(str) >= max {
            return fmt.Errorf("expected %s to be between %d and %d, but got %d", name, min, max, len(str))
        }
        return nil
    }
}

func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func isDigits(name string) errorHandler {
    return func(str string) error {

        for _, c := range str {
            if !(c >= '0' && c <= '9') {
                return fmt.Errorf("expected only digits but got %s", str)
            }
        }

        return nil
    }
}


func mustBeLen(length int, name string) errorHandler {
    return func(str string) error {
        if len(str) != length {
            return fmt.Errorf("Expected %s to be length %d but got %d", name, length, len(str))
        }
        return nil
    }
}

func notEmpty(name string) errorHandler {
    return func(str string) error {
        if len(str) == 0 {
            return fmt.Errorf("%s cannot be empty", name)
        }
        return nil
    }
}

func compose(a, b errorHandler) errorHandler {
    return func(str string) error {
        err := a(str)
        if err != nil {
            return err
        }
        return b(str)
    }
}


