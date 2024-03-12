package main

import (
	"fmt"
	"os"

	"github.com/vfc2/tax-calculator/internal/money"
	"github.com/vfc2/tax-calculator/internal/tax"
)

func main() {
	fmt.Print("Enter income: ")

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	wage, err := money.NewFromString(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	taxRates := tax.LoadTaxBands()

	allowance := tax.CalculateTaxAllowance(wage, taxRates)
	ni := tax.CalculateNationalInsurance(wage.Div(52)).Mul(52)
	tax := tax.CalculateIncomeTax(wage, allowance, taxRates)

	fmt.Printf("For a gross income of £%v:\n", wage.Format(2))
	fmt.Printf("Taxable income £%v\n", tax.Taxable.Format(2))
	fmt.Printf("National Insurance £%v\n", ni.Format(2))
	fmt.Printf("Allowance £%v\n", allowance.Format(2))
	fmt.Printf("Tax Due £%v\n", tax.Taxed.Format(2))
	fmt.Printf("Take Home £%v\n", (wage - tax.Taxed - ni).Format(2))
}
