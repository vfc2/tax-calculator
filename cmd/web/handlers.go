package main

import (
	"log/slog"
	"net/http"

	"github.com/vfc2/tax-calculator/internal/money"
)

type Handlers struct {
	logger *slog.Logger
	views  *Views
	models Models
}

type TaxInputValidation struct {
	Errors map[string]string
}

func (h Handlers) home(w http.ResponseWriter, r *http.Request) {
	h.views.render(w, "home", "layout", nil, h.logger)
}

func (h Handlers) inputPage(w http.ResponseWriter, r *http.Request) {
	h.views.render(w, "tax_input", "view", nil, h.logger)
}

func (h Handlers) outputPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		serverError(w, r, err, h.logger)
	}

	income := r.PostForm.Get("income")
	period := r.PostForm.Get("period")

	wage, err := money.NewFromString(income)
	if err != nil {
		val := TaxInputValidation{
			Errors: map[string]string{
				"income": "The value must be a valid number.",
			},
		}

		h.views.render(w, "tax_input", "view", val, h.logger)
		return
	}

	switch period {
	case "Month":
		wage = wage.Mul(12)
	case "Week":
		wage = wage.Mul(52)
	}

	tax, err := h.models.calc.CalculateTakeHome(wage, "A")
	if err != nil {
		serverError(w, r, err, h.logger)
	}

	h.views.render(w, "tax_output", "view", tax, h.logger)
}
