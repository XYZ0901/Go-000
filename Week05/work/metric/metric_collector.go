package metric

import (
	"sync"
	"time"
)

// MetricCollector represents the contract that all collectors must fulfill to gather circuit statistics.
// Implementations of this interface do not have to maintain locking around thier data stores so long as
// they are not modified outside of the hystrix context.
type MetricCollector interface {
	// Update accepts a set of metrics from a command execution for remote instrumentation
	Update(MetricResult)
	// Reset resets the internal counters and timers.
	Reset()
}
type DefaultMetricCollector struct {
	mutex *sync.RWMutex

	numRequests *Number
	errors      *Number

	successes *Number

	totalDuration *Timing
	runDuration   *Timing
}

func NewDefaultMetricCollector() MetricCollector {
	m := &DefaultMetricCollector{}
	m.mutex = &sync.RWMutex{}
	m.Reset()
	return m
}

// NumRequests returns the rolling number of requests
func (d *DefaultMetricCollector) NumRequests() *Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.numRequests
}

// Errors returns the rolling number of errors
func (d *DefaultMetricCollector) Errors() *Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.errors
}

// Successes returns the rolling number of successes
func (d *DefaultMetricCollector) Successes() *Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.successes
}

// TotalDuration returns the rolling total duration
func (d *DefaultMetricCollector) TotalDuration() *Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.totalDuration
}

// RunDuration returns the rolling run duration
func (d *DefaultMetricCollector) RunDuration() *Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.runDuration
}

func (d *DefaultMetricCollector) Update(r MetricResult) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	d.numRequests.Increment(r.Attempts)
	d.errors.Increment(r.Errors)
	d.successes.Increment(r.Successes)

	d.totalDuration.Add(r.TotalDuration)
	d.runDuration.Add(r.RunDuration)
}

// Reset resets all metrics in this collector to 0.
func (d *DefaultMetricCollector) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.numRequests = NewNumber()
	d.errors = NewNumber()
	d.successes = NewNumber()

	d.totalDuration = NewTiming()
	d.runDuration = NewTiming()
}

type MetricResult struct {
	Attempts  float64
	Errors    float64
	Successes float64

	TotalDuration time.Duration
	RunDuration   time.Duration
}
