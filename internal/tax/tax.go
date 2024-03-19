package tax

import (
	"fmt"

	"github.com/vfc2/tax-calculator/internal/money"
)

type Money = money.Money

type Band struct {
	Min  Money
	Max  Money
	Rate float64
}

type IncomeTaxRates struct {
	PersonalAllowance          Money
	PersonalAllowanceThreshold Money
	Basic                      Band
	Higher                     Band
	Additional                 Band
}

type NationalInsuranceRates struct {
	Band1 Band
	Band2 Band
	Band3 Band
}

type IncomeTaxBreakdown struct {
	GrossIncome       Money
	BasicRate         Money
	HigherRate        Money
	AdditionalRate    Money
	Taxable           Money
	Taxed             Money
	NationalInsurance Money
	TakeHome          Money
}

type TaxCalculator struct {
	IncomeTaxRates         IncomeTaxRates
	NationalInsuranceRates map[string]NationalInsuranceRates
}

// Calculate the National Insurance amount due weekly for Category A.
// Requirements from https://www.gov.uk/national-insurance-rates-letters
func (t TaxCalculator) calculateNationalInsurance(weekIncome Money, category string) (Money, error) {
	cat, ok := t.NationalInsuranceRates[category]
	if !ok {
		return 0, fmt.Errorf("the requested %s Category does not exist", category)
	}

	c := max(weekIncome-cat.Band2.Max, 0)
	b := max(weekIncome-c-cat.Band1.Max, 0)

	tax := c.Mul(cat.Band3.Rate) + b.Mul(cat.Band2.Rate)

	return tax, nil
}

// Calculate the Taxable Income of yearly gross income.
// Requirements from https://www.gov.uk/income-tax-rates
func (t TaxCalculator) calculateIncomeTax(income Money, allowance Money) IncomeTaxBreakdown {
	hrLimit := t.IncomeTaxRates.Higher.Min + (allowance - t.IncomeTaxRates.PersonalAllowance)

	ar := max(income-t.IncomeTaxRates.Additional.Min, 0)
	hr := max(income-ar-hrLimit, 0)
	br := max(income-ar-hr-allowance, 0)

	tax := ar.Mul(t.IncomeTaxRates.Additional.Rate) + hr.Mul(t.IncomeTaxRates.Higher.Rate) + br.Mul(t.IncomeTaxRates.Basic.Rate)

	return IncomeTaxBreakdown{
		GrossIncome:    income,
		BasicRate:      br.Mul(t.IncomeTaxRates.Basic.Rate),
		HigherRate:     hr.Mul(t.IncomeTaxRates.Higher.Rate),
		AdditionalRate: ar.Mul(t.IncomeTaxRates.Additional.Rate),
		Taxed:          tax,
		Taxable:        max(income-allowance, 0),
	}
}

// Calculate the Tax Allowance based on a yearly gross income.
// Requirements from https://www.gov.uk/income-tax-rates/income-over-100000
func (t TaxCalculator) calculateTaxAllowance(annumIncome Money) Money {
	over := max((annumIncome - t.IncomeTaxRates.PersonalAllowanceThreshold).Mul(0.5), 0)

	return max(t.IncomeTaxRates.PersonalAllowance-over, 0)
}

// Calculate the full income tax and return breakdown.
func (t TaxCalculator) CalculateTakeHome(income Money, niCategory string) (IncomeTaxBreakdown, error) {
	allowance := t.calculateTaxAllowance(income)
	ni, err := t.calculateNationalInsurance(income.Div(52), niCategory)
	if err != nil {
		return IncomeTaxBreakdown{}, err
	}
	tax := t.calculateIncomeTax(income, allowance)

	tax.NationalInsurance = ni.Mul(52)
	tax.TakeHome = income - tax.Taxed - tax.NationalInsurance

	return tax, nil
}
