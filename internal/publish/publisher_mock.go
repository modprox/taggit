package publish

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"gophers.dev/pkgs/semantic"
)

// PublisherMock implements Publisher
type PublisherMock struct {
	t minimock.Tester

	funcPublish          func(t1 semantic.Tag) (err error)
	inspectFuncPublish   func(t1 semantic.Tag)
	afterPublishCounter  uint64
	beforePublishCounter uint64
	PublishMock          mPublisherMockPublish
}

// NewPublisherMock returns a mock for Publisher
func NewPublisherMock(t minimock.Tester) *PublisherMock {
	m := &PublisherMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PublishMock = mPublisherMockPublish{mock: m}
	m.PublishMock.callArgs = []*PublisherMockPublishParams{}

	return m
}

type mPublisherMockPublish struct {
	mock               *PublisherMock
	defaultExpectation *PublisherMockPublishExpectation
	expectations       []*PublisherMockPublishExpectation

	callArgs []*PublisherMockPublishParams
	mutex    sync.RWMutex
}

// PublisherMockPublishExpectation specifies expectation struct of the Publisher.Publish
type PublisherMockPublishExpectation struct {
	mock    *PublisherMock
	params  *PublisherMockPublishParams
	results *PublisherMockPublishResults
	Counter uint64
}

// PublisherMockPublishParams contains parameters of the Publisher.Publish
type PublisherMockPublishParams struct {
	t1 semantic.Tag
}

// PublisherMockPublishResults contains results of the Publisher.Publish
type PublisherMockPublishResults struct {
	err error
}

// Expect sets up expected params for Publisher.Publish
func (mmPublish *mPublisherMockPublish) Expect(t1 semantic.Tag) *mPublisherMockPublish {
	if mmPublish.mock.funcPublish != nil {
		mmPublish.mock.t.Fatalf("PublisherMock.Publish mock is already set by Set")
	}

	if mmPublish.defaultExpectation == nil {
		mmPublish.defaultExpectation = &PublisherMockPublishExpectation{}
	}

	mmPublish.defaultExpectation.params = &PublisherMockPublishParams{t1}
	for _, e := range mmPublish.expectations {
		if minimock.Equal(e.params, mmPublish.defaultExpectation.params) {
			mmPublish.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmPublish.defaultExpectation.params)
		}
	}

	return mmPublish
}

// Inspect accepts an inspector function that has same arguments as the Publisher.Publish
func (mmPublish *mPublisherMockPublish) Inspect(f func(t1 semantic.Tag)) *mPublisherMockPublish {
	if mmPublish.mock.inspectFuncPublish != nil {
		mmPublish.mock.t.Fatalf("Inspect function is already set for PublisherMock.Publish")
	}

	mmPublish.mock.inspectFuncPublish = f

	return mmPublish
}

// Return sets up results that will be returned by Publisher.Publish
func (mmPublish *mPublisherMockPublish) Return(err error) *PublisherMock {
	if mmPublish.mock.funcPublish != nil {
		mmPublish.mock.t.Fatalf("PublisherMock.Publish mock is already set by Set")
	}

	if mmPublish.defaultExpectation == nil {
		mmPublish.defaultExpectation = &PublisherMockPublishExpectation{mock: mmPublish.mock}
	}
	mmPublish.defaultExpectation.results = &PublisherMockPublishResults{err}
	return mmPublish.mock
}

//Set uses given function f to mock the Publisher.Publish method
func (mmPublish *mPublisherMockPublish) Set(f func(t1 semantic.Tag) (err error)) *PublisherMock {
	if mmPublish.defaultExpectation != nil {
		mmPublish.mock.t.Fatalf("Default expectation is already set for the Publisher.Publish method")
	}

	if len(mmPublish.expectations) > 0 {
		mmPublish.mock.t.Fatalf("Some expectations are already set for the Publisher.Publish method")
	}

	mmPublish.mock.funcPublish = f
	return mmPublish.mock
}

