package service

import (
	"errors"
	"main/internal/models"
	"regexp"
)

const (
	pairRegex     = `^[\d\w]+([\-\/\_]{1})?[A-Za-z]+$`
	exchangeRegex = `^(binance_spot|binance_futures|binance_us|bybit_spot|bybit_futures)$`
	emailRegex    = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z])?)*$"
	directoryPath = "internal.service."
)

var (
	errGettingFoundVolume        = errors.New("error getting found volumes")
	errEmailIsEmpty              = errors.New("email data is empty")
	errPairNameIsEmpty           = errors.New("pair name is empty")
	errExchangeNameIsEmpty       = errors.New("exchange name is empty")
	errPasswordIsEmpty           = errors.New("user password value is empty")
	errEmailInvalidFormat        = errors.New("invalid email format")
	errPairNameInvalidFormat     = errors.New("invalid pair name format")
	errExchangeNameInvalidFormat = errors.New("invalid exchange name format")
	errIdBelowOne                = errors.New("user id must be above zero")
	errExactValueBelowZero       = errors.New("exact value must be above zero")
)

// CheckUserData validates the user data before operations like signing up and logging in.
// It performs the following checks:
//   - the Email field is not empty
//   - the Password field is not empty
//   - the email format matches the predefined regex pattern
//
// If any of these checks fail, an error is returned indicating the specific problem.
// If all checks pass, nil is returned indicating that the user data is valid.
func CheckUserData(user models.User) error {
	// Check if the Email field of the user struct is empty
	if user.Email == "" {
		// Return an error indicating that the email data must be provided
		return errEmailIsEmpty
	}

	// Check if the Password field length is zero
	if len(user.Password) == 0 {
		// Return an error indicating that the password must be provided
		return errPasswordIsEmpty
	}

	// Use a regular expression to validate the format of the email against a predefined regex pattern
	isMatch, err := regexp.MatchString(emailRegex, user.Email)
	if err != nil || !isMatch {
		// If there was an error during regex matching or if the email does not match the expected format,
		// return an error indicating that the email format is invalid
		return errEmailInvalidFormat
	}

	// If all checks pass without errors, return nil indicating that the user data is valid
	return nil
}

// CheckPairData checks if the provided pairData satisfies the following criteria:
//   - the Pair field is not empty
//   - the Exchange field is not empty
//   - the ExactValue is greater than or equal to 1
//   - the UserID is greater than 0
//   - the pair name matches a predefined regex pattern
//   - the exchange name matches a predefined regex pattern
//
// If any of these checks fail, an error is returned indicating the specific problem.
// If all checks pass, nil is returned indicating that the pairData is valid.
func CheckPairData(pairData models.UserPairs) error {
	// Check if the Pair field of the pairData struct is empty
	if pairData.Pair == "" {
		// Return an error indicating that the pair name must be provided
		return errPairNameIsEmpty
	}

	// Check if the Exchange field of the pairData struct is empty
	if pairData.Exchange == "" {
		// Return an error indicating that the exchange name must be provided
		return errExchangeNameIsEmpty
	}

	// Check if ExactValue is less than 1
	if pairData.ExactValue < 1 {
		// Return an error indicating that the exact value must be above zero
		return errExactValueBelowZero
	}

	// Check if UserID is less than 1
	if pairData.UserID < 1 {
		// Return an error indicating that a valid user ID must be provided
		return errIdBelowOne
	}

	// Use a regular expression to validate the format of the trading pair name against a predefined pattern
	isMatch, err := regexp.MatchString(pairRegex, pairData.Pair)
	if err != nil || !isMatch {
		// If there was an error during regex matching or if the pair name does not match the expected format,
		// return an error indicating that the pair name format is invalid
		return errPairNameInvalidFormat
	}

	// Use a regular expression to validate the format of the exchange name against a predefined pattern
	isMatch, err = regexp.MatchString(exchangeRegex, pairData.Exchange)
	if err != nil || !isMatch {
		// If there was an error during regex matching or if the exchange name does not match the expected format,
		// return an error indicating that the exchange name format is invalid
		return errExchangeNameInvalidFormat
	}

	// If all checks pass without errors, return nil indicating that the trading pair data is valid
	return nil
}
