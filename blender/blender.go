package blender

import (
	"fmt"
	imageColor "image/color"
	"math"

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
	// handle period of zero
	if period == 0 {
		b.step = 0
		return
	}
	// set the step
	b.step = (b.step + numSteps) % period
	// don't let the step position go negative
	if b.step < 0 {
		b.step = 0
	}
}

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

// GetColor calculates the color for the current step position
func (b *Blender) GetColor() *color.Color {
	// create a new Color object to hold the result
	result := &color.Color{}
	// get the color func value
	cfv, cf := b.colorFuncs.GetFuncValue(b.step)
	// get the base color resulting from the func value
	result.SetColor(b.getTransitionColor(cf, cfv))
	// get the brightness func value
	bfv, ok := b.brightnessFuncs.GetFuncValue(b.step)
	// apply the brightness to the base color
	if ok {
		result.SetBrightness(bfv)
	}
	// get the white level func value
	wlfv, ok := b.whiteLevelFuncs.GetFuncValue(b.step)
	// apply the white level to the base color
	if ok {
		result.SetWhiteLevel(wlfv)
	}
	// return the resulting color
	return result
}

func (b *Blender) getPeriod() int {
	// I don't want to write an LCM function
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

// getColorTransTypeFunc selects the transition function
func (b *Blender) getColorTransTypeFunc(transType transfunc.TransType) func(colorFunc *transfunc.ColorFunc, transPercent float32) imageColor.RGBA {
	switch transType {
	case transfunc.OneAtATime:
		return b.oneAtATimeColorTransition
	case transfunc.AllAtOnce:
		return b.allAtOnceColorTransition
	case transfunc.ToWhite:
		return b.whiteColorTransition
	case transfunc.ToBlack:
		return b.blackColorTransition
	default:
		panic(fmt.Sprintf("Invalid color transition type: %v", transType))
	}
}

// getTransitionColor retreives the transition function and gets the transition color
func (b *Blender) getTransitionColor(colorFunc *transfunc.ColorFunc, transPercent float32) imageColor.RGBA {
	transTypeFunc := b.getColorTransTypeFunc(colorFunc.TransType)
	color := transTypeFunc(colorFunc, transPercent)
	return color
}

// oneAtATimeColorTransition transitions between colors by changing only one component value at a time
func (b *Blender) oneAtATimeColorTransition(colorFunc *transfunc.ColorFunc, transPercent float32) imageColor.RGBA {
	// get the full transition distance if we don't already have it
	if colorFunc.TransDist <= 0 {
		_, colorFunc.TransDist = b._oneAtATimeColorTransition(colorFunc.Color1, colorFunc.Color2, 4*math.MaxUint8)
	}
	maxDist := int(math.Round(float64(transPercent * float32(colorFunc.TransDist))))
	color, _ := b._oneAtATimeColorTransition(colorFunc.Color1, colorFunc.Color2, maxDist)
	return color
}

// oneAtATimeColorTransition transitions between colors by changing only one component value at a time
func (b *Blender) _oneAtATimeColorTransition(color1 imageColor.RGBA, color2 imageColor.RGBA, maxDist int) (resultingColor imageColor.RGBA, distance int) {
	// track the transition distance
	var dist int
	// initialize the RGBA to color1
	result := color.NewColor(color1)
	// get the components by dmominance
	c1DomPtrs, c1DomNames := color.NewColor(color1).GetColorDominance(&color1)
	c2DomPtrs, c2DomNames := color.NewColor(color2).GetColorDominance(&color2)
	// avoid backtracking across the same path
	if *c1DomPtrs[1] != *c2DomPtrs[0] {
		// result[c1.d1] => 0
		dist += b.setComponentWithConstraint(result, c1DomNames[1], 0, maxDist-dist)
	}
	// result[c2.d0] => 255
	dist += b.setComponentWithConstraint(result, c2DomNames[0], 255, maxDist-dist)
	// result[c2.d2] => 0
	dist += b.setComponentWithConstraint(result, c2DomNames[2], 0, maxDist-dist)
	// result[c2.d1] => c2[c2.d1]
	dist += b.setComponentWithConstraint(result, c2DomNames[1], *c2DomPtrs[1], maxDist-dist)
	// return the color
	return result.GetColor(), dist
}

// allAtOnceColorTransition transitions between colors by changing all component values at once, moving them directly toward the target values
func (b *Blender) allAtOnceColorTransition(cf *transfunc.ColorFunc, transPercent float32) imageColor.RGBA {
	return imageColor.RGBA{
		R: cf.Color1.R + uint8(float32((cf.Color2.R-cf.Color1.R))*transPercent),
		G: cf.Color1.G + uint8(float32((cf.Color2.G-cf.Color1.G))*transPercent),
		B: cf.Color1.B + uint8(float32((cf.Color2.B-cf.Color1.B))*transPercent),
	}
}

// whiteColorTransition similar to allAtOnceColorTransition but transitions to white before transitioning to the target values
func (b *Blender) whiteColorTransition(colorFunc *transfunc.ColorFunc, transPercent float32) imageColor.RGBA {
	return imageColor.RGBA{}
}

// blackColorTransition similar to allAtOnceColorTransition but transitions to black before transitioning to the target values
func (b *Blender) blackColorTransition(colorFunc *transfunc.ColorFunc, transPercent float32) imageColor.RGBA {
	return imageColor.RGBA{}
}

// setComponentWithConstraint sets a color component value but restricts the amount the value can deviate from its current value
func (b *Blender) setComponentWithConstraint(color *color.Color, compName string, value uint8, maxDist int) (distTraveled int) {
	// check the base case
	if maxDist <= 0 {
		return 0
	}
	// calculate the component change
	compValue := color.GetComponentValue(compName)
	change := int(value) - int(compValue)
	// get the distance
	dist := int(math.Abs(float64(change)))
	// set the component value
	if dist <= maxDist {
		color.SetComponentValue(compName, value)
		distTraveled = dist
	} else {
		if change > 0 {
			color.SetComponentValue(compName, color.GetComponentValue(compName)+uint8(maxDist))
		} else {
			color.SetComponentValue(compName, color.GetComponentValue(compName)-uint8(maxDist))
		}
		distTraveled = maxDist
	}
	return distTraveled
}

// // GetColorWindow calculates the colors for the next n step positions, where n is the length of the slice pointer provided
// func (b *Blender) GetColorWindow(window *[]color.Color) {
// 	panic("NotImplemented")
// }
