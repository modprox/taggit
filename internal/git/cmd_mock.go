package git

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

// CmdMock implements Cmd
type CmdMock struct {
	t minimock.Tester

	funcRun          func(args []string, timeout time.Duration) (s1 string, err error)
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

	return m
}

type mCmdMockRun struct {
	mock               *CmdMock
	defaultExpectation *CmdMockRunExpectation
	expectations       []*CmdMockRunExpectation
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
func (m *mCmdMockRun) Expect(args []string, timeout time.Duration) *mCmdMockRun {
	if m.mock.funcRun != nil {
		m.mock.t.Fatalf("CmdMock.Run mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CmdMockRunExpectation{}
	}

	m.defaultExpectation.params = &CmdMockRunParams{args, timeout}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Cmd.Run
func (m *mCmdMockRun) Return(s1 string, err error) *CmdMock {
	if m.mock.funcRun != nil {
		m.mock.t.Fatalf("CmdMock.Run mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &CmdMockRunExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &CmdMockRunResults{s1, err}
	return m.mock
}

//Set uses given function f to mock the Cmd.Run method
func (m *mCmdMockRun) Set(f func(args []string, timeout time.Duration) (s1 string, err error)) *CmdMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Cmd.Run method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Cmd.Run method")
	}

	m.mock.funcRun = f
	return m.mock
}

// When sets expectation for the Cmd.Run which will trigger the result defined by the following
// Then helper
func (m *mCmdMockRun) When(args []string, timeout time.Duration) *CmdMockRunExpectation {
	if m.mock.funcRun != nil {
		m.mock.t.Fatalf("CmdMock.Run mock is already set by Set")
	}

	expectation := &CmdMockRunExpectation{
		mock:   m.mock,
		params: &CmdMockRunParams{args, timeout},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Cmd.Run return parameters for the expectation previously defined by the When method
func (e *CmdMockRunExpectation) Then(s1 string, err error) *CmdMock {
	e.results = &CmdMockRunResults{s1, err}
	return e.mock
}

// Run implements Cmd
func (m *CmdMock) Run(args []string, timeout time.Duration) (s1 string, err error) {
	atomic.AddUint64(&m.beforeRunCounter, 1)
	defer atomic.AddUint64(&m.afterRunCounter, 1)

	for _, e := range m.RunMock.expectations {
		if minimock.Equal(*e.params, CmdMockRunParams{args, timeout}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.s1, e.results.err
		}
	}

	if m.RunMock.defaultExpectation != nil {
		atomic.AddUint64(&m.RunMock.defaultExpectation.Counter, 1)
		want := m.RunMock.defaultExpectation.params
		got := CmdMockRunParams{args, timeout}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("CmdMock.Run got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.RunMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the CmdMock.Run")
		}
		return (*results).s1, (*results).err
	}
	if m.funcRun != nil {
		return m.funcRun(args, timeout)
	}
	m.t.Fatalf("Unexpected call to CmdMock.Run. %v %v", args, timeout)
	return
}

// RunAfterCounter returns a count of finished CmdMock.Run invocations
func (m *CmdMock) RunAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterRunCounter)
}

// RunBeforeCounter returns a count of CmdMock.Run invocations
func (m *CmdMock) RunBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeRunCounter)
}

// MinimockRunDone returns true if the count of the Run invocations corresponds
// the number of defined expectations
func (m *CmdMock) MinimockRunDone() bool {
	for _, e := range m.RunMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRunCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRun != nil && atomic.LoadUint64(&m.afterRunCounter) < 1 {
		return false
	}
	return true
}

// MinimockRunInspect logs each unmet expectation
func (m *CmdMock) MinimockRunInspect() {
	for _, e := range m.RunMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CmdMock.Run with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterRunCounter) < 1 {
		m.t.Errorf("Expected call to CmdMock.Run with params: %#v", *m.RunMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRun != nil && atomic.LoadUint64(&m.afterRunCounter) < 1 {
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
func (m *CmdMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (m *CmdMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRunDone()
}