// When sets expectation for the Publisher.Publish which will trigger the result defined by the following
// Then helper
func (mmPublish *mPublisherMockPublish) When(t1 semantic.Tag) *PublisherMockPublishExpectation {
	if mmPublish.mock.funcPublish != nil {
		mmPublish.mock.t.Fatalf("PublisherMock.Publish mock is already set by Set")
	}

	expectation := &PublisherMockPublishExpectation{
		mock:   mmPublish.mock,
		params: &PublisherMockPublishParams{t1},
	}
	mmPublish.expectations = append(mmPublish.expectations, expectation)
	return expectation
}

// Then sets up Publisher.Publish return parameters for the expectation previously defined by the When method
func (e *PublisherMockPublishExpectation) Then(err error) *PublisherMock {
	e.results = &PublisherMockPublishResults{err}
	return e.mock
}

// Publish implements Publisher
func (mmPublish *PublisherMock) Publish(t1 semantic.Tag) (err error) {
	mm_atomic.AddUint64(&mmPublish.beforePublishCounter, 1)
	defer mm_atomic.AddUint64(&mmPublish.afterPublishCounter, 1)

	if mmPublish.inspectFuncPublish != nil {
		mmPublish.inspectFuncPublish(t1)
	}

	mm_params := &PublisherMockPublishParams{t1}

	// Record call args
	mmPublish.PublishMock.mutex.Lock()
	mmPublish.PublishMock.callArgs = append(mmPublish.PublishMock.callArgs, mm_params)
	mmPublish.PublishMock.mutex.Unlock()

	for _, e := range mmPublish.PublishMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmPublish.PublishMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmPublish.PublishMock.defaultExpectation.Counter, 1)
		mm_want := mmPublish.PublishMock.defaultExpectation.params
		mm_got := PublisherMockPublishParams{t1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmPublish.t.Errorf("PublisherMock.Publish got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmPublish.PublishMock.defaultExpectation.results
		if mm_results == nil {
			mmPublish.t.Fatal("No results are set for the PublisherMock.Publish")
		}
		return (*mm_results).err
	}
	if mmPublish.funcPublish != nil {
		return mmPublish.funcPublish(t1)
	}
	mmPublish.t.Fatalf("Unexpected call to PublisherMock.Publish. %v", t1)
	return
}

// PublishAfterCounter returns a count of finished PublisherMock.Publish invocations
func (mmPublish *PublisherMock) PublishAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmPublish.afterPublishCounter)
}

// PublishBeforeCounter returns a count of PublisherMock.Publish invocations
func (mmPublish *PublisherMock) PublishBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmPublish.beforePublishCounter)
}

// Calls returns a list of arguments used in each call to PublisherMock.Publish.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmPublish *mPublisherMockPublish) Calls() []*PublisherMockPublishParams {
	mmPublish.mutex.RLock()

	argCopy := make([]*PublisherMockPublishParams, len(mmPublish.callArgs))
	copy(argCopy, mmPublish.callArgs)

	mmPublish.mutex.RUnlock()

	return argCopy
}

// MinimockPublishDone returns true if the count of the Publish invocations corresponds
// the number of defined expectations
func (m *PublisherMock) MinimockPublishDone() bool {
	for _, e := range m.PublishMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PublishMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPublish != nil && mm_atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		return false
	}
	return true
}

// MinimockPublishInspect logs each unmet expectation
func (m *PublisherMock) MinimockPublishInspect() {
	for _, e := range m.PublishMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PublisherMock.Publish with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PublishMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		if m.PublishMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PublisherMock.Publish")
		} else {
			m.t.Errorf("Expected call to PublisherMock.Publish with params: %#v", *m.PublishMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPublish != nil && mm_atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		m.t.Error("Expected call to PublisherMock.Publish")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PublisherMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockPublishInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PublisherMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *PublisherMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockPublishDone()
}
