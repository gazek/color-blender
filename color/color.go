package color

import (
	ic "image/color"
	"math"
	"sort"
)

const (
	// arbitrary color used as the base color of true white
	whiteBaseR = 0xff
	whiteBaseG = 0x00
	whiteBaseB = 0xff
)

// Color type RGBA
type Color struct {
	color      ic.RGBA
	whiteLevel uint8
}

// NewColor create a new color instance and sets the color value
func NewColor(color ic.RGBA) *Color {
	c := &Color{}
	c.SetColor(color)
	return c
}

// SetColor sets the RGBA color
func (c *Color) SetColor(color ic.RGBA) {
	c.color = color
	// c.baseColor = c.getBaseColor(color)
	c.whiteLevel = c.getWhiteLevel(color, nil)
}

// SetBrightness applies a brightness level to the current color
func (c *Color) SetBrightness(alpha uint8) {
	c.color.A = alpha
}

// SetWhiteLevel applies a white level to the current color
func (c *Color) SetWhiteLevel(whiteLevel uint8) {
	c.whiteLevel = whiteLevel
	c.color = c.applyWhiteLevel(c.color, whiteLevel)
}

// GetColor returns the current color
func (c *Color) GetColor() ic.RGBA {
	return c.color
}

// GetBaseColor removes white and black from an rgb color
func (c *Color) GetBaseColor(color ic.RGBA) ic.RGBA {
	// check if the color is true white
	if c.isWhite(color) {
		return ic.RGBA{R: whiteBaseR, G: whiteBaseG, B: whiteBaseB, A: color.A}
	}
	// general case
	// remove white level components from color
	couleurSansBlanc := c.applyWhiteLevel(color, 0)
	// scale the two non-zero component colors so that the dominant color is 255
	baseColor := c.normalizeRGBLevels(couleurSansBlanc)
	// return the base color
	return baseColor
}

// GetColorDominance returns a slice of pointers and color component names sorted descending by color component value
func (c *Color) GetColorDominance(color *ic.RGBA) (domPointers []*uint8, names []string) {
	// create a slice of pointers to the RGB values
	domPointers = []*uint8{&color.R, &color.G, &color.B}
	// sort the slice by the underlying uint8 values
	sort.SliceStable(domPointers, func(i, j int) bool { return *domPointers[j] < *domPointers[i] })
	// get the names of each location
	names = c.getcolorDominanceNames(color, domPointers)
	// return the sorted slice and names
	return domPointers, names
}

// getcolorDominanceNames returns the color component names sorted descending by color component value
func (c *Color) getcolorDominanceNames(color *ic.RGBA, domPointers []*uint8) []string {
	var result []string
	for d := range domPointers {
		switch domPointers[d] {
		case &color.R:
			result = append(result, "R")
		case &color.G:
			result = append(result, "G")
		case &color.B:
			result = append(result, "B")
		default:
			panic("Failed to match dominance component")
		}
	}
	return result
}

// GetComponentValue returns the component value of the color by component string
func (c *Color) GetComponentValue(compName string) uint8 {
	switch compName {
	case "R":
		return c.color.R
	case "G":
		return c.color.G
	case "B":
		return c.color.B
	default:
		panic("Invalid color component name")
	}
}

// SetComponentValue set the R, G or B component value
func (c *Color) SetComponentValue(compName string, value uint8) {
	// set the color component
	switch compName {
	case "R":
		c.color.R = value
	case "G":
		c.color.G = value
	case "B":
		c.color.B = value
	}
	// set the color
	c.SetColor(c.color)
}

// isWhite check if a color is true white
func (c *Color) isWhite(color ic.RGBA) bool {
	// it's true white if each of the three color components are equal
	if color.R == color.G && color.G == color.B {
		return true
	}
	return false
}

// getWhiteLevelComponent calculates an RGB tuple that represents the white components of the color
func (c *Color) getWhiteLevelComponent(color ic.RGBA) ic.RGBA {
	// check if the color is true white
	if c.isWhite(color) {
		return color
	}
	// general case
	whiteLevelComp := ic.RGBA{R: color.R, G: color.G, B: color.B}
	// get component dominaces if needed
	dominance, _ := c.GetColorDominance(&whiteLevelComp)
	// calculate the middle dominnce component's white level amount
	whiteLevel := c.getWhiteLevel(whiteLevelComp, dominance)
	*dominance[1] = *dominance[1] - uint8(((int(*dominance[1])*math.MaxUint8)-int(*dominance[0])*int(whiteLevel))/(math.MaxUint8-int(whiteLevel)))
	// set dominant color component to 0 since none of it contributes to white
	*dominance[0] = 0
	// we leave the least dominant color component as is, since all of its value contributes to the white level
	// return the result
	return whiteLevelComp
}

// getWhiteLevel calculates the white level of the color
func (c *Color) getWhiteLevel(color ic.RGBA, dominance []*uint8) uint8 {
	if dominance == nil {
		dominance, _ = c.GetColorDominance(&color)
	}
	if *dominance[0] == 0 {
		return 0
	}
	return uint8(math.MaxUint8 * float32(*dominance[2]) / float32(*dominance[0]))
}

// applyWhiteLevel applies a white level to a color
func (c *Color) applyWhiteLevel(color ic.RGBA, whiteLevel uint8) ic.RGBA {
	// get the white level components of the color
	wlc := c.getWhiteLevelComponent(color)
	// create a new RGB object that contains the color with zero white level
	result := ic.RGBA{
		R: color.R - wlc.R,
		G: color.G - wlc.G,
		B: color.B - wlc.B,
		A: color.A,
	}
	// if the white level is zero then we are done
	if whiteLevel == 0 {
		return result
	}
	// get the color dominance
	dom, _ := c.GetColorDominance(&result)
	// modify the white level
	for i := 0; i < 3; i++ {
		adj := uint8((int(*dom[0] - *dom[i])) * int(whiteLevel) / math.MaxUint8)
		*dom[i] += adj
	}
	// return the new color
	return result
}

// applyWhiteLevel applies a white level to a color
func (c *Color) normalizeRGBLevels(color ic.RGBA) ic.RGBA {
	// create a new RGB object to store the resulting color
	result := ic.RGBA{
		R: color.R,
		G: color.G,
		B: color.B,
		A: color.A,
	}
	// get the color dominance
	dom, _ := c.GetColorDominance(&result)
	// modify the white level
	adj := math.MaxUint8 / float32(*dom[0])
	for i := 0; i < 3; i++ {
		*dom[i] = uint8(float32(*dom[i]) * adj)
	}
	return result
}
