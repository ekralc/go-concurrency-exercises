//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	ch := make(chan bool)
	go func() {
		process()
		ch <- true
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ch:
			return true
		case <-ticker.C:
			if timeUsed := atomic.AddInt64(&u.TimeUsed, 1); timeUsed >= 10 {
				return false
			} else {
				fmt.Printf("User %v has accumulated %v seconds\n", u.ID, timeUsed)
			}
		}
	}
}

func main() {
	RunMockServer()
}
