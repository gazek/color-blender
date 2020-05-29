package blender

import (
	imageColor "image/color"

	"github.com/gazek/color-blender/color"
	"github.com/gazek/color-blender/transfunc"
)

// Blender modifies a color over time according to the provided color, brightness and white level transition functions
type Blender struct {
	colorFuncs      transfunc.ColorFuncSlice
	brightnessFuncs transfunc.BrightnessFuncSlice
	whiteLevelFuncs transfunc.WhiteLevelFuncSlice
	step            int
}

// ResetStep sets the step position to zero
func (b *Blender) ResetStep() {
	b.step = 0
}

// AdvanceStep changes the current step position by the numSteps amount
func (b *Blender) AdvanceStep(numSteps int) {
	// get a common period
	period := b.getPeriod()
	// set the step
	b.step = (b.step + numSteps) % period
}

func (b *Blender) getPeriod() int {
	// I don't want to write a LCM function
	// so I'll just do it the easy way
	isZero := true
	period := 1
	// multiply all the periods
	if b.colorFuncs.GetPeriod() > 0 {
		period *= b.colorFuncs.GetPeriod()
		isZero = false
	}
	if b.brightnessFuncs.GetPeriod() > 0 {
		period *= b.brightnessFuncs.GetPeriod()
		isZero = false
	}
	if b.whiteLevelFuncs.GetPeriod() > 0 {
		period *= b.whiteLevelFuncs.GetPeriod()
		isZero = false
	}
	// return the result
	if isZero {
		return 0
	}
	return period
}

// GetColor calculates the color for the current step position
func (b *Blender) GetColor() *color.Color {
	// create a new Color object to hold the result
	result := &color.Color{}
	// get the color func value
	cfv, c1, c2, transType := b.colorFuncs.GetFuncValue(b.step)
	// get the base color resulting from the func value
	result.SetColor(b.getColorTransition(c1, c2, cfv, transType))
	// get the brightness func value
	bfv := b.brightnessFuncs.GetFuncValue(b.step)
	// apply the brightness to the base color
	result.SetBrightness(bfv)
	// get the white level func value
	wlfv := b.whiteLevelFuncs.GetFuncValue(b.step)
	// apply the white level to the base color
	result.SetWhiteLevel(wlfv)
	// return the resulting color
	return result
}

func (b *Blender) getColorTransition(color1 imageColor.RGBA, color2 imageColor.RGBA, transVal float32, transType string) imageColor.RGBA {
	panic("NotImplemented")
}

// // GetColorWindow calculates the colors for the next n step positions, where n is the length of the slice pointer provided
// func (b *Blender) GetColorWindow(window *[]color.Color) {
// 	panic("NotImplemented")
// }

// AppendColorFunc appends the ColorFunc to the ColorFuncSlice
func (b *Blender) AppendColorFunc(f transfunc.ColorFunc) {
	b.colorFuncs.AppendFunc(&f)
}

// AppendBrightnessFunc appends the ColorFunc to the ColorFuncSlice
func (b *Blender) AppendBrightnessFunc(f transfunc.BrightnessFunc) {
	b.brightnessFuncs.AppendFunc(&f)
}

// AppendWhiteLevelFunc appends the ColorFunc to the ColorFuncSlice
func (b *Blender) AppendWhiteLevelFunc(f transfunc.WhiteLevelFunc) {
	b.whiteLevelFuncs.AppendFunc(&f)
}
