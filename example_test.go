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

import "fmt"

// ExampleRateLimiter example for RateLimiter
func ExampleRateLimiter() {

	limiter, err := NewRateLimiter(10)
	if err != nil {
		fmt.Println(err)
		return
	}

	// acquire one token
	cost, err := limiter.Acquire()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cost) // print cose time in milliseconds

	ok, err := limiter.TryAcquire() // if has no token return false immediately
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ok) // true
}
