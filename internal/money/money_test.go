package money

import (
	"testing"
)

func TestMoneyMul(t *testing.T) {
	tests := map[string]struct {
		base      float64
		mul       float64
		expected  string
		precision int
	}{
		"453*0.02": {
			base:      453,
			mul:       0.02,
			expected:  "9.06",
			precision: 2,
		},
		"120*0": {
			base:     120,
			mul:      0,
			expected: "0",
		},
		"23453*0.34592": {
			base:      23453,
			mul:       0.34592,
			expected:  "8112.86176",
			precision: 5,
		},
		"4.33334*13.456": {
			base:      4.33334,
			mul:       13.456,
			expected:  "58.31",
			precision: 2,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := New(test.base).Mul(test.mul).Format(test.precision)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestMoneyDiv(t *testing.T) {
	tests := map[string]struct {
		base      float64
		div       float64
		expected  string
		precision int
	}{
		"400/4": {
			base:      400,
			div:       4,
			expected:  "100.0",
			precision: 1,
		},
		"763.45/0": {
			base:      763.45,
			div:       0,
			expected:  "763.45",
			precision: 2,
		},
		"1567.45/3.67": {
			base:      1567.45,
			div:       3.67,
			expected:  "427.098093",
			precision: 6,
		},
		"-45639.56784/936.1078": {
			base:      -45639.56784,
			div:       936.1078,
			expected:  "-48.755",
			precision: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := New(test.base).Div(test.div).Format(test.precision)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestMoneySub(t *testing.T) {
	tests := map[string]struct {
		base     int64
		sub      float64
		expected int64
	}{
		"657-17": {
			base:     657,
			sub:      17,
			expected: 640000000,
		},
		"1456-2491": {
			base:     1456,
			sub:      2491,
			expected: -1035000000,
		},
		"44-0": {
			base:     44,
			sub:      0,
			expected: 44000000,
		},
		"149-3.04567": {
			base:     149,
			sub:      3.04567,
			expected: 145954330,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := New(test.base).Sub(test.sub)

			if actual != Money(test.expected) {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestMoneyAdd(t *testing.T) {
	tests := map[string]struct {
		base     float64
		add      float64
		expected int64
	}{
		"45+67": {
			base:     67,
			add:      45,
			expected: 112000000,
		},
		"34+1678.68": {
			base:     34,
			add:      1678.68,
			expected: 1712680000,
		},
		"113+0": {
			base:     113,
			add:      0,
			expected: 113000000,
		},
		"-9067.34+45.678": {
			base:     -9067.34,
			add:      45.678,
			expected: -9021662000,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := New(test.base).Add(test.add)

			if actual != Money(test.expected) {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestMoneyFormat(t *testing.T) {
	tests := map[string]struct {
		base     string
		digits   int
		expected string
	}{
		"453": {
			base:     "453",
			digits:   0,
			expected: "453",
		},
		"66.5": {
			base:     "66.498",
			digits:   1,
			expected: "66.5",
		},
		"129.90": {
			base:     "129.9",
			digits:   2,
			expected: "129.90",
		},
		"142": {
			base:     "141.657",
			digits:   0,
			expected: "142",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v, _ := NewFromString(test.base)
			actual := v.Format(test.digits)

			if actual != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestMoneyNewFromString(t *testing.T) {
	tests := map[string]struct {
		base     string
		expected string
	}{
		"129.90": {
			base:     "345.22654",
			expected: "345.23",
		},
	}

	tests_fail := map[string]struct {
		base string
	}{
		"InvalidCharacters": {
			base: "6hgX.e",
		},
		"Empty": {
			base: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, _ := NewFromString(test.base)

			if actual.Format(2) != test.expected {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}

	for name, test := range tests_fail {
		t.Run(name, func(t *testing.T) {
			actual, err := NewFromString(test.base)

			if actual != 0 || err == nil {
				t.Error("an error was expected but not returned")
			}
		})
	}
}
