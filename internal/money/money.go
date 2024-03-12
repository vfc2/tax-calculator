// Package money implements utility routines for manipulating monetary values.
//
// This package is a basic implementation of Google's micros
// https://developers.google.com/standard-payments/reference/glossary#micros
// Storing monetary values as integer by multiplying / dividing it by 1,000,000.
// For example £1.23 is stored as 1230000.
package money

import (
	"math"
	"strconv"
)

const unit = 1000000

// Store monetary values as integer.
type Money int64

// New initializes and return a Money. It accepts integer and
// float types. For float types, he input will be rounded to
// the nearest integer using math.RoundToEven().
func New[T int | int32 | int64 | float64](amount T) Money {
	var m Money = 0

	switch any(amount).(type) {
	case int, int32, int64:
		m = Money(amount * unit)
	case float64:
		m = Money(math.RoundToEven(float64(amount) * unit))
	}

	return m
}

// NewForString initializes and return a Money. It accepts an arbitrary
// string containing a number.
// The number is converted to float64 then rounded to the nearest
// integer using math.RoundToEven().
// If the conversion fails, an error is returned.
func NewFromString(amount string) (Money, error) {
	v, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	return New(v), nil
}

// Mul returns a Money multiplied by a provided float64.
func (m Money) Mul(mul float64) Money {
	factor := math.RoundToEven(mul * unit)

	return m * Money(factor) / unit
}

// Div returns a Money divided by a provided float64.
// If the divisor is 0, the Money will be returned as-is.
func (m Money) Div(div float64) Money {
	if div == 0 {
		return m
	}

	factor := math.RoundToEven(div * unit)
	v := math.RoundToEven(float64(m) / factor * unit)

	return Money(v)
}

// Sub returns a Money substracted by a provided float64.
func (m Money) Sub(sub float64) Money {
	v := math.RoundToEven(sub * unit)

	return m - Money(v)
}

// Format returns a Money as a string with the specified digits.
// For example 66498000 with digits = 2 is returned as 66.50.
func (m Money) Format(digits int) string {
	return strconv.FormatFloat(float64(m)/unit, 'f', digits, 64)
}
