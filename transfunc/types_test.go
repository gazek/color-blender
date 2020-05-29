package transfunc

import (
	imageColor "image/color"
	"testing"
)

func TestBrightnessGetFuncValue(t *testing.T) {
	tests := []struct {
		funcVal float32
		want    uint8
	}{
		{1.0, 255},
		{0, 0},
		{0.55, 140},
	}

	for test := range tests {
		funcs := []transFuncer{
			&transFunc{Period: 1, Function: func(x float32) float32 { return tests[test].funcVal }},
			&transFunc{Period: 1, Function: func(x float32) float32 { return tests[test].funcVal }},
		}
		s := &BrightnessFuncSlice{}
		s.SetFuncs(funcs)
		if result := s.GetFuncValue(1); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestWhiteLevelGetFuncValue(t *testing.T) {
	tests := []struct {
		funcVal float32
		want    uint8
	}{
		{1.0, 255},
		{0, 0},
		{0.55, 140},
	}

	for test := range tests {
		funcs := []transFuncer{
			&transFunc{Period: 1, Function: func(x float32) float32 { return tests[test].funcVal }},
			&transFunc{Period: 1, Function: func(x float32) float32 { return tests[test].funcVal }},
		}
		s := &WhiteLevelFuncSlice{}
		s.SetFuncs(funcs)
		if result := s.GetFuncValue(1); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestColorGetFuncValue(t *testing.T) {
	tests := []struct {
		funcVal   float32
		color1    imageColor.RGBA
		color2    imageColor.RGBA
		transType string
	}{
		{1.0, imageColor.RGBA{R: 50}, imageColor.RGBA{B: 100}, "base"},
		{1.0, imageColor.RGBA{G: 150}, imageColor.RGBA{B: 200}, "base"},
		{1.0, imageColor.RGBA{R: 75}, imageColor.RGBA{B: 175}, "base"},
	}

	for test := range tests {
		funcs := []transFuncer{
			&ColorFunc{
				tests[test].color1,
				tests[test].color2,
				tests[test].transType,
				transFunc{
					Period:   1,
					Function: func(x float32) float32 { return tests[test].funcVal },
				},
			},
		}
		s := &ColorFuncSlice{}
		s.SetFuncs(funcs)
		funcVal, color1, color2, transType := s.GetFuncValue(1)
		if funcVal != tests[test].funcVal {
			t.Errorf("Wanted %v, got: %v", tests[test].funcVal, funcVal)
		}
		if color1 != tests[test].color1 {
			t.Errorf("Wanted %v, got: %v", tests[test].color1, color1)
		}
		if color2 != tests[test].color2 {
			t.Errorf("Wanted %v, got: %v", tests[test].color2, color2)
		}
		if transType != tests[test].transType {
			t.Errorf("Wanted %v, got: %v", tests[test].transType, transType)
		}
	}
}
