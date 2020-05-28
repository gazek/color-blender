package transfunc

type transFuncer interface {
	GetFuncValue(stepNum int) float32
	GetFuncPeriod() int
}

type transFunc struct {
	Function   func(x float32) float32
	Period     int
	InputRange []float32 // left inclusive, right exclusive
	TransType  string    // "base", "white", "black"
}

func (f *transFunc) GetFuncPeriod() int {
	return f.Period
}

func (f *transFunc) GetFuncValue(stepNum int) float32 {
	// make sure the input range is valid
	if len(f.InputRange) != 2 {
		// set to zero value
		f.InputRange = []float32{0, 0}
	}
	// find the input value
	stepMin := stepNum % f.Period
	posInRange := float32(stepMin) / float32(f.Period)
	inputValue := posInRange*(f.InputRange[1]-f.InputRange[0]) + f.InputRange[0]
	// get the function value
	return f.Function(inputValue)
}
