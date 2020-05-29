package transfunc

import (
	"testing"
)

func TestSetFuncs(t *testing.T) {
	funcs := []transFuncer{
		&transFunc{Period: 1},
		&transFunc{Period: 2},
		&transFunc{Period: 3},
	}
	want := 6
	s := transFuncSlice{}
	s.SetFuncs(funcs)
	// check that the array was stored
	if s.funcs == nil {
		t.Errorf("Wanted: %v, found: %v", funcs, s.funcs)
	}
	// check that the period was set
	if s.period != want {
		t.Errorf("Wanted: %v, found: %v", want, s.period)
	}
}

func TestAppendFunc(t *testing.T) {
	funcs := []transFuncer{
		&transFunc{Period: 1},
		&transFunc{Period: 2},
	}
	newfunc := &transFunc{Period: 3}
	want := 6
	s := transFuncSlice{funcs: funcs}
	s.AppendFunc(newfunc)
	// check that the func was appended was stored
	if len(s.funcs) != len(funcs)+1 {
		t.Errorf("Wanted: %v, found: %v", len(funcs), len(s.funcs))
	}
	// check that the new period was calculated
	if s.period != want {
		t.Errorf("Wanted: %v, found: %v", want, s.period)
	}
}

func TestAppendFuncWithSetFuncs(t *testing.T) {
	funcs := []transFuncer{
		&transFunc{Period: 1},
		&transFunc{Period: 2},
	}
	newfunc := &transFunc{Period: 3}
	want := 6
	s := transFuncSlice{funcs: funcs}
	s.SetFuncs(funcs)
	s.AppendFunc(newfunc)
	// check that the func was appended was stored
	if len(s.funcs) != len(funcs)+1 {
		t.Errorf("Wanted: %v, found: %v", len(funcs), len(s.funcs))
	}
	// check that the new period was calculated
	if s.period != want {
		t.Errorf("Wanted: %v, found: %v", want, s.period)
	}
}

func TestGetFunctionIndex(t *testing.T) {
	tests := []struct {
		periods   []int
		stepNum   int
		index     int
		localStep int
	}{
		{[]int{5, 15, 10}, 7, 1, 2},
		{[]int{5, 15, 10}, 33, 0, 3},
		{[]int{5, 15, 10}, 5, 0, 5},
		{[]int{5, 15, 10}, 20, 1, 15},
	}

	for test := range tests {
		s := transFuncSlice{}
		for f := range tests[test].periods {
			s.AppendFunc(&transFunc{Period: tests[test].periods[f]})
		}
		index, localStep := s.getFunctionIndex(tests[test].stepNum)
		if index != tests[test].index {
			t.Errorf("index Wanted %v, got: %v", tests[test].index, index)
		}
		if localStep != tests[test].localStep {
			t.Errorf("localStep Wanted %v, got: %v", tests[test].localStep, localStep)
		}
	}
}

func TestGetFunctionValue(t *testing.T) {
	tests := []struct {
		periods []int
		stepNum int
		want    float32
	}{
		{[]int{5, 15, 10}, 7, 15},
		{[]int{5, 15, 10}, 33, 5},
		{[]int{5, 15, 10}, 5, 5},
		{[]int{5, 15, 10}, 20, 15},
	}

	for test := range tests {
		s := transFuncSlice{}
		for f := range tests[test].periods {
			returnValue := float32(tests[test].periods[f])
			s.AppendFunc(&transFunc{
				Period:   tests[test].periods[f],
				Function: func(x float32) float32 { return returnValue },
			})
		}
		result, _ := s.GetFuncValue(tests[test].stepNum)
		if result != tests[test].want {
			t.Errorf("index Wanted %v, got: %v", tests[test].want, result)
		}
	}
}

func TestGetPeriod(t *testing.T) {
	period := 42
	s := transFuncSlice{}
	s.AppendFunc(&transFunc{Period: period})
	// check that the array was stored
	if result := s.GetPeriod(); result != period {
		t.Errorf("Wanted: %v, found: %v", period, result)
	}
}
