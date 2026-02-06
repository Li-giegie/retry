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
