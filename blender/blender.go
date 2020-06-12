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
	// set the step
	b.step = (b.step + numSteps) % period
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
	result.SetColor(b.getTransitionColor(cfv, cf))
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

// getTransitionColor calculates the max tranversal distance and calls getColorTransitionColor
func (b *Blender) getTransitionColor(transVal float32, colorFunc *transfunc.ColorFunc) imageColor.RGBA {
	// get the full transition distance if we don't already have it
	if colorFunc.TransDist <= 0 {
		b.getColorTransitionDistance(colorFunc)
	}
	// find the distance for the given transVal
	dist := int(math.Round(float64(transVal * float32(colorFunc.TransDist))))
	// call the function for the given TransType
	// and get the transition color
	return b.getColorTransitionColor(colorFunc, dist)
}

// getColorTransTypeFunc selects the transition function
func (b *Blender) getColorTransTypeFunc(transType transfunc.TransType) func(color1 imageColor.RGBA, color2 imageColor.RGBA, maxDist int) (color imageColor.RGBA, distTraveled int) {
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

// getColorTransitionDistance sets colorFunc.TransDist to the full transition distance
func (b *Blender) getColorTransitionDistance(colorFunc *transfunc.ColorFunc) {
	// if the colors are equal, return 0
	if colorFunc.Color1 == colorFunc.Color2 {
		colorFunc.TransDist = 0
		return
	}
	// walk the traversal and get the distance
	transTypeFunc := b.getColorTransTypeFunc(colorFunc.TransType)
	_, colorFunc.TransDist = transTypeFunc(colorFunc.Color1, colorFunc.Color2, math.MaxUint8*4)
}

// getColorTransitionColor retreives the transition function and gets the transition color
func (b *Blender) getColorTransitionColor(colorFunc *transfunc.ColorFunc, transDist int) imageColor.RGBA {
	transTypeFunc := b.getColorTransTypeFunc(colorFunc.TransType)
	color, _ := transTypeFunc(colorFunc.Color1, colorFunc.Color2, transDist)
	return color
}

// oneAtATimeColorTransition transitions between colors by changing only one component value at a time
func (b *Blender) oneAtATimeColorTransition(color1 imageColor.RGBA, color2 imageColor.RGBA, maxDist int) (resultingColor imageColor.RGBA, distance int) {
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
func (b *Blender) allAtOnceColorTransition(color1 imageColor.RGBA, color2 imageColor.RGBA, maxDist int) (color imageColor.RGBA, distance int) {
	return imageColor.RGBA{}, 0
}

// whiteColorTransition similar to allAtOnceColorTransition but transitions to white before transitioning to the target values
func (b *Blender) whiteColorTransition(color1 imageColor.RGBA, color2 imageColor.RGBA, maxDist int) (color imageColor.RGBA, distance int) {
	return imageColor.RGBA{}, 0
}

// blackColorTransition similar to allAtOnceColorTransition but transitions to black before transitioning to the target values
func (b *Blender) blackColorTransition(color1 imageColor.RGBA, color2 imageColor.RGBA, maxDist int) (color imageColor.RGBA, distance int) {
	return imageColor.RGBA{}, 0
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
