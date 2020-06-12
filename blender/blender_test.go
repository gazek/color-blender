package blender

import (
	ic "image/color"
	"math"
	"testing"

	"github.com/gazek/color-blender/color"
	"github.com/gazek/color-blender/transfunc"
)

func TestSetComponentWithConstraint(t *testing.T) {
	tests := []struct {
		color     ic.RGBA
		component string
		value     uint8
		maxDist   int
		want      ic.RGBA
	}{
		{ic.RGBA{R: 0, G: 0, B: 0}, "R", 200, 255, ic.RGBA{R: 200, G: 0, B: 0}},
		{ic.RGBA{R: 0, G: 0, B: 0}, "G", 200, 255, ic.RGBA{R: 0, G: 200, B: 0}},
		{ic.RGBA{R: 0, G: 0, B: 0}, "B", 200, 255, ic.RGBA{R: 0, G: 0, B: 200}},
		{ic.RGBA{R: 0, G: 0, B: 0}, "B", 200, 0, ic.RGBA{R: 0, G: 0, B: 0}},
		{ic.RGBA{R: 0, G: 0, B: 0}, "B", 200, 75, ic.RGBA{R: 0, G: 0, B: 75}},
		{ic.RGBA{R: 0, G: 0, B: 200}, "B", 0, 75, ic.RGBA{R: 0, G: 0, B: 125}},
	}

	b := Blender{}
	for test := range tests {
		c := color.Color{}
		c.SetColor(tests[test].color)
		b.setComponentWithConstraint(&c, tests[test].component, tests[test].value, tests[test].maxDist)
		if c.GetColor() != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, c.GetColor())
		}
	}
}

