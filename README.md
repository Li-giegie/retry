# retry
A simple go retry method function package

```go
package main

import (
	"context"
	"errors"
	"github.com/Li-giegie/retry"
	"log"
	"time"
)

func Todo() (any, error) {
	log.Println("todo......")
	return nil, errors.New("err: xxxxxx")
}

func main() {
	Do()
	DoContest()
	DoContestExponentialBackoff()
	DoContestRandExponentialBackoff()
}

// Do 基础重试，最大重试次数3次，随机1~2s（最大值为3，但这里为开区间）执行失败间隔时间。
func Do() {
	// 执行函数fn，出错最大重试次数为3
	var result any
	err := retry.Do(3, retry.RandInterval(1, 3, time.Second), func(ctx context.Context) (err error) {
		result, err = Todo()
		return
	})
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(result)
}

// DoContest 带上下文重试，不限重试次数，使用上上下文超时限定执行时间，随机1~2s（最大值为3，但这里为开区间）执行失败间隔时间。
func DoContest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var result any
	err := retry.DoContext(ctx, 0, retry.RandInterval(1, 3, time.Second), func(ctx context.Context) (err error) {
		result, err = Todo()
		return
	})
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(result)
}

// DoContestExponentialBackoff 带上下文重试，不限重试次数，使用上上下文超时限定执行时间，指数退避执行失败间隔时间（例如16：1、2、4、8、16...16）。
func DoContestExponentialBackoff() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var result any
	err := retry.DoContext(ctx, 0, retry.ExponentialBackoff(16, time.Second), func(ctx context.Context) (err error) {
		result, err = Todo()
		return
	})
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(result)
}

// DoContestRandExponentialBackoff 带上下文重试，不限重试次数，使用上上下文超时限定执行时间，指数退避+随机震荡时间，执行失败间隔时间（例如16：1s + 190ms、2s + 150ms...16s + 800ms）。
func DoContestRandExponentialBackoff() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var result any
	err := retry.DoContext(ctx, 0, retry.RandExponentialBackoff(16, time.Second, 100, 1000, time.Millisecond), func(ctx context.Context) (err error) {
		result, err = Todo()
		return
	})
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(result)
}
```