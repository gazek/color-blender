package transfunc

import "image/color"

// BrightnessFunc stores a function that describes how to modify the brightness (alpha) of a Color
type BrightnessFunc struct{ transFunc }

// BrightnessFuncSlice holds a slice of BrightnessFuncs
type BrightnessFuncSlice struct{ transFuncSlice }

// GetFuncValue returns the function value for the given step and the anchor colors
func (b *BrightnessFuncSlice) GetFuncValue(stepNum int) uint8 {
	funcVal, _ := b.transFuncSlice.GetFuncValue(stepNum)
	return uint8(0xff * funcVal)
}

// WhiteLevelFunc stores a function that describes how to modify the white level of a Color
type WhiteLevelFunc struct{ transFunc }

// GetFuncValue returns the function value for the given step and the anchor colors
func (w *WhiteLevelFuncSlice) GetFuncValue(stepNum int) uint8 {
	funcVal, _ := w.transFuncSlice.GetFuncValue(stepNum)
	return uint8(0xff * funcVal)
}

// WhiteLevelFuncSlice holds a slice of WhiteLevelFuncs
type WhiteLevelFuncSlice struct{ transFuncSlice }

// ColorFunc stores a function that describes the transition from one Color to another
type ColorFunc struct {
	Color1    color.RGBA
	Color2    color.RGBA
	TransType TransType
	TransDist int
	transFunc
}

// ColorFuncSlice holds a slice of ColorFuncs
type ColorFuncSlice struct{ transFuncSlice }

// GetFuncValue returns the function value for the given step and the anchor colors
func (c *ColorFuncSlice) GetFuncValue(stepNum int) (float32, *ColorFunc) {
	funcVal, tf := c.transFuncSlice.GetFuncValue(stepNum)
	cf := tf.(*ColorFunc)
	return funcVal, cf
}

// TransType defines the type of transition
type TransType int

const (
	// OneAtATime changes only one color component at a time
	OneAtATime TransType = iota
	// AllAtOnce changes all of the color component at the same time
	AllAtOnce
	// ToWhite tansitions from color1 to white to color2
	ToWhite
	// ToBlack tansitions from color1 to black to color2
	ToBlack
)

func (t TransType) String() string {
	return [...]string{"OneAtATime", "AllAtOnce", "ToWhite", "ToBlack"}[t]
}