func TestOneAtATimeColorTransition(t *testing.T) {
	tests := []struct {
		color1    ic.RGBA
		color2    ic.RGBA
		percent   float64
		wantDist  int
		wantColor ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, 0.5, 510, ic.RGBA{R: 255, G: 255, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, 0, 510, ic.RGBA{R: 255, G: 0, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, 1, 510, ic.RGBA{R: 0, G: 255, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, 0.33333, 510, ic.RGBA{R: 255, G: 170, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 45, B: 255}, 0.75, 555, ic.RGBA{R: 94, G: 0, B: 255}},
	}

	b := &Blender{}
	for index := range tests {
		test := tests[index]
		_, dist := b.oneAtATimeColorTransition(test.color1, test.color2, 10*math.MaxUint8)
		if dist != test.wantDist {
			t.Errorf("Wanted: %v, got: %v", test.wantDist, dist)
		}
		color, _ := b.oneAtATimeColorTransition(test.color1, test.color2, int(math.Round(float64(test.wantDist)*test.percent)))
		if color != test.wantColor {
			t.Errorf("Wanted %v, got: %v", test.wantColor, color)
		}
	}
}

func TestGetColorTransitionColorAndDistance(t *testing.T) {
	tests := []struct {
		color1    ic.RGBA
		color2    ic.RGBA
		transType transfunc.TransType
		percent   float64
		wantDist  int
		wantColor ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 255, G: 0, B: 0}, transfunc.OneAtATime, 0.5, 0, ic.RGBA{R: 255, G: 0, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 0.5, 510, ic.RGBA{R: 255, G: 255, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 0, 510, ic.RGBA{R: 255, G: 0, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 1, 510, ic.RGBA{R: 0, G: 255, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 0.33333, 510, ic.RGBA{R: 255, G: 170, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 45, B: 255}, transfunc.OneAtATime, 0.75, 555, ic.RGBA{R: 94, G: 0, B: 255}},
	}

	b := &Blender{}
	for index := range tests {
		test := tests[index]
		cf := &transfunc.ColorFunc{
			Color1:    test.color1,
			Color2:    test.color2,
			TransType: test.transType,
		}
		b.getColorTransitionDistance(cf)
		if cf.TransDist != test.wantDist {
			t.Errorf("Wanted: %v, got: %v", test.wantDist, cf.TransDist)
		}
		color := b.getColorTransitionColor(cf, int(math.Round(float64(test.wantDist)*test.percent)))
		if color != test.wantColor {
			t.Errorf("Wanted %v, got: %v", test.wantColor, color)
		}
	}
}

func TestGetTransitionColor(t *testing.T) {
	tests := []struct {
		color1    ic.RGBA
		color2    ic.RGBA
		transType transfunc.TransType
		percent   float64
		wantDist  int
		wantColor ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 0.5, 510, ic.RGBA{R: 255, G: 255, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 0, 510, ic.RGBA{R: 255, G: 0, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 1, 510, ic.RGBA{R: 0, G: 255, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 255, B: 0}, transfunc.OneAtATime, 0.33333, 510, ic.RGBA{R: 255, G: 170, B: 0}},
		{ic.RGBA{R: 255, G: 0, B: 0}, ic.RGBA{R: 0, G: 45, B: 255}, transfunc.OneAtATime, 0.75, 555, ic.RGBA{R: 94, G: 0, B: 255}},
	}

	b := &Blender{}
	for index := range tests {
		test := tests[index]
		cf := &transfunc.ColorFunc{
			Color1:    test.color1,
			Color2:    test.color2,
			TransType: test.transType,
		}
		color := b.getTransitionColor(float32(test.percent), cf)
		if color != test.wantColor {
			t.Errorf("Wanted %v, got: %v", test.wantColor, color)
		}
	}
}

func TestGetPeriod(t *testing.T) {
	tests := []struct {
		colorFuncPeriod      int
		brightnessFuncPeriod int
		whiteLevelFuncPeriod int
		want                 int
	}{
		{0, 0, 0, 0},
		{1, 2, 3, 6},
	}

	for index := range tests {
		test := tests[index]
		b := Blender{}
		cf := transfunc.ColorFunc{}
		cf.Period = test.colorFuncPeriod
		b.AppendColorFunc(cf)
		bf := transfunc.BrightnessFunc{}
		bf.Period = test.brightnessFuncPeriod
		b.AppendBrightnessFunc(bf)
		wf := transfunc.WhiteLevelFunc{}
		wf.Period = test.whiteLevelFuncPeriod
		b.AppendWhiteLevelFunc(wf)
		period := b.getPeriod()
		if period != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, period)
		}
	}
}

func TestAdvanceStepAndResetStep(t *testing.T) {
	tests := []struct {
		colorFuncPeriod      int
		brightnessFuncPeriod int
		whiteLevelFuncPeriod int
		stepStart            int
		step                 int
		want                 int
	}{
		{1, 2, 3, 0, 1, 1},
		{1, 2, 3, 1, 1, 2},
		{1, 2, 3, 5, 2, 1},
		{1, 2, 3, -3, 7, 4},
		{1, 2, 3, 4, -9999, 0},
		{0, 0, 0, 0, 10, 0},
	}

	for index := range tests {
		test := tests[index]
		// build the blender
		b := Blender{}
		cf := transfunc.ColorFunc{}
		cf.Period = test.colorFuncPeriod
		b.AppendColorFunc(cf)
		bf := transfunc.BrightnessFunc{}
		bf.Period = test.brightnessFuncPeriod
		b.AppendBrightnessFunc(bf)
		wf := transfunc.WhiteLevelFunc{}
		wf.Period = test.whiteLevelFuncPeriod
		b.AppendWhiteLevelFunc(wf)
		// set the step
		b.step = test.stepStart
		// advance the step
		b.AdvanceStep(test.step)
		if b.step != test.want {
			t.Errorf("Wanted: %v, found: %v", test.want, b.step)
		}
		// reset step
		b.ResetStep()
		if b.step != 0 {
			t.Errorf("Wanted: %v, found: %v", 0, b.step)
		}
	}
}

func TestGetGetColor(t *testing.T) {
	tests := []struct {
		color1               ic.RGBA
		color2               ic.RGBA
		transType            transfunc.TransType
		colorFunc            func(x float32) float32
		colorFuncRange       []float32
		colorFuncPeriod      int
		whiteLevelFunc       func(x float32) float32
		whiteLevelFuncRange  []float32
		whiteLevelFuncPeriod int
		brightnessFunc       func(x float32) float32
		brightnessFuncRange  []float32
		brightnessFuncPeriod int
		step                 int
		want                 ic.RGBA
	}{
		{
			color1:    ic.RGBA{R: 255},
			color2:    ic.RGBA{B: 255},
			transType: transfunc.OneAtATime,
			colorFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			colorFuncRange:  []float32{0, 10},
			colorFuncPeriod: 10,
			whiteLevelFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			whiteLevelFuncRange:  []float32{0, 10},
			whiteLevelFuncPeriod: 10,
			brightnessFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			brightnessFuncRange:  []float32{0, 10},
			brightnessFuncPeriod: 10,
			step:                 0,
			want:                 ic.RGBA{R: 255},
		},
		{
			color1:    ic.RGBA{R: 255},
			color2:    ic.RGBA{B: 255},
			transType: transfunc.OneAtATime,
			colorFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			colorFuncRange:  []float32{0, 10},
			colorFuncPeriod: 10,
			whiteLevelFunc: func(x float32) float32 {
				if x < 10 {
					return 0
				}
				return 1
			},
			whiteLevelFuncRange:  []float32{0, 10},
			whiteLevelFuncPeriod: 10,
			brightnessFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			brightnessFuncRange:  []float32{0, 10},
			brightnessFuncPeriod: 10,
			step:                 9,
			want:                 ic.RGBA{B: 255, A: 255},
		},
		{
			color1:    ic.RGBA{R: 255},
			color2:    ic.RGBA{B: 255},
			transType: transfunc.OneAtATime,
			colorFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			colorFuncRange:  []float32{0, 10},
			colorFuncPeriod: 10,
			whiteLevelFunc: func(x float32) float32 {
				if x < 5 {
					return 0
				}
				return 1
			},
			whiteLevelFuncRange:  []float32{0, 10},
			whiteLevelFuncPeriod: 10,
			brightnessFunc: func(x float32) float32 {
				if x < 10 {
					return 0
				}
				return 1
			},
			brightnessFuncRange:  []float32{0, 10},
			brightnessFuncPeriod: 10,
			step:                 9,
			want:                 ic.RGBA{R: 255, G: 255, B: 255, A: 0},
		},
	}

	for index := range tests {
		test := tests[index]
		// create the blender
		b := Blender{}
		// add color func
		cf := transfunc.NewColorFunc(test.color1, test.color2, test.transType, test.colorFunc, test.colorFuncPeriod, test.colorFuncRange)
		b.AppendColorFunc(cf)
		// add white level func
		wf := transfunc.NewWhiteLevelFunc(test.whiteLevelFunc, test.whiteLevelFuncPeriod, test.whiteLevelFuncRange)
		b.AppendWhiteLevelFunc(wf)
		// add brightness func
		bf := transfunc.NewBrightnessFunc(test.brightnessFunc, test.brightnessFuncPeriod, test.brightnessFuncRange)
		b.AppendBrightnessFunc(bf)
		// set the step
		b.AdvanceStep(test.step)
		// get the color
		color := b.GetColor()
		// check the result
		if color.GetColor() != test.want {
			t.Errorf("Wanted: %v, found: %v", test.want, color.GetColor())
		}
	}
}
