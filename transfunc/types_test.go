package transfunc

import (
	imageColor "image/color"
	"testing"
)

func TestBrightnessGetFuncValue(t *testing.T) {
	tests := []struct {
		funcVal float64
		want    uint8
	}{
		{1.0, 255},
		{0, 0},
		{0.55, 140},
	}

	for _, test := range tests {
		funcs := []transFuncer{
			&transFunc{Period: 1, Function: func(x float64) float64 { return test.funcVal }},
			&transFunc{Period: 1, Function: func(x float64) float64 { return test.funcVal }},
		}
		s := &BrightnessFuncSlice{}
		s.SetFuncs(funcs)
		if result, _ := s.GetFuncValue(1); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
		}
	}
}

func TestWhiteLevelGetFuncValue(t *testing.T) {
	tests := []struct {
		funcVal float64
		want    uint8
	}{
		{1.0, 255},
		{0, 0},
		{0.55, 140},
	}

	for _, test := range tests {
		funcs := []transFuncer{
			&transFunc{Period: 1, Function: func(x float64) float64 { return test.funcVal }},
			&transFunc{Period: 1, Function: func(x float64) float64 { return test.funcVal }},
		}
		s := &WhiteLevelFuncSlice{}
		s.SetFuncs(funcs)
		if result, _ := s.GetFuncValue(1); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
		}
	}
}

func TestColorGetFuncValue(t *testing.T) {
	tests := []struct {
		funcVal   float64
		color1    imageColor.RGBA
		color2    imageColor.RGBA
		transType TransType
	}{
		{1.0, imageColor.RGBA{R: 50}, imageColor.RGBA{B: 100}, AllAtOnce},
		{1.0, imageColor.RGBA{G: 150}, imageColor.RGBA{B: 200}, AllAtOnce},
		{1.0, imageColor.RGBA{R: 75}, imageColor.RGBA{B: 175}, AllAtOnce},
	}

	for _, test := range tests {
		funcs := []transFuncer{
			&ColorFunc{
				test.color1,
				test.color2,
				test.transType,
				0,
				transFunc{
					Period:   1,
					Function: func(x float64) float64 { return test.funcVal },
				},
			},
		}
		s := &ColorFuncSlice{}
		s.SetFuncs(funcs)
		funcVal, cf := s.GetFuncValue(1)
		if funcVal != test.funcVal {
			t.Errorf("Wanted %v, got: %v", test.funcVal, funcVal)
		}
		if cf.Color1 != test.color1 {
			t.Errorf("Wanted %v, got: %v", test.color1, cf.Color1)
		}
		if cf.Color2 != test.color2 {
			t.Errorf("Wanted %v, got: %v", test.color2, cf.Color2)
		}
		if cf.TransType != test.transType {
			t.Errorf("Wanted %v, got: %v", test.transType, cf.TransType)
		}
	}
}

func TestTransTypeString(t *testing.T) {
	tests := map[string]TransType{
		"OneAtATime": OneAtATime,
		"AllAtOnce":  AllAtOnce,
		"ToWhite":    ToWhite,
		"ToBlack":    ToBlack,
	}
	want := []string{
		"OneAtATime",
		"AllAtOnce",
		"ToWhite",
		"ToBlack",
	}
	for index := range want {
		w := want[index]
		if tests[w].String() != w {
			t.Errorf("Wanted %v, got: %v", w, tests[w].String())

		}
	}
}
