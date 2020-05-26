package color

import (
	ic "image/color"
	"testing"
)

func TestGetColorDominanceSort(t *testing.T) {
	color := ic.RGBA{R: 1, G: 2, B: 3, A: 4}
	c := Color{}
	d := c.getColorDominance(&color)
	want := []uint8{3, 2, 1}
	for i := range want {
		if *d[i] != want[i] {
			t.Errorf("Want: %v, found: %v", want[i], *d[i])
		}
	}
}

func TestGetColorDominanceSortStability(t *testing.T) {
	color := ic.RGBA{R: 0, G: 0, B: 0, A: 4}
	c := Color{}
	d := c.getColorDominance(&color)
	want := []uint8{1, 2, 3}
	for i := range want {
		*d[i] = want[i]
	}
	comp := []uint8{color.R, color.G, color.B}
	for i := range comp {
		if comp[i] != want[i] {
			t.Errorf("Want: %v, found: %v", want[i], comp[i])
		}
	}
}

func TestIsWhite(t *testing.T) {
	tests := []struct {
		color ic.RGBA
		want  bool
	}{
		{ic.RGBA{R: 0xff, G: 0xff, B: 0xff}, true},
		{ic.RGBA{R: 0xaa, G: 0xaa, B: 0xaa}, true},
		{ic.RGBA{R: 0xaa, G: 0x00, B: 0xaa}, false},
	}

	c := Color{}
	for test := range tests {
		if result := c.isWhite(tests[test].color); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestGetWhiteLevel(t *testing.T) {
	tests := []struct {
		color ic.RGBA
		want  uint8
	}{
		{ic.RGBA{R: 0xff, G: 0xff, B: 0xff}, 0xff},
		{ic.RGBA{R: 0xff, G: 0x0, B: 0x0}, 0x0},
		{ic.RGBA{R: 0x0, G: 0x0, B: 0x0}, 0x0},
		{ic.RGBA{R: 175, G: 200, B: 160}, 204},
	}

	c := Color{}
	for test := range tests {
		if result := c.getWhiteLevel(tests[test].color, nil); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestGetWhiteLevelComponent(t *testing.T) {
	tests := []struct {
		color ic.RGBA
		want  ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 255, B: 255}, ic.RGBA{R: 255, G: 255, B: 255}},
		{ic.RGBA{R: 255, G: 0x0, B: 0x0}, ic.RGBA{R: 0, G: 0, B: 0}},
		{ic.RGBA{R: 255, G: 150, B: 0x0}, ic.RGBA{R: 0, G: 0, B: 0}},
		{ic.RGBA{R: 0x0, G: 0x0, B: 0x0}, ic.RGBA{R: 0x0, G: 0x0, B: 0x0}},
		{ic.RGBA{R: 175, G: 255, B: 160}, ic.RGBA{R: 135, G: 0, B: 160}},
	}

	c := Color{}
	for test := range tests {
		if result := c.getWhiteLevelComponent(tests[test].color); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestApplyWhiteLevel(t *testing.T) {
	tests := []struct {
		color      ic.RGBA
		whiteLevel uint8
		want       ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 150, B: 0}, 0, ic.RGBA{R: 255, G: 150, B: 0}},
		{ic.RGBA{R: 255, G: 150, B: 0}, 255, ic.RGBA{R: 255, G: 255, B: 255}},
		{ic.RGBA{R: 255, G: 150, B: 0}, 127, ic.RGBA{R: 255, G: 202, B: 127}},
		{ic.RGBA{R: 150, G: 75, B: 25}, 0, ic.RGBA{R: 150, G: 60, B: 0}},
	}

	c := Color{}
	for test := range tests {
		if result := c.applyWhiteLevel(tests[test].color, tests[test].whiteLevel); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestNormalizeRGBLevels(t *testing.T) {
	tests := []struct {
		color ic.RGBA
		want  ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 150, B: 0}, ic.RGBA{R: 255, G: 150, B: 0}},
		{ic.RGBA{R: 150, G: 75, B: 25}, ic.RGBA{R: 255, G: 127, B: 42}},
	}

	c := Color{}
	for test := range tests {
		if result := c.normalizeRGBLevels(tests[test].color); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestGetBaseColor(t *testing.T) {
	tests := []struct {
		color ic.RGBA
		want  ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 150, B: 0, A: 200}, ic.RGBA{R: 255, G: 150, B: 0, A: 200}},
		{ic.RGBA{R: 150, G: 75, B: 25, A: 155}, ic.RGBA{R: 255, G: 102, B: 0, A: 155}},
		{ic.RGBA{R: 150, G: 150, B: 150, A: 123}, ic.RGBA{R: whiteBaseR, G: whiteBaseG, B: whiteBaseB, A: 123}},
	}

	c := Color{}
	for test := range tests {
		if result := c.getBaseColor(tests[test].color); result != tests[test].want {
			t.Errorf("Wanted %v, got: %v", tests[test].want, result)
		}
	}
}
