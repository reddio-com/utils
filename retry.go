package utils

import (
	"fmt"
	"time"
)

func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			fmt.Println("retrying", i, "attempts after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %v", attempts, err)
}

func RetryOnIgnore(attempts int, sleep time.Duration, f func() (bool, error)) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			fmt.Println("retrying", i, "attempts after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		ignore, err1 := f()
		if err1 == nil {
			return nil
		}
		if ignore {
			return fmt.Errorf("not retrying error: %v", err1)
		}
		err = err1
	}
	return fmt.Errorf("after %d attempts, last error: %v", attempts, err)
}
