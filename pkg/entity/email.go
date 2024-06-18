package entity

import (
	"errors"
	"regexp"
)

// EmailAddress errors
var (
	ErrInvalidEmailAddress = errors.New("Not a valid email address")
)

// EmailAddress represents a valid email address
type EmailAddress struct {
	value string
}

// NewEmailAddress creates a new email address
func NewEmailAddress(email string) EmailAddress {
	var n EmailAddress
	match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email)
	if !match {
		panic("Email invalid")
		// return n, ErrInvalidEmailAddress
	}
	n.value = email

	// return n, nil
	return n
}

// String returns string representation of the email address
func (n EmailAddress) String() string {
	return n.value
}

// Equals checks that two email addresses are the same
// func (n EmailAddress) Equals(value Value) bool {
// 	otherEmailAddress, ok := value.(EmailAddress)
// 	return ok && n.value == otherEmailAddress.value
// }
