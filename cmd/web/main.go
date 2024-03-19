package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/vfc2/tax-calculator/internal/tax"
)

type Sub struct {
	Band1 tax.Band
	Band2 tax.Band
	Band3 tax.Band
}

type Test struct {
}

func main() {
	logOptions := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, logOptions))

	views, err := NewViews()
	if err != nil {
		logger.Error("error initialising views", "error", err.Error())
		os.Exit(1)
	}

	taxConfig, err := loadTaxConfig("./assets/config/income_tax/2024_2025.json")
	if err != nil {
		logger.Error("error loading tax rates config", "error", err.Error())
		os.Exit(1)
	}

	niConfig, err := loadNIConfig("./assets/config/national_insurance/2024_2025.json")
	if err != nil {
		logger.Error("error loading national insurance rates config", "error", err.Error())
		os.Exit(1)
	}

	models := Models{
		calc: tax.TaxCalculator{
			IncomeTaxRates:         taxConfig,
			NationalInsuranceRates: niConfig,
		},
	}

	handlers := &Handlers{
		logger: logger,
		views:  views,
		models: models,
	}

	mw := &Middlewares{logger: logger}

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes(handlers, mw),
	}

	log.Fatal(server.ListenAndServe())
}

func routes(h *Handlers, mw *Middlewares) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./assets/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("GET /inputs", h.inputPage)
	mux.HandleFunc("POST /calculate", h.outputPage)

	return mw.recovery(mw.logRequest(mw.secureHeaders(mux)))
}

func serverError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	logger.Error("server error", "error", err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func loadTaxConfig(filename string) (tax.IncomeTaxRates, error) {
	tr := &tax.IncomeTaxRates{}

	f, err := os.ReadFile(filename)
	if err != nil {
		return *tr, err
	}

	err = json.Unmarshal(f, tr)
	if err != nil {
		return *tr, err
	}

	return *tr, nil
}

func loadNIConfig(filename string) (map[string]tax.NationalInsuranceRates, error) {
	ni := map[string]tax.NationalInsuranceRates{}

	f, err := os.ReadFile(filename)
	if err != nil {
		return ni, err
	}

	err = json.Unmarshal(f, &ni)
	if err != nil {
		return ni, err
	}

	return ni, nil
}
