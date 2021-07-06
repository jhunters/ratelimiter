// Copyright 2021 The baidu Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	errStopped = errors.New("bad operation cause rate limiter is stopped")
)

// token
type token struct {
	ctime time.Time
}

// RateLimiter rate limiter struct
type RateLimiter struct {
	tokenBucket chan token
	// QPS set
	quotePerSeconds uint32
	interval        time.Duration // the interval time to produce a request token
	ticker          *time.Ticker
	stopChannel     chan bool // stop controller
	stop            bool

	lock          sync.Mutex
	lastTokenTime time.Time
}

// NewRateLimiter return a new rate limiter
func NewRateLimiter(quotePerSeconds uint32) (*RateLimiter, error) {
	if quotePerSeconds == 0 {
		return nil, fmt.Errorf("quotePerSeconds should be large than zero")
	}

	interval := time.Duration(int64(time.Second) / int64(quotePerSeconds))

	ticker := time.NewTicker(interval)
	limiter := &RateLimiter{
		quotePerSeconds: quotePerSeconds,
		tokenBucket:     make(chan token, quotePerSeconds),
		interval:        interval,
		ticker:          ticker,
		stopChannel:     make(chan bool),
		stop:            false,
	}
	go limiter.start()
	return limiter, nil
}

// start rate limiter
func (limiter *RateLimiter) start() {
	for {
		select {
		case <-limiter.ticker.C:
			now := time.Now()
			limiter.tokenBucket <- token{ctime: now}
			limiter.lastTokenTime = now
		case <-limiter.stopChannel:
			limiter.ticker.Stop()
			limiter.stop = true
			return
		}
	}
}

// Acquire Acquires a single permit from this {@code RateLimiter}, blocking until the
// request can be granted. Tells the amount of time slept, if any.
func (limiter *RateLimiter) Acquire() (time.Duration, error) {
	return limiter.AcquireBatch(uint32(1))
}

// AcquireBatch Acquires the specified permit count from this {@code RateLimiter}, blocking until the
// request can be granted. Tells the amount of time slept, if any.
func (limiter *RateLimiter) AcquireBatch(permits uint32) (time.Duration, error) {
	if permits == 0 {
		return time.Duration(0), (fmt.Errorf("permits should be large than zero"))
	}
	if limiter.stop {
		return time.Duration(0), errStopped
	}
	now := time.Now()
	r := make(chan bool, 1)
	limiter.lock.Lock()
	defer limiter.lock.Unlock()
	limiter.acquireBatch(context.Background(), permits, r)
	return time.Since(now), nil
}

func (limiter *RateLimiter) acquireBatch(context context.Context, permits uint32, result chan<- bool) {
	for i := 0; i < int(permits); i++ {
		select {
		case <-context.Done():
			// time out and break
			result <- false
			return
		default:
			<-limiter.tokenBucket
		}

	}
	result <- true
}

func (limiter *RateLimiter) acquireBatchWithLock(context context.Context, permits uint32, result chan<- bool) {
	limiter.lock.Lock()
	defer limiter.lock.Unlock()
	limiter.acquireBatch(context, permits, result)
}

// AcquireBatchWithTimeout It return true if acquires a permit from this {@link RateLimiter}
// if it can be acquired with in timeout. otherwise return false
func (limiter *RateLimiter) TryAcquireWithTimeout(timeout time.Duration) (bool, error) {
	return limiter.TryAcquireBatchWithTimeout(1, timeout)

}

// AcquireBatchWithTimeout It return true if acquires specified permit count from this {@link RateLimiter}
// if it can be acquired with in timeout. otherwise return false
func (limiter *RateLimiter) TryAcquireBatchWithTimeout(permits uint32, timeout time.Duration) (bool, error) {
	if limiter.stop {
		return false, errStopped
	}
	// check can acquire
	if len(limiter.tokenBucket)+(int(timeout)/int(limiter.interval)) < int(permits) {
		return false, nil
	}

	r := make(chan bool, 1)
	context, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	limiter.acquireBatchWithLock(context, 1, r) // bug  超时后，超时前的已被取走token

	return <-r, nil

}

// TryAcquire Acquires a permit from this {@link RateLimiter} if it can be acquired immediately without delay
func (limiter *RateLimiter) TryAcquire() (bool, error) {
	return limiter.TryAcquireBatch(1)
}

// TryAcquireBatch Acquires permits from this {@link RateLimiter} if it can be acquired immediately without delay.
func (limiter *RateLimiter) TryAcquireBatch(permits uint32) (bool, error) {
	if limiter.stop {
		return false, errStopped
	}
	limiter.lock.Lock()
	defer limiter.lock.Unlock()
	r := len(limiter.tokenBucket) < int(permits)
	if r {
		return false, nil
	}
	ch := make(chan bool, 1)
	limiter.acquireBatch(context.Background(), permits, ch)

	return true, nil
}

// Stop stop the rate limiter
func (limiter *RateLimiter) Stop() {
	limiter.stopChannel <- false
}
