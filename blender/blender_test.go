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
