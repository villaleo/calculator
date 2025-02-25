package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// CalculationRequest is a request object for receiving client requests for
// calculation requests.
type CalculationRequest struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// SumRequest is a request object for receiving client requests for summing a
// sequence of numbers.
type SumRequest struct {
	Numbers []float64 `json:"numbers"`
}

// CalculationResponse is a response object for sending responses back to the
// client for calculation requests.
type CalculationResponse struct {
	Interpretation string  `json:"interpretation"`
	Result         float64 `json:"result"`
}

// Interpret constructs a string interpretation from the CalculationRequest
// using the given operation. For example, if cr.X = 5.1 and cr.Y = 1, then
// "5.1 + 1" is returned. Unknown operations will result in an empty string.
func (cr CalculationRequest) Interpret(operation rune) string {
	x := strconv.FormatFloat(cr.X, 'f', -1, 64)
	y := strconv.FormatFloat(cr.Y, 'f', -1, 64)

	validOperations := []rune{'+', '-', '*', '/'}
	if !slices.Contains(validOperations, operation) {
		return ""
	}

	return fmt.Sprintf("%s %s %s", x, string(operation), y)
}

// Interpret constructs a string interpretation from the SumRequest,
// constructing an expression in the form of "x0 + x1 + ... + xN" where xI
// represents the the ith number, x, in sr.Numbers and N represents the amount
// of numbers in sr.Numbers.
//
// If sr.Numbers is empty (N = 0), then the interpretation will be an empty
// string. If a xI number is negative, then the appropriate operation symbols
// are displayed. For example, if sr.Numbers = [-1, 3.14, -8] then the
// interpretation would be "-1 + 3.14 - 8".
func (sr SumRequest) Interpret() string {
	len := len(sr.Numbers)
	if len == 1 {
		return strconv.FormatFloat(sr.Numbers[0], 'f', -1, 64)
	}

	var result strings.Builder

	for i, num := range sr.Numbers {
		formattedNum := strconv.FormatFloat(num, 'f', -1, 64)

		var sep string
		if i+1 != len {
			// Separate numbers by either a '+', or '-' depending on the value of the
			// next number.
			if sr.Numbers[i+1] < 0 {
				sep = " - "
			} else {
				sep = " + "
			}
		}

		// Remove the '-' from the formatted number since the seperator is already
		// set to '+' or '-'. Skip the first number since there are no numbers before
		// it.
		if i != 0 {
			formattedNum = strings.TrimLeft(formattedNum, "-")
		}
		result.WriteString(formattedNum + sep)
	}

	return result.String()
}
