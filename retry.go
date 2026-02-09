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

// RandInterval 生成一个基于min开区间，max闭区间的随机数 * base，如果min = max则直接返回min*unit
func RandInterval(min, max int, base time.Duration) func() time.Duration {
	if min > max {
		max, min = min, max
	}
	return func() time.Duration {
		if min == max {
			return time.Duration(min) * base
		}
		return time.Duration(rand.Intn(max-min)+min) * base
	}
}

// RandExponentialBackoff 指数退避值+随机值
func RandExponentialBackoff(max int, base time.Duration, rMin, rMax int, rBase time.Duration) func() time.Duration {
	num := 1
	fn := RandInterval(rMin, rMax, rBase)
	return func() time.Duration {
		d := time.Duration(num) * base
		if num <<= 1; num > max && max > -1 {
			num = max
		}
		return d + fn()
	}
}

// ExponentialBackoff 指数退避，max计算次幂的最大值，base 时间单位
func ExponentialBackoff(max int, base time.Duration) func() time.Duration {
	num := 1
	return func() time.Duration {
		d := time.Duration(num) * base
		if num <<= 1; num > max && max > -1 {
			num = max
		}
		return d
	}
}
