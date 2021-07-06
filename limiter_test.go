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
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestRateLimiterCreate test NewRateLimiter
func TestRateLimiterCreate(t *testing.T) {

	Convey("test NewRateLimiter function", t, func() {

		Convey("NewRateLimiter with error paramters", func() {

			limiter, err := NewRateLimiter(0)
			So(err, ShouldNotBeNil)
			So(limiter, ShouldBeNil)
		})

		Convey("NewRateLimiter success", func() {

			limiter, err := NewRateLimiter(1)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			defer limiter.Stop()

			So(limiter.quotePerSeconds, ShouldEqual, 1)
			So(limiter.tokenBucket, ShouldNotBeNil)
			So(limiter.interval, ShouldEqual, time.Second)

		})
	})
}

// TestAcquireOper test acquire operation
func TestAcquireOper(t *testing.T) {

	Convey("test acquire", t, func() {

		Convey("test default acquire ", func() {

			limiter, err := NewRateLimiter(10)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			defer limiter.Stop()

			So(10*limiter.interval, ShouldEqual, time.Second)
			cost, err := limiter.Acquire()
			So(cost, ShouldBeLessThan, 150*time.Millisecond)
			So(err, ShouldBeNil)

		})

		Convey("test AcquireBatch", func() {
			limiter, err := NewRateLimiter(100)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			defer limiter.Stop()

			So(100*limiter.interval, ShouldEqual, time.Second)

			cost, err := limiter.AcquireBatch(10)
			So(cost, ShouldBeLessThan, 180*time.Millisecond) // for some case delayed
			So(err, ShouldBeNil)

			Convey("test delayed acquire case", func() {
				// if rate limiter
				time.Sleep(2 * time.Second)
				cost, err = limiter.AcquireBatch(300)
				So(cost, ShouldBeGreaterThan, 1*time.Second)
				So(err, ShouldBeNil)

			})
		})
	})
}

// TestTryAcquireOper test case for TryAcquire functions
func TestTryAcquireOper(t *testing.T) {

	Convey("test TryAcquire", t, func() {

		Convey("test default TryAcquire", func() {
			limiter, err := NewRateLimiter(2)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			defer limiter.Stop()
			// fetch the early token
			limiter.Acquire()

			permit, err := limiter.TryAcquire()
			So(permit, ShouldBeFalse)
			ShouldBeNil(err)

			time.Sleep(1 * time.Second)
			permit, err = limiter.TryAcquire()
			So(err, ShouldBeNil)
			So(permit, ShouldBeTrue)
		})

		Convey("test TryAcquireBatch", func() {
			limiter, err := NewRateLimiter(15)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			defer limiter.Stop()

			permit, err := limiter.TryAcquireBatch(20)
			So(permit, ShouldBeFalse)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)
			fmt.Println(len(limiter.tokenBucket))
			permit, err = limiter.TryAcquireBatch(10)
			So(permit, ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("test TryAcquireWithTimeout", func() {
			limiter, err := NewRateLimiter(2)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			defer limiter.Stop()

			now := time.Now()
			permit, err := limiter.TryAcquireWithTimeout(100 * time.Millisecond)
			So(permit, ShouldBeFalse)
			So(time.Since(now), ShouldBeLessThan, 110*time.Millisecond)
			So(err, ShouldBeNil)

			permit, err = limiter.TryAcquireBatchWithTimeout(2, 2*time.Second)
			So(permit, ShouldBeTrue)
			So(err, ShouldBeNil)

		})
	})
}

// TestAcquireOperOnStopped test acquire action after stopped
func TestAcquireOperOnStopped(t *testing.T) {
	Convey("TestAcquireOperOnStopped", t, func() {
		Convey("test acquire after stopped", func() {
			limiter, err := NewRateLimiter(2)
			So(err, ShouldBeNil)
			So(limiter, ShouldNotBeNil)
			So(limiter.stop, ShouldBeFalse)
			limiter.Stop()

			time.Sleep(100 * time.Millisecond) // wait a while due to stop is async way

			So(limiter.stop, ShouldBeTrue)
			elasped, err := limiter.Acquire()
			So(err, ShouldEqual, errStopped)
			So(elasped, ShouldEqual, time.Duration(0))
		})
	})
}
