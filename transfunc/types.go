package transfunc

import "image/color"

// BrightnessFunc stores a function that describes how to modify the brightness (alpha) of a Color
type BrightnessFunc struct{ transFunc }

// BrightnessFuncSlice holds a slice of BrightnessFuncs
type BrightnessFuncSlice struct{ transFuncSlice }

// WhiteLevelFunc stores a function that describes how to modify the white level of a Color
type WhiteLevelFunc struct{ transFunc }

// WhiteLevelFuncSlice holds a slice of WhiteLevelFuncs
type WhiteLevelFuncSlice struct{ transFuncSlice }

// ColorFunc stores a function that describes the transition from one Color to another
type ColorFunc struct {
	Color1 color.RGBA
	Color2 color.RGBA
	transFunc
}

// ColorFuncSlice holds a slice of ColorFuncs
type ColorFuncSlice struct{ transFuncSlice }
