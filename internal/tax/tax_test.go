package tax

import (
	"testing"

	"github.com/vfc2/tax-calculator/internal/money"
)

var taxRates = IncomeTaxRates{
	PersonalAllowance:          money.New(12570),
	PersonalAllowanceThreshold: money.New(100000),
	Basic: Band{
		Min:  money.New(12571),
		Max:  money.New(50270),
		Rate: 0.2,
	},
	Higher: Band{
		Min:  money.New(50271),
		Max:  money.New(125140),
		Rate: 0.4,
	},
	Additional: Band{
		Min:  money.New(125141),
		Max:  money.New(0),
		Rate: 0.45,
	},
}

var niRates = map[string]NationalInsuranceRates{
	"A": {
		Band1: Band{
			Min:  123000000,
			Max:  242000000,
			Rate: 0,
		},
		Band2: Band{
			Min:  242010000,
			Max:  967000000,
			Rate: 0.1,
		},
		Band3: Band{
			Min:  967010000,
			Rate: 0.02,
		},
	},
}

func TestNationalInsurance(t *testing.T) {
	tests := map[string]struct {
		income   Money
		expected Money
	}{
		"NoTax": {
			income:   money.New(120),
			expected: 0,
		},
		"Mid": {
			income:   money.New(731),
			expected: money.New(48.9),
		},
		"High": {
			income:   money.New(1058),
			expected: money.New(74.32),
		},
	}

	tests_fail := map[string]struct {
		category string
	}{
		"DoesntExist": {
			category: "ZZ",
		},
		"Empty": {
			category: "",
		},
	}

	tax := TaxCalculator{
		IncomeTaxRates:         taxRates,
		NationalInsuranceRates: niRates,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, _ := tax.calculateNationalInsurance(test.income, "A")

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}

	for name, test := range tests_fail {
		t.Run(name, func(t *testing.T) {
			actual, err := tax.calculateNationalInsurance(0, test.category)

			if actual != 0 || err == nil {
				t.Error("an error was expected but not returned")
			}
		})
	}
}

func TestIncomeTax(t *testing.T) {
	tests := map[string]struct {
		income    Money
		allowance Money
		expected  IncomeTaxBreakdown
	}{
		"NoTax": {
			income:    money.New(7543),
			allowance: money.New(12570),
			expected: IncomeTaxBreakdown{
				GrossIncome: money.New(7543),
			},
		},
		"BasicRate": {
			income:    money.New(35000),
			allowance: money.New(12570),
			expected: IncomeTaxBreakdown{
				GrossIncome: money.New(35000),
				BasicRate:   money.New(4486),
				Taxable:     money.New(22430),
				Taxed:       money.New(4486),
			},
		},
		"HigherRate": {
			income:    money.New(63450),
			allowance: money.New(12570),
			expected: IncomeTaxBreakdown{
				GrossIncome: money.New(63450),
				BasicRate:   money.New(7540.20),
				HigherRate:  money.New(5271.60),
				Taxable:     money.New(50880),
				Taxed:       money.New(12811.80),
			},
		},
		"AdditionalRate": {
			income:    money.New(143000),
			allowance: 0,
			expected: IncomeTaxBreakdown{
				GrossIncome:    money.New(143000),
				BasicRate:      money.New(7540.20),
				HigherRate:     money.New(34976),
				AdditionalRate: money.New(8036.55),
				Taxable:        money.New(143000),
				Taxed:          money.New(50552.75),
			},
		},
	}

	tax := TaxCalculator{
		IncomeTaxRates:         taxRates,
		NationalInsuranceRates: niRates,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tax.calculateIncomeTax(test.income, test.allowance)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestTaxAllowance(t *testing.T) {
	tests := map[string]struct {
		income   Money
		expected Money
	}{
		"NoAllowance": {
			income:   money.New(145000),
			expected: money.New(0),
		},
		"FullAllowance": {
			income:   money.New(65000),
			expected: money.New(12570),
		},
		"PartialAllowance": {
			income:   money.New(112000),
			expected: money.New(6570),
		},
	}

	tax := TaxCalculator{
		IncomeTaxRates:         taxRates,
		NationalInsuranceRates: niRates,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tax.calculateTaxAllowance(test.income)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestTakeHome(t *testing.T) {
	tests := map[string]struct {
		income           Money
		expectedTakeHome string
		expectedNI       string
	}{
		"NoTax": {
			income:           money.New(7543),
			expectedTakeHome: "7543.00",
			expectedNI:       "0.00",
		},
		"HigherRate": {
			income:           money.New(63450),
			expectedTakeHome: "46604.88",
			expectedNI:       "4033.32",
		},
	}

	tests_fail := map[string]struct {
		niCategory string
	}{
		"NIDoesntExist": {
			niCategory: "ZZ",
		},
		"NIEmpty": {
			niCategory: "",
		},
	}

	tax := TaxCalculator{
		IncomeTaxRates:         taxRates,
		NationalInsuranceRates: niRates,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, _ := tax.CalculateTakeHome(test.income, "A")

			takeHome := actual.TakeHome.Format(2)
			ni := actual.NationalInsurance.Format(2)

			if takeHome != test.expectedTakeHome || ni != test.expectedNI {
				t.Errorf("got {TakeHome: %s, National Insurance: %s}, want {TakeHome: %s, National Insurance: %s}",
					takeHome, ni, test.expectedTakeHome, test.expectedNI)
			}
		})
	}

	for name, test := range tests_fail {
		t.Run(name, func(t *testing.T) {
			actual, err := tax.CalculateTakeHome(0, test.niCategory)
			expected := IncomeTaxBreakdown{}

			if actual != expected || err == nil {
				t.Error("an error was expected but not returned")
			}
		})
	}
}
