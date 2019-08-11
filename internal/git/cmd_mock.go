package git

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	"time"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// CmdMock implements Cmd
type CmdMock struct {
	t minimock.Tester

	funcRun          func(args []string, timeout time.Duration) (s1 string, err error)
	inspectFuncRun   func(args []string, timeout time.Duration)
	afterRunCounter  uint64
	beforeRunCounter uint64
	RunMock          mCmdMockRun
}

// NewCmdMock returns a mock for Cmd
func NewCmdMock(t minimock.Tester) *CmdMock {
	m := &CmdMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RunMock = mCmdMockRun{mock: m}
	m.RunMock.callArgs = []*CmdMockRunParams{}

	return m
}

type mCmdMockRun struct {
	mock               *CmdMock
	defaultExpectation *CmdMockRunExpectation
	expectations       []*CmdMockRunExpectation

	callArgs []*CmdMockRunParams
	mutex    sync.RWMutex
}

// CmdMockRunExpectation specifies expectation struct of the Cmd.Run
type CmdMockRunExpectation struct {
	mock    *CmdMock
	params  *CmdMockRunParams
	results *CmdMockRunResults
	Counter uint64
}

// CmdMockRunParams contains parameters of the Cmd.Run
type CmdMockRunParams struct {
	args    []string
	timeout time.Duration
}

// CmdMockRunResults contains results of the Cmd.Run
type CmdMockRunResults struct {
	s1  string
	err error
}

// Expect sets up expected params for Cmd.Run
func (mmRun *mCmdMockRun) Expect(args []string, timeout time.Duration) *mCmdMockRun {
	if mmRun.mock.funcRun != nil {
		mmRun.mock.t.Fatalf("CmdMock.Run mock is already set by Set")
	}

	if mmRun.defaultExpectation == nil {
		mmRun.defaultExpectation = &CmdMockRunExpectation{}
	}

	mmRun.defaultExpectation.params = &CmdMockRunParams{args, timeout}
	for _, e := range mmRun.expectations {
		if minimock.Equal(e.params, mmRun.defaultExpectation.params) {
			mmRun.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRun.defaultExpectation.params)
		}
	}

	return mmRun
}

// Inspect accepts an inspector function that has same arguments as the Cmd.Run
func (mmRun *mCmdMockRun) Inspect(f func(args []string, timeout time.Duration)) *mCmdMockRun {
	if mmRun.mock.inspectFuncRun != nil {
		mmRun.mock.t.Fatalf("Inspect function is already set for CmdMock.Run")
	}

	mmRun.mock.inspectFuncRun = f

	return mmRun
}

// Return sets up results that will be returned by Cmd.Run
func (mmRun *mCmdMockRun) Return(s1 string, err error) *CmdMock {
	if mmRun.mock.funcRun != nil {
		mmRun.mock.t.Fatalf("CmdMock.Run mock is already set by Set")
	}

	if mmRun.defaultExpectation == nil {
		mmRun.defaultExpectation = &CmdMockRunExpectation{mock: mmRun.mock}
	}
	mmRun.defaultExpectation.results = &CmdMockRunResults{s1, err}
	return mmRun.mock
}

//Set uses given function f to mock the Cmd.Run method
func (mmRun *mCmdMockRun) Set(f func(args []string, timeout time.Duration) (s1 string, err error)) *CmdMock {
	if mmRun.defaultExpectation != nil {
		mmRun.mock.t.Fatalf("Default expectation is already set for the Cmd.Run method")
	}

	if len(mmRun.expectations) > 0 {
		mmRun.mock.t.Fatalf("Some expectations are already set for the Cmd.Run method")
	}

	mmRun.mock.funcRun = f
	return mmRun.mock
}

// When sets expectation for the Cmd.Run which will trigger the result defined by the following
// Then helper
func (mmRun *mCmdMockRun) When(args []string, timeout time.Duration) *CmdMockRunExpectation {
	if mmRun.mock.funcRun != nil {
		mmRun.mock.t.Fatalf("CmdMock.Run mock is already set by Set")
	}

	expectation := &CmdMockRunExpectation{
		mock:   mmRun.mock,
		params: &CmdMockRunParams{args, timeout},
	}
	mmRun.expectations = append(mmRun.expectations, expectation)
	return expectation
}

// Then sets up Cmd.Run return parameters for the expectation previously defined by the When method
func (e *CmdMockRunExpectation) Then(s1 string, err error) *CmdMock {
	e.results = &CmdMockRunResults{s1, err}
	return e.mock
}

// Run implements Cmd
func (mmRun *CmdMock) Run(args []string, timeout time.Duration) (s1 string, err error) {
	mm_atomic.AddUint64(&mmRun.beforeRunCounter, 1)
	defer mm_atomic.AddUint64(&mmRun.afterRunCounter, 1)

	if mmRun.inspectFuncRun != nil {
		mmRun.inspectFuncRun(args, timeout)
	}

	mm_params := &CmdMockRunParams{args, timeout}

	// Record call args
	mmRun.RunMock.mutex.Lock()
	mmRun.RunMock.callArgs = append(mmRun.RunMock.callArgs, mm_params)
	mmRun.RunMock.mutex.Unlock()

	for _, e := range mmRun.RunMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.s1, e.results.err
		}
	}

	if mmRun.RunMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRun.RunMock.defaultExpectation.Counter, 1)
		mm_want := mmRun.RunMock.defaultExpectation.params
		mm_got := CmdMockRunParams{args, timeout}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRun.t.Errorf("CmdMock.Run got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRun.RunMock.defaultExpectation.results
		if mm_results == nil {
			mmRun.t.Fatal("No results are set for the CmdMock.Run")
		}
		return (*mm_results).s1, (*mm_results).err
	}
	if mmRun.funcRun != nil {
		return mmRun.funcRun(args, timeout)
	}
	mmRun.t.Fatalf("Unexpected call to CmdMock.Run. %v %v", args, timeout)
	return
}

// RunAfterCounter returns a count of finished CmdMock.Run invocations
func (mmRun *CmdMock) RunAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRun.afterRunCounter)
}

// RunBeforeCounter returns a count of CmdMock.Run invocations
func (mmRun *CmdMock) RunBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRun.beforeRunCounter)
}

// Calls returns a list of arguments used in each call to CmdMock.Run.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRun *mCmdMockRun) Calls() []*CmdMockRunParams {
	mmRun.mutex.RLock()

	argCopy := make([]*CmdMockRunParams, len(mmRun.callArgs))
	copy(argCopy, mmRun.callArgs)

	mmRun.mutex.RUnlock()

	return argCopy
}

// MinimockRunDone returns true if the count of the Run invocations corresponds
// the number of defined expectations
func (m *CmdMock) MinimockRunDone() bool {
	for _, e := range m.RunMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRun != nil && mm_atomic.LoadUint64(&m.afterRunCounter) < 1 {
		return false
	}
	return true
}

// MinimockRunInspect logs each unmet expectation
func (m *CmdMock) MinimockRunInspect() {
	for _, e := range m.RunMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CmdMock.Run with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunCounter) < 1 {
		if m.RunMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CmdMock.Run")
		} else {
			m.t.Errorf("Expected call to CmdMock.Run with params: %#v", *m.RunMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRun != nil && mm_atomic.LoadUint64(&m.afterRunCounter) < 1 {
		m.t.Error("Expected call to CmdMock.Run")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CmdMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockRunInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CmdMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *CmdMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRunDone()
}
