package validate

import (
	"fmt"
	"net/mail"
)

type ErrorHandler func(str string) error

func WithinLen(min, max int, name string) ErrorHandler {
	return func(str string) error {
		if len(str) >= min && len(str) <= max {
			return nil
		}
		return fmt.Errorf(
			"expected %s to be between %d and %d, but got %d",
			name,
			min,
			max,
			len(str),
		)
	}
}

func CcnValidator(cardNumber string) error {
	if len(cardNumber) < 13 {
		return fmt.Errorf("invalid credit card number")
	}

	total := 0
	isSecondDigit := false

	// Iterate through the card number digits in reverse order
	for i := len(cardNumber) - 1; i >= 0; i-- {
		// Skip spaces or any other non-digit characters
		if cardNumber[i] < '0' || cardNumber[i] > '9' {
			continue
		}

		// Convert the digit character to an integer
		digit := int(cardNumber[i] - '0')

		if isSecondDigit {
			// Double the digit for each second digit from the right
			digit *= 2
			if digit > 9 {
				// If doubling the digit results in a two-digit number, subtract 9 to get the sum of digits
				digit -= 9
			}
		}

		// Add the current digit to the total sum
		total += digit

		// Toggle the flag for the next iteration
		isSecondDigit = !isSecondDigit
	}

	// Check if the total sum mod 10 is not zero (valid card numbers should result in total mod 10 == 0)
	if total%10 != 0 {
		return fmt.Errorf("invalid credit card number")
	}
	return nil
}

func IsDigits(name string) ErrorHandler {
	return func(str string) error {
		for _, c := range str {
			if !(c >= '0' && c <= '9') {
				return fmt.Errorf("[%s] expected only digits but got %s", name, str)
			}
		}

		return nil
	}
}

func MustBeLen(length int, name string) ErrorHandler {
	return func(str string) error {
		if str == "" {
			return nil
		}
		if len(str) != length {
			return fmt.Errorf("Expected %s to be length %d but got %d", name, length, len(str))
		}
		return nil
	}
}

func NotEmpty(name string) ErrorHandler {
	return func(str string) error {
		if len(str) == 0 {
			return fmt.Errorf("%s cannot be empty", name)
		}
		return nil
	}
}

func EmailValidator(str string) error {
	_, err := mail.ParseAddress(str)
	if err != nil {
		return fmt.Errorf("not a valid email address")
	}
	return nil
}

func Compose(input ...ErrorHandler) ErrorHandler {
	return func(str string) error {
		for _, f := range input {
			err := f(str)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
