package transfunc

import (
	"testing"
)

func TestGetFuncValue(t *testing.T) {
	tests := []struct {
		function   func(x float32) float32
		period     int
		inputRange []float32
		stepNum    int
		want       float32
	}{
		{func(x float32) float32 { return 10 - x }, 10, []float32{0, 10}, 5, 5},
		{func(x float32) float32 { return 10 - x }, 10, []float32{0, 10}, 17, 3},
	}

	for test := range tests {
		f := transFunc{Function: tests[test].function, Period: tests[test].period, InputRange: tests[test].inputRange}
		if result := f.GetFuncValue(tests[test].stepNum); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestGetFuncPeriod(t *testing.T) {
	tests := []struct {
		period int
	}{
		{10},
	}

	for test := range tests {
		f := transFunc{Period: tests[test].period}
		if result := f.GetFuncPeriod(); result != tests[test].period {
			t.Errorf("Wanted %v, got: %v", tests[test].period, result)
		}
	}
}
