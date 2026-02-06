package retry

import (
	"context"
	"math/rand"
	"time"
)

// RetryContext 执行fn回调函数最大次数n，每次执行间隔时间interval，num必须大于0否则panic
func RetryContext(ctx context.Context, n int, interval func() time.Duration, fn func() error) (err error) {
	if n <= 0 {
		panic("num must be > 0")
	}
	curCtx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
				if err = fn(); err == nil {
					return
				}
				if n--; n < 1 {
					return
				}
				time.Sleep(interval())
			}
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-curCtx.Done():
		return
	}
}

// Retry 执行fn回调函数最大次数n，每次执行间隔时间interval，num必须大于0否则panic
func Retry(num int, interval func() time.Duration, fn func() error) (err error) {
	return RetryContext(context.Background(), num, interval, fn)
}

// Interval 返回 d
func Interval(d time.Duration) func() time.Duration {
	return func() time.Duration {
		return d
	}
}

// RandInterval 生成一个基于min开区间，max闭区间的随机数*unit，如果min = max则直接返回min*unit
func RandInterval(min, max int, unit time.Duration) func() time.Duration {
	if min > max {
		max, min = min, max
	}
	return func() time.Duration {
		if min == max {
			return time.Duration(min) * unit
		}
		return time.Duration(rand.Intn(max-min)+min) * unit
	}
}
