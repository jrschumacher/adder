package adder

import (
	"fmt"
	"strings"
)

// Request interface that all generated request structs implement
// This ensures adder package is always imported, eliminating conditional imports
type Request interface {
	GetRawArguments() []string
}

// ValidateEnum validates that a value is in the allowed enum list
// Returns nil if valid, error with helpful message if invalid
func ValidateEnum(flagName, value string, validValues []string) error {
	for _, validValue := range validValues {
		if value == validValue {
			return nil
		}
	}

	var enumList string
	if len(validValues) <= 1 {
		enumList = strings.Join(validValues, "")
	} else if len(validValues) == 2 {
		enumList = validValues[0] + " or " + validValues[1]
	} else {
		enumList = strings.Join(validValues[:len(validValues)-1], ", ") + ", or " + validValues[len(validValues)-1]
	}

	return fmt.Errorf("invalid %s: %s (must be %s)", flagName, value, enumList)
}