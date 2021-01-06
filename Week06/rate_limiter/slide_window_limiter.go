package rate_limiter

import (
	"errors"
	"sync"
	"time"
)

type RateLimiter interface {
	// Allow 询问 limiter 是否可以访问，失败返回 error
	Allow() error
}

var ErrExceededLimit = errors.New("限流")

var _ RateLimiter = &SlidingWindowLimiter{}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	startTime   time.Time // win 开始的时间
	preMoveTime time.Time // 上一次 win 移动的时间
	interval    int64     // 移动间隔，按毫秒计
	buckets     Ring      // 每个间隔的请求数
	max         int       // 每秒最大请求数量
	lock        sync.Mutex
}

// NewSlidingWindowLimiter 为了简单起见，大部分配置写死
// 总共 10 个 bucket
// 每个 bucket 存储 100ms 中请求的数量 ，也就是说 interval 为 100ms
func NewSlidingWindowLimiter(maxPerSecond int) *SlidingWindowLimiter {
	now := time.Now()
	return &SlidingWindowLimiter{
		startTime:   now,
		preMoveTime: now,
		interval:    100,
		buckets:     Ring{Data: make([]int, 10)},
		max:         maxPerSecond,
	}
}

func (s *SlidingWindowLimiter) Allow() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	currTime := time.Now()

	// 在调用接口时更新滑动窗口的状态
	s.update(currTime)

	// 判断当前是否可以执行任务
	if s.buckets.Sum() >= s.max {
		return ErrExceededLimit
	}

	// 将相应的 Bucket 的值加 1
	s.addBucket(currTime, 1)
	return nil
}

func (s *SlidingWindowLimiter) update(currTime time.Time) {
	// 计算出距离上一次更新，滑动窗口应该移动的步数
	steps := int(currTime.Sub(s.preMoveTime).Milliseconds() / s.interval)
	s.buckets.Move(steps)
	s.startTime = s.startTime.Add(time.Duration(int64(steps)*s.interval) * time.Millisecond)
	s.preMoveTime = currTime
}

func (s *SlidingWindowLimiter) addBucket(currTime time.Time, inc int) {
	steps := int(currTime.Sub(s.startTime).Milliseconds() / s.interval)
	s.buckets.Access(steps, steps, func(v *int) { *v += inc })
}

// ring 特制的循环数组
type Ring struct {
	Data  []int
	HeadP int
}

func (r *Ring) Sum() (sum int) {
	for i := range r.Data {
		sum += r.Data[i]
	}
	return
}

func (r *Ring) Size() int {
	return len(r.Data)
}

// move 将 headP 循环移动 steps 长度
// 同时，将新旧 headP 之间的元素置 0
// 例如:
// 对于 [1 2 3] headP = 0
// 如果调用 move(2) ，则：[0 0 3] headP = 2
func (r *Ring) Move(steps int) {
	end := steps - 1
	if steps > r.Size() {
		end = r.Size() - 1
	}
	r.Access(0, end, func(v *int) { *v = 0 })
	r.HeadP = (r.HeadP + steps) % r.Size()
}

// 访问数组 [start end] 中的元素
// 对于调用者来说，index 是从 0 开始的，调用者不用关系它实际在数组中的位置
// 例如：
// 对于 [1 2 3] headP = 1
// access(1, 2, fn(...)) 其实访问的是 [3 1]
func (r *Ring) Access(start int, end int, fn func(v *int)) {
	for i := start; i <= end; i++ {
		fn(&r.Data[(r.HeadP+i)%r.Size()])
	}
}
