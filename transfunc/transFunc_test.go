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

	for _, test := range tests {
		f := transFunc{Function: test.function, Period: test.period, InputRange: test.inputRange}
		if result := f.GetFuncValue(test.stepNum); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
		}
	}
}

func TestGetFuncPeriod(t *testing.T) {
	tests := []struct {
		period int
	}{
		{10},
	}

	for _, test := range tests {
		f := transFunc{Period: test.period}
		if result := f.GetFuncPeriod(); result != test.period {
			t.Errorf("Wanted %v, got: %v", test.period, result)
		}
	}
}
