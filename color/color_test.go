package color

import (
	ic "image/color"
	"testing"
)

func TestGetColorDominanceSort(t *testing.T) {
	color := ic.RGBA{R: 1, G: 2, B: 3, A: 4}
	c := Color{}
	d, names := c.GetColorDominance(&color)
	wantNames := []string{"B", "G", "R"}
	want := []uint8{3, 2, 1}
	for i := range wantNames {
		if wantNames[i] != names[i] {
			t.Errorf("Want: %v, found: %v", wantNames[i], names[i])
		}
	}
	for i := range want {
		if *d[i] != want[i] {
			t.Errorf("Want: %v, found: %v", want[i], *d[i])
		}
	}
}

func TestGetColorDominanceSortStability(t *testing.T) {
	color := ic.RGBA{R: 0, G: 0, B: 0, A: 4}
	c := Color{}
	d, names := c.GetColorDominance(&color)
	want := []uint8{1, 2, 3}
	wantNames := []string{"R", "G", "B"}
	for i := range wantNames {
		if wantNames[i] != names[i] {
			t.Errorf("Want: %v, found: %v", wantNames[i], names[i])
		}
	}
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

func TestGetComponentValue(t *testing.T) {
	color := ic.RGBA{R: 1, G: 2, B: 3, A: 4}
	c := Color{color: color}
	want := []uint8{1, 2, 3}
	keys := []string{"R", "G", "B"}
	for i := range keys {
		result := c.GetComponentValue(keys[i])
		if want[i] != result {
			t.Errorf("Want: %v, found: %v", want[i], result)
		}
	}
}

func TestSetComponentValue(t *testing.T) {
	color := ic.RGBA{R: 1, G: 2, B: 3, A: 4}
	c := Color{color: color}
	want := []uint8{10, 9, 8}
	keys := []string{"R", "G", "B"}
	for i := range keys {
		c.SetComponentValue(keys[i], want[i])
		result := c.GetComponentValue(keys[i])
		if want[i] != result {
			t.Errorf("Want: %v, found: %v", want[i], result)
		}
	}
}

func TestGetColor(t *testing.T) {
	color := ic.RGBA{R: 1, G: 2, B: 3, A: 4}
	c := Color{color: color}
	result := c.GetColor()
	if color != result {
		t.Errorf("Want: %v, found: %v", color, result)
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
	for _, test := range tests {
		if result := c.isWhite(test.color); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
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
	for _, test := range tests {
		if result := c.getWhiteLevel(test.color, nil); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
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
	for _, test := range tests {
		if result := c.getWhiteLevelComponent(test.color); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
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
	for _, test := range tests {
		if result := c.applyWhiteLevel(test.color, test.whiteLevel); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
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
		{ic.RGBA{R: 0, G: 0, B: 0}, ic.RGBA{R: whiteBaseR, G: whiteBaseG, B: whiteBaseB}},
	}

	c := Color{}
	for _, test := range tests {
		if result := c.normalizeRGBLevels(test.color); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
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
	for _, test := range tests {
		if result := c.GetBaseColor(test.color); result != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, result)
		}
	}
}

func TestSetBrightness(t *testing.T) {
	tests := []struct {
		color ic.RGBA
		alpha uint8
	}{
		{ic.RGBA{R: 255, G: 150, B: 0, A: 200}, 123},
	}

	for _, test := range tests {
		c := Color{color: test.color}
		c.SetBrightness(test.alpha)
		if c.color.A != test.alpha {
			t.Errorf("Wanted %v, got: %v", test.alpha, c.color.A)
		}
	}
}

func TestSetColor(t *testing.T) {
	tests := []struct {
		color ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 150, B: 25, A: 200}},
	}

	for _, test := range tests {
		c := Color{}
		c.SetColor(test.color)
		// the color should be stored
		if c.color != test.color {
			t.Errorf("Wanted %v, got: %v", test.color, c.color)
		}
		// the white level should be set
		if c.whiteLevel == 0 {
			t.Error("Failed to set whiteLevel")
		}
	}
}

func TestNewColor(t *testing.T) {
	tests := []struct {
		color ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 150, B: 25, A: 200}},
	}

	for _, test := range tests {
		// call the func
		c := NewColor(test.color)
		// the colors should be the same
		if test.color != c.color {
			t.Errorf("Wanted %v, got: %v", test.color, c.color)
		}
		// the white level should be set
		if c.whiteLevel == 0 {
			t.Error("Failed to set whiteLevel")
		}
	}
}

func TestSetWhiteLevel(t *testing.T) {
	tests := []struct {
		color      ic.RGBA
		whiteLevel uint8
		want       ic.RGBA
	}{
		{ic.RGBA{R: 255, G: 150, B: 25, A: 200}, 0, ic.RGBA{R: 255, G: 138, B: 0, A: 200}},
		{ic.RGBA{R: 255, G: 138, B: 0, A: 200}, 25, ic.RGBA{R: 255, G: 149, B: 25, A: 200}},
	}

	for _, test := range tests {
		c := Color{color: test.color}
		c.SetWhiteLevel(test.whiteLevel)
		// the color should be stored
		if c.color != test.want {
			t.Errorf("Wanted %v, got: %v", test.want, c.color)
		}
		// the white level should be set
		if c.whiteLevel != test.whiteLevel {
			t.Error("Failed to set whiteLevel")
		}
	}
}
