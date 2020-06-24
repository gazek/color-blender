package transfunc

type transFuncSlice struct {
	funcs  []transFuncer
	period int
}

// SetFuncs overwrite the current function slice with a new one
func (s *transFuncSlice) SetFuncs(funcs []transFuncer) {
	// store the value
	s.funcs = funcs
	// calculate the period
	s.setPeriod()
}

// AppendFunc appends a transFuncer to the slice
func (s *transFuncSlice) AppendFunc(f transFuncer) {
	// append the transFuncer
	s.funcs = append(s.funcs, f)
	// calculate the period
	if s.period == 0 {
		s.setPeriod()
	} else {
		s.period = s.period + f.GetFuncPeriod()
	}
}

func (s *transFuncSlice) setPeriod() {
	// get full period
	var period int
	for f := range s.funcs {
		period += s.funcs[f].GetFuncPeriod()
	}
	// store it
	s.period = period
}

// GetFuncValue returns the value of the function at the given step
func (s *transFuncSlice) GetFuncValue(stepNum int) (float64, transFuncer) {
	if len(s.funcs) == 0 {
		return 0, nil
	}
	// mod the step number
	minStep := stepNum % s.period
	// find the function index
	index, localStep := s.getFunctionIndex(minStep)
	// get the function value
	return s.funcs[index].GetFuncValue(localStep), s.funcs[index]
}

func (s *transFuncSlice) getFunctionIndex(stepNum int) (index int, localStep int) {
	// mod the step number
	minStep := stepNum % s.period
	// find the function index
	boundary := 0
	for f := range s.funcs {
		if minStep <= boundary+s.funcs[f].GetFuncPeriod() {
			localStep = minStep - boundary
			break
		}
		index++
		boundary += s.funcs[f].GetFuncPeriod()
	}
	return index, localStep
}

func (s *transFuncSlice) GetPeriod() int {
	return s.period
}
