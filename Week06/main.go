package main

import (
	. "Week06/rate_limiter"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func TestRing() {
	size := 4
	r := &Ring{
		Data:  make([]int, size),
		HeadP: 0,
	}

	// [0 1 2 3]
	// p = 0
	for i := 0; i < r.Size(); i++ {
		r.Access(i, i, func(v *int) { *v = i })
	}
	if s := r.Sum(); s != 6 {
		fmt.Printf("expected 6, nut got %d\n", s)
	}

	// [0 0 2 3]
	// p = 2
	r.Move(2)
	if s := r.Sum(); s != 5 {
		fmt.Printf("expected 5, nut got %d\n", s)
	}

	// [10 5 2 3]
	// p = 2
	r.Access(2, 2, func(v *int) { *v = 10 })
	r.Access(3, 3, func(v *int) { *v = 5 })
	if s := r.Sum(); s != 20 {
		fmt.Printf("expected 20, nut got %d\n", s)
	}

	// 不移动
	r.Move(0)
	if s := r.Sum(); s != 20 {
		fmt.Printf("expected 20, nut got %d\n", s)
	}

	// [0 0 0 0]
	// p = 2
	r.Move(r.Size())
	if s := r.Sum(); s != 0 {
		fmt.Printf("expected 0, nut got %d\n", s)
	}
}

func Init() {
	// 限速器每秒接受10次访问
	// 并发访问 12 次，失败两次
	var (
		s1        = NewSlidingWindowLimiter(10)
		errCount  int64
		wg1       sync.WaitGroup
		taskCount = 12
	)
	wg1.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			if err := s1.Allow(); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
			wg1.Done()
		}()
	}
	wg1.Wait()
	if errCount != 2 {
		fmt.Printf("expect 2, but got %d\n", errCount)
	}
}

func WithInterval() {
	// 限速器每秒可接受3个访问
	// 第一个 100ms，并发访问3次，都能成功访问
	// 过去 100ms 后，并发访问4次，失败一次
	var (
		s1        = NewSlidingWindowLimiter(3)
		errCount  int64
		wg1       sync.WaitGroup
		taskCount = 3
	)
	wg1.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			if err := s1.Allow(); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
			wg1.Done()
		}()
	}
	wg1.Wait()
	if errCount != 0 {
		fmt.Printf("errcount should be 0,but got %v\n", errCount)
	}
	time.Sleep(time.Millisecond * 100)
	taskCount += 1
	wg1.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			if err := s1.Allow(); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
			wg1.Done()
		}()
	}
	wg1.Wait()
	if errCount != 1 {
		fmt.Printf("errcount should be 1, but got %d\n", errCount)
	}
}

func LongInterval() {
	// 限速器每秒可访问 10 次
	// 测试梅 100ms 访问一次，时长 2s ，共 20 个请求，不应该报错
	var (
		s1 = NewSlidingWindowLimiter(10)
	)
	for i := 0; i < 20; i++ {
		if err := s1.Allow(); err != nil {
			fmt.Printf("unexpect err in loop %d\n", i)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
func main() {
	TestRing()
	Init()
	WithInterval()
	LongInterval()

}
