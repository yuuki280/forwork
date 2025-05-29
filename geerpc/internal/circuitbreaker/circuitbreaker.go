// internal/circuitbreaker/circuitbreaker.go
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// 熔断器状态
const (
	StateClosed   = iota // 关闭状态，允许请求通过
	StateOpen            // 开启状态，请求被拒绝
	StateHalfOpen        // 半开状态，允许有限请求通过以测试服务是否恢复
)

var (
	ErrOpenState       = errors.New("熔断器开启")
	ErrTooManyRequests = errors.New("请求过多")
)

// Stats 表示熔断器统计信息
type Stats struct {
	Requests             int64
	TotalSuccesses       int64
	TotalFailures        int64
	ConsecutiveSuccesses int64
	ConsecutiveFailures  int64
}

// CircuitBreaker 表示熔断器
type CircuitBreaker struct {
	name                  string
	mutex                 sync.Mutex
	state                 int
	timeout               time.Duration // 熔断器开启状态持续时间
	failureThreshold      int64         // 连续失败阈值
	successThreshold      int64         // 半开状态下连续成功阈值
	maxConcurrentRequests int64         // 最大并发请求数
	currentRequests       int64         // 当前并发请求数
	stats                 Stats
	lastStateChangeTime   time.Time
}

// NewCircuitBreaker 创建一个新的熔断器
func NewCircuitBreaker(name string, failureThreshold, successThreshold, maxConcurrentRequests int64, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:                  name,
		state:                 StateClosed,
		failureThreshold:      failureThreshold,
		successThreshold:      successThreshold,
		maxConcurrentRequests: maxConcurrentRequests,
		timeout:               timeout,
		lastStateChangeTime:   time.Now(),
	}
}

// AllowRequest 判断是否允许请求通过
func (cb *CircuitBreaker) AllowRequest() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// 超过超时时间，从开启状态转为半开状态
	if cb.state == StateOpen && time.Since(cb.lastStateChangeTime) > cb.timeout {
		cb.changeState(StateHalfOpen)
	}

	// 根据状态决定是否允许请求
	switch cb.state {
	case StateClosed:
		if cb.currentRequests >= cb.maxConcurrentRequests {
			return ErrTooManyRequests
		}
		cb.currentRequests++
		return nil
	case StateHalfOpen:
		if cb.currentRequests >= 1 { // 半开状态只允许一个请求
			return ErrTooManyRequests
		}
		cb.currentRequests++
		return nil
	default: // StateOpen
		return ErrOpenState
	}
}

// 记录请求成功
func (cb *CircuitBreaker) Success() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.currentRequests--
	cb.stats.Requests++
	cb.stats.TotalSuccesses++
	cb.stats.ConsecutiveSuccesses++
	cb.stats.ConsecutiveFailures = 0

	// 半开状态下连续成功次数达到阈值，转为关闭状态
	if cb.state == StateHalfOpen && cb.stats.ConsecutiveSuccesses >= cb.successThreshold {
		cb.changeState(StateClosed)
	}
}

// 记录请求失败
func (cb *CircuitBreaker) Failure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.currentRequests--
	cb.stats.Requests++
	cb.stats.TotalFailures++
	cb.stats.ConsecutiveFailures++
	cb.stats.ConsecutiveSuccesses = 0

	// 关闭状态下连续失败次数达到阈值，转为开启状态
	if cb.state == StateClosed && cb.stats.ConsecutiveFailures >= cb.failureThreshold {
		cb.changeState(StateOpen)
	}

	// 半开状态下出现失败，立即转为开启状态
	if cb.state == StateHalfOpen {
		cb.changeState(StateOpen)
	}
}

// 改变熔断器状态
func (cb *CircuitBreaker) changeState(newState int) {
	cb.state = newState
	cb.lastStateChangeTime = time.Now()

	// 状态变更时重置一些计数
	if newState == StateOpen || newState == StateClosed {
		cb.stats.ConsecutiveSuccesses = 0
		cb.stats.ConsecutiveFailures = 0
	}
}

// 获取当前状态
func (cb *CircuitBreaker) State() int {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.state
}

// 获取统计信息
func (cb *CircuitBreaker) Stats() Stats {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.stats
}

// 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.stats = Stats{}
	cb.state = StateClosed
	cb.lastStateChangeTime = time.Now()
}
