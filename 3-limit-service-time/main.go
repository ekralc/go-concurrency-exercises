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
	ch := make(chan bool)
	go func() {
		process()
		ch <- true
	}()

	ticker := time.NewTicker(time.Second)
	go func() {
		if u.IsPremium {
			ch <- true
			return
		}
		for {
			<-ticker.C

			atomic.AddInt64(&u.TimeUsed, 1)
			fmt.Printf("User %v has accumulated %v seconds\n", u.ID, u.TimeUsed)

			if atomic.LoadInt64(&u.TimeUsed) >= 10 {
				ch <- false
				break
			}
		}
	}()

	res := <-ch

	return res
}

func main() {
	RunMockServer()
}
