package tax

import (
	"fmt"
	"testing"

	"github.com/vfc2/tax-calculator/internal/money"
)

func TestNationalInsurance(t *testing.T) {
	tests := map[string]struct {
		wage     int64
		expected string
	}{
		"NoTax": {
			wage:     120,
			expected: "0.00",
		},
		"Mid": {
			wage:     731,
			expected: "48.90",
		},
		"High": {
			wage:     1058,
			expected: "74.32",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := money.New(test.wage)
			actual := CalculateNationalInsurance(v).Format(2)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestIncomeTax(t *testing.T) {
	tests := map[string]struct {
		wage      int64
		allowance int64
		expected  IncomeTaxBreakdown
	}{
		"NoTax": {
			wage:      7543,
			allowance: 12570,
			expected:  IncomeTaxBreakdown{},
		},
		"BasicRate": {
			wage:      35000,
			allowance: 12570,
			expected: IncomeTaxBreakdown{
				BasicRate: money.New(4486),
				Taxable:   money.New(22430),
				Taxed:     money.New(4486),
			},
		},
		"HigherRate": {
			wage:      63450,
			allowance: 12570,
			expected: IncomeTaxBreakdown{
				BasicRate:  money.New(7540.20),
				HigherRate: money.New(5271.60),
				Taxable:    money.New(50880),
				Taxed:      money.New(12811.80),
			},
		},
		"AdditionalRate": {
			wage:      143000,
			allowance: 0,
			expected: IncomeTaxBreakdown{
				BasicRate:      money.New(7540.20),
				HigherRate:     money.New(34976),
				AdditionalRate: money.New(8036.55),
				Taxable:        money.New(143000),
				Taxed:          money.New(50552.75),
			},
		},
	}

	taxRates := LoadTaxBands()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := money.New(test.wage)
			a := money.New(test.allowance)

			actual := CalculateIncomeTax(w, a, taxRates)

			fmt.Println(actual)
			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestTaxAllowance(t *testing.T) {
	tests := map[string]struct {
		wage     int64
		expected string
	}{
		"NoAllowance": {
			wage:     145000,
			expected: "0.00",
		},
		"FullAllowance": {
			wage:     65000,
			expected: "12570.00",
		},
		"PartialAllowance": {
			wage:     112000,
			expected: "6570.00",
		},
	}

	taxRates := LoadTaxBands()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := money.New(test.wage)

			actual := CalculateTaxAllowance(v, taxRates).Format(2)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestCalculateAllowanceAndTax(t *testing.T) {
	tests := map[string]struct {
		wage     int64
		expected IncomeTaxBreakdown
	}{
		"NoTax": {
			wage:     7543,
			expected: IncomeTaxBreakdown{},
		},
		"BasicRate": {
			wage: 35000,
			expected: IncomeTaxBreakdown{
				BasicRate: money.New(4486),
				Taxable:   money.New(22430),
				Taxed:     money.New(4486),
			},
		},
		"HigherRate": {
			wage: 63450,
			expected: IncomeTaxBreakdown{
				BasicRate:  money.New(7540.20),
				HigherRate: money.New(5271.60),
				Taxable:    money.New(50880),
				Taxed:      money.New(12811.80),
			},
		},
		"AdditionalRate": {
			wage: 143000,
			expected: IncomeTaxBreakdown{
				BasicRate:      money.New(7540.20),
				HigherRate:     money.New(34976),
				AdditionalRate: money.New(8036.55),
				Taxable:        money.New(143000),
				Taxed:          money.New(50552.75),
			},
		},
	}

	taxRates := LoadTaxBands()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := money.New(test.wage)

			a := CalculateTaxAllowance(w, taxRates)
			actual := CalculateIncomeTax(w, a, taxRates)

			fmt.Println(actual)
			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}
