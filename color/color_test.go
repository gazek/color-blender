package color

import (
	ic "image/color"
	"testing"
)

func TestGetColorDominanceSort(t *testing.T) {
	c := Color{
		color: ic.RGBA{R: 1, G: 2, B: 3, A: 4},
	}
	d := c.getColorDominance()
	want := []uint8{3, 2, 1}
	for i := range want {
		if *d[i] != want[i] {
			t.Errorf("Want: %v, found: %v", want[i], *d[i])
		}
	}
	*d[0]++
	if c.color.B != 4 {
		t.Errorf("Want: %v, found: %v", 4, *d[0])
	}
}

func TestGetColorDominanceSortStability(t *testing.T) {
	c := Color{
		color: ic.RGBA{R: 0, G: 0, B: 0, A: 4},
	}
	d := c.getColorDominance()
	want := []uint8{1, 2, 3}
	for i := range want {
		*d[i] = want[i]
	}
	comp := []uint8{c.color.R, c.color.G, c.color.B}
	for i := range comp {
		if comp[i] != want[i] {
			t.Errorf("Want: %v, found: %v", want[i], comp[i])
		}
	}
}
