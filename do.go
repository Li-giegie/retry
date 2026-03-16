package retry

import (
	"context"
	"math/rand"
	"time"
)

// DoContext 执行fn回调函数最大次数n，每次执行间隔时间interval，num <= 0无限次数重试
func DoContext(ctx context.Context, n int, interval func() time.Duration, fn func(ctx context.Context) error) (err error) {
	curCtx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	go func() {
		defer cancel()
		count := 0
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
				if err = fn(curCtx); err == nil {
					return
				}
				if count++; n > 0 && count >= n {
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

// Do 执行fn回调函数最大次数n，每次执行间隔时间interval，num <= 0无限次数重试
func Do(num int, interval func() time.Duration, fn func(ctx context.Context) error) (err error) {
	return DoContext(context.Background(), num, interval, fn)
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
