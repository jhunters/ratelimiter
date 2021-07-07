# ratelimiter
pure golang implments for token bucket

[![Go Report Card](https://goreportcard.com/badge/github.com/jhunters/ratelimiter?style=flat-square)](https://goreportcard.com/report/github.com/jhunters/ratelimiter)
[![Build Status](https://travis-ci.com/jhunters/ratelimiter.svg?branch=main&status=started)](https://travis-ci.com/jhunters/ratelimiter)
[![codecov](https://codecov.io/gh/jhunters/ratelimiter/branch/main/graph/badge.svg?token=ATQhFv91YP)](https://codecov.io/gh/jhunters/ratelimiter)
[![Releases](https://img.shields.io/github/release/jhunters/ratelimiter/all.svg?style=flat-square)](https://github.com/jhunters/ratelimiter/releases)
[![Go Reference](https://golang.com.cn/badge/github.com/jhunters/ratelimiter.svg)](https://golang.com.cn/github.com/jhunters/ratelimiter)
[![LICENSE](https://img.shields.io/github/license/jhunters/ratelimiter.svg?style=flat-square)](https://github.com/jhunters/ratelimiter/blob/master/LICENSE)

## Usage
### Installing 

To start using ratelimiter, install Go and run `go get`:

```sh
$ go get github.com/jhunters/ratelimiter
```

### base method

create RateLimiter

```go
// 初始化令牌桶, 控制并发 100 / 秒
limiter, err := ratelimiter.NewRateLimiter(100)
if err != nil {
    panic(err)
}

```

serval ways to acquire tokens

```go
// acquire a token block unitl it reached
cost, err := limiter.Acquire()
fmt.Prinlnt(cost) // print cost time

```

```go
// acquire tokens block unitl it reached
cost, err := limiter.AcquireBatch(10)
```

```go
// try acquire tokens if ready or return false immediately
permit, err := limiter.TryAcquireBatch(20)
```

```go

// try acquire a token with time out
permit, err := limiter.TryAcquireWithTimeout(100 * time.Millisecond)

```

close RateLimiter

```go
limiter.Stop()
```