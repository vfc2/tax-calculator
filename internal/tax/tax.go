package tax

import (
	"github.com/vfc2/tax-calculator/internal/money"
)

type Money = money.Money

type TaxBand struct {
	Min  Money
	Max  Money
	Rate float64
}

type IncomeTaxRates struct {
	PersonalAllowance          Money
	PersonalAllowanceThreshold Money
	Basic                      TaxBand
	Higher                     TaxBand
	Additional                 TaxBand
}

type IncomeTaxBreakdown struct {
	BasicRate      Money
	HigherRate     Money
	AdditionalRate Money
	Taxable        Money
	Taxed          Money
}

func LoadTaxBands() IncomeTaxRates {
	return IncomeTaxRates{
		PersonalAllowance:          money.New(12570),
		PersonalAllowanceThreshold: money.New(100000),
		Basic: TaxBand{
			Min:  money.New(12571),
			Max:  money.New(50270),
			Rate: 0.2,
		},
		Higher: TaxBand{
			Min:  money.New(50271),
			Max:  money.New(125140),
			Rate: 0.4,
		},
		Additional: TaxBand{
			Min:  money.New(125141),
			Max:  money.New(0),
			Rate: 0.45,
		},
	}
}

// Calculate the National Insurance amount due weekly for Category A.
// Requirements from https://www.gov.uk/national-insurance-rates-letters
func CalculateNationalInsurance(weekIncome Money) Money {
	c := max(weekIncome.Sub(967), 0)
	b := max((weekIncome - c).Sub(242), 0)

	tax := c.Mul(0.02) + b.Mul(0.1)

	return tax
}

// Calculate the Taxable Income of yearly gross income.
// Requirements from https://www.gov.uk/income-tax-rates
func CalculateIncomeTax(income Money, allowance Money, tr IncomeTaxRates) IncomeTaxBreakdown {
	hrLimit := tr.Higher.Min + (allowance - tr.PersonalAllowance)

	ar := max(income-tr.Additional.Min, 0)
	hr := max(income-ar-hrLimit, 0)
	br := max(income-ar-hr-allowance, 0)

	tax := ar.Mul(tr.Additional.Rate) + hr.Mul(tr.Higher.Rate) + br.Mul(tr.Basic.Rate)

	return IncomeTaxBreakdown{
		BasicRate:      br.Mul(tr.Basic.Rate),
		HigherRate:     hr.Mul(tr.Higher.Rate),
		AdditionalRate: ar.Mul(tr.Additional.Rate),
		Taxed:          tax,
		Taxable:        max(income-allowance, 0),
	}
}

// Calculate the Tax Allowance based on a yearly gross income.
// Requirements from https://www.gov.uk/income-tax-rates/income-over-100000
func CalculateTaxAllowance(annumIncome Money, tr IncomeTaxRates) Money {
	over := max((annumIncome - tr.PersonalAllowanceThreshold).Mul(0.5), 0)

	return max(tr.PersonalAllowance-over, 0)
}
