package transfunc

type transFuncer interface {
	GetFuncValue(stepNum int) float64
	GetFuncPeriod() int
}

type transFunc struct {
	Function   func(x float64) float64
	Period     int
	InputRange []float64 // left inclusive, right exclusive
}

func (f *transFunc) GetFuncPeriod() int {
	return f.Period
}

func (f *transFunc) GetFuncValue(stepNum int) float64 {
	// make sure the input range is valid
	if len(f.InputRange) != 2 {
		// set to zero value
		f.InputRange = []float64{0, 0}
	}
	// find the input value
	stepMin := stepNum % f.Period
	posInRange := float64(stepMin) / float64(f.Period)
	inputValue := posInRange*(f.InputRange[1]-f.InputRange[0]) + f.InputRange[0]
	// get the function value
	return f.Function(inputValue)
}
