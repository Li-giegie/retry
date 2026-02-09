# retry
A simple go retry method function package

```go
package retry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	err := Retry(3, RandInterval(1, 3, time.Second), func() error {
		n := rand.Int()
		log.Println("n", n)
		if n%2 == 0 {
			return errors.New("error")
		}
		return nil
	})
	fmt.Println(err)
}

func TestRetryContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	err := RetryContext(ctx, 10, RandInterval(1, 3, time.Second), func() error {
		log.Println("running......")
		return errors.New("error")
	})
	fmt.Println(err)
	time.Sleep(time.Second * 6)
}

func TestExponentialBackoff(t *testing.T) {
	// 重试6次，时间间隔为指数退幂最大休眠时间为16秒（1、2、4、8、16、16......）
	err := Retry(6, ExponentialBackoff(16, time.Second), func() error {
		log.Println("running......")
		return errors.New("error")
	})
	log.Println(err)
	// out
	// go test -run TestExponentialBackoff
	// 2026/02/09 13:45:50 running......
	// 2026/02/09 13:45:51 running......
	// 2026/02/09 13:45:53 running......
	// 2026/02/09 13:45:57 running......
	// 2026/02/09 13:46:05 running......
	// 2026/02/09 13:46:21 running......
	// 2026/02/09 13:46:21 error
}

func TestRandExponentialBackoff(t *testing.T) {
	// 重试6次，时间间隔为指数退幂最大休眠时间为16秒（1、2、4、8、16、16......）加上随机时间
	err := Retry(6, RandExponentialBackoff(16, time.Second, 100, 300, time.Millisecond), func() error {
		log.Println("running......")
		return errors.New("error")
	})
	log.Println(err)
	// out
	// go test -run TestRandExponentialBackoff
	// 2026/02/09 13:48:05 running......
	// 2026/02/09 13:48:06 running......
	// 2026/02/09 13:48:08 running......
	// 2026/02/09 13:48:12 running......
	// 2026/02/09 13:48:21 running......
	// 2026/02/09 13:48:37 running......
	// 2026/02/09 13:48:37 error
}
```