package publish

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

// PublisherMock implements Publisher
type PublisherMock struct {
	t minimock.Tester

	funcPublish          func(version string) (err error)
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

	return m
}

type mPublisherMockPublish struct {
	mock               *PublisherMock
	defaultExpectation *PublisherMockPublishExpectation
	expectations       []*PublisherMockPublishExpectation
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
	version string
}

// PublisherMockPublishResults contains results of the Publisher.Publish
type PublisherMockPublishResults struct {
	err error
}

// Expect sets up expected params for Publisher.Publish
func (m *mPublisherMockPublish) Expect(version string) *mPublisherMockPublish {
	if m.mock.funcPublish != nil {
		m.mock.t.Fatalf("PublisherMock.Publish mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &PublisherMockPublishExpectation{}
	}

	m.defaultExpectation.params = &PublisherMockPublishParams{version}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Publisher.Publish
func (m *mPublisherMockPublish) Return(err error) *PublisherMock {
	if m.mock.funcPublish != nil {
		m.mock.t.Fatalf("PublisherMock.Publish mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &PublisherMockPublishExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &PublisherMockPublishResults{err}
	return m.mock
}

//Set uses given function f to mock the Publisher.Publish method
func (m *mPublisherMockPublish) Set(f func(version string) (err error)) *PublisherMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Publisher.Publish method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Publisher.Publish method")
	}

	m.mock.funcPublish = f
	return m.mock
}

// When sets expectation for the Publisher.Publish which will trigger the result defined by the following
// Then helper
func (m *mPublisherMockPublish) When(version string) *PublisherMockPublishExpectation {
	if m.mock.funcPublish != nil {
		m.mock.t.Fatalf("PublisherMock.Publish mock is already set by Set")
	}

	expectation := &PublisherMockPublishExpectation{
		mock:   m.mock,
		params: &PublisherMockPublishParams{version},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Publisher.Publish return parameters for the expectation previously defined by the When method
func (e *PublisherMockPublishExpectation) Then(err error) *PublisherMock {
	e.results = &PublisherMockPublishResults{err}
	return e.mock
}

// Publish implements Publisher
func (m *PublisherMock) Publish(version string) (err error) {
	atomic.AddUint64(&m.beforePublishCounter, 1)
	defer atomic.AddUint64(&m.afterPublishCounter, 1)

	for _, e := range m.PublishMock.expectations {
		if minimock.Equal(*e.params, PublisherMockPublishParams{version}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if m.PublishMock.defaultExpectation != nil {
		atomic.AddUint64(&m.PublishMock.defaultExpectation.Counter, 1)
		want := m.PublishMock.defaultExpectation.params
		got := PublisherMockPublishParams{version}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("PublisherMock.Publish got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.PublishMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the PublisherMock.Publish")
		}
		return (*results).err
	}
	if m.funcPublish != nil {
		return m.funcPublish(version)
	}
	m.t.Fatalf("Unexpected call to PublisherMock.Publish. %v", version)
	return
}

// PublishAfterCounter returns a count of finished PublisherMock.Publish invocations
func (m *PublisherMock) PublishAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterPublishCounter)
}

// PublishBeforeCounter returns a count of PublisherMock.Publish invocations
func (m *PublisherMock) PublishBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforePublishCounter)
}

// MinimockPublishDone returns true if the count of the Publish invocations corresponds
// the number of defined expectations
func (m *PublisherMock) MinimockPublishDone() bool {
	for _, e := range m.PublishMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PublishMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPublish != nil && atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		return false
	}
	return true
}

// MinimockPublishInspect logs each unmet expectation
func (m *PublisherMock) MinimockPublishInspect() {
	for _, e := range m.PublishMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PublisherMock.Publish with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PublishMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterPublishCounter) < 1 {
		m.t.Errorf("Expected call to PublisherMock.Publish with params: %#v", *m.PublishMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPublish != nil && atomic.LoadUint64(&m.afterPublishCounter) < 1 {
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
func (m *PublisherMock) MinimockWait(timeout time.Duration) {
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

func (m *PublisherMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockPublishDone()
}
