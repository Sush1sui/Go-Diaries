package helper

import (
	"reflect"
	"strings"
)

func ValidateUserInput(firstName, lastName, email string, userTickets, remainingTickets uint) (bool, bool, bool) {
	isValidName := len(firstName) > 2 && len(lastName) > 2
	isValidEmail := strings.Contains(email, "@") && len(email) > 5
	isValidTicketNumber := reflect.TypeOf(userTickets).Kind() == reflect.Uint && userTickets > 0 && userTickets <= remainingTickets
	return isValidName, isValidEmail, isValidTicketNumber
}