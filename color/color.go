package color

import (
	ic "image/color"
	"sort"
)

// Color type RGBA
type Color struct {
	baseColor  ic.RGBA
	color      ic.RGBA
	whiteLevel uint8
}

// NewColor create a new color instance and sets the color value
// func NewColor(color ic.RGBA) *Color {
// 	c := &Color{}
// 	c.SetColor(color)
// 	return c
// }

// // SetColor sets the RGBA color
// func (c *Color) SetColor(color ic.RGBA) {
// 	c.color = color
// 	c.baseColor = c.getBaseColor(color)
// 	c.whiteLevel = c.getWhiteLevel(color)
// }

// // SetBrightness applies a brightness level to the current color
// func (c *Color) SetBrightness(alpha uint8) {
// 	c.color.A = alpha
// }

// // SetWhiteLevel applies a white level to the current color
// func (c *Color) SetWhiteLevel(whiteLevel uint8) {
// 	c.whiteLevel = whiteLevel
// 	c.color = c.applyWhiteLevelToBase(whiteLevel)
// }

// // getBaseColor removes white and black from an rgb color
// func (c *Color) getBaseColor(color ic.RGBA) ic.RGBA {

// }

// // applyWhiteLevelToBase applies the white level to the current base color store in c.color
// func (c *Color) applyWhiteLevelToBase(whiteLevel uint8) {

// }

// // getWhiteLevel calculates the white level of the color
// func (c *Color) getWhiteLevel(color ic.RGBA) uint8 {

// }

func (c *Color) getColorDominance() []*uint8 {
	// create a slice of pointers to the RGB values
	rgbPointers := []*uint8{&c.color.R, &c.color.G, &c.color.B}
	// sort the slice by the underlying uint8 values
	sort.SliceStable(rgbPointers, func(i, j int) bool { return *rgbPointers[j] < *rgbPointers[i] })
	// return the sorted slice
	return rgbPointers
}
