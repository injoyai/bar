package bar

import (
	"fmt"
	"sync"
	"time"

	"github.com/injoyai/bar/internal/volume"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
)

// WithPlan 进度条,例 [>>>   ]
func WithPlan(op ...PlanOption) Format {
	p := NewPlan(op...)
	return func(b *Bar) string {
		return p.String(b.Rate())
	}
}

// WithText 文本
func WithText(a ...any) Format {
	return func(b *Bar) string {
		return fmt.Sprint(a...)
	}
}

// WithTime 时间
func WithTime() Format {
	return func(b *Bar) string {
		return time.Now().Format(time.TimeOnly)
	}
}

// WithDate 日期
func WithDate() Format {
	return func(b *Bar) string {
		return time.Now().Format(time.DateOnly)
	}
}

// WithDateTime 日期时间
func WithDateTime() Format {
	return func(b *Bar) string {
		return time.Now().Format(time.DateTime)
	}
}

// WithAnimation 进度动画
func WithAnimation(ls []string) Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%s", ls[int(b.Current())%len(ls)])
	}
}

// WithAnimationSnake 进度动画: 贪吃蛇
func WithAnimationSnake() Format {
	return WithAnimation(Animations[11])
}

// WithAnimationMoon 进度动画: 月亮
func WithAnimationMoon() Format {
	return WithAnimation(Animations[70])
}

// WithRate 进度百分比,例 58%
func WithRate() Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%0.1f%%", float64(b.Current())*100/float64(b.Total()))
	}
}

// WithRateSize //进度数量,例 58/100
func WithRateSize() Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), b.Total())
	}
}

// WithRateSizeUnit 进度数量带单位,例 58B/100B
func WithRateSizeUnit() Format {
	return func(b *Bar) string {
		currentNum, currentUnit := volume.SizeUnit(b.Current())
		totalNum, totalUnit := volume.SizeUnit(b.Total())
		return fmt.Sprintf("%0.1f%s/%0.1f%s", currentNum, currentUnit, totalNum, totalUnit)
	}
}

func speed(cache *maps.Safe, key string, size int64, expiration time.Duration, f func(float64) string) string {

	timeKey := "time_" + key
	cacheKey := "speed_" + key
	//最后的数据时间
	lastTime, _ := cache.GetOrSetByHandler(timeKey, func() (any, error) {
		return time.Time{}, nil
	})

	//记录这次时间,用于下次计算时间差
	now := time.Now()
	cache.Set(timeKey, now)

	//尝试从缓存获取速度,存在则直接返回,由expiration控制
	if val, ok := cache.Get(cacheKey); ok {
		return val.(string)
	}

	//计算速度
	size = conv.Select(size >= 0, size, 0)
	spendSize := float64(size) / now.Sub(lastTime.(time.Time)).Seconds()
	s := f(spendSize)
	cache.Set(cacheKey, s, expiration)
	return s
}

// WithSpeed //进度速度,例 13/s
func WithSpeed(expiration ...time.Duration) Format {
	cache := maps.NewSafe()
	return func(b *Bar) string {
		return speed(cache, "Speed", b.Last(), conv.Default(time.Millisecond*500, expiration...), func(size float64) string {
			return fmt.Sprintf("%0.1f/s", size)
		})
	}
}

// WithSpeedUnit //进度速度带单位,例 13MB/s
func WithSpeedUnit(expiration ...time.Duration) Format {
	cache := maps.NewSafe()
	return func(b *Bar) string {
		return speed(cache, "SpeedUnit", b.Last(), conv.Default(time.Millisecond*500, expiration...), func(size float64) string {
			f, unit := volume.SizeUnit(int64(size))
			return fmt.Sprintf("%0.1f%s/s", f, unit)
		})
	}
}

// WithSpeedAvg //进度平均速度,例 13/s
func WithSpeedAvg() Format {
	return func(b *Bar) string {
		speedSize := float64(b.Current()) / time.Since(b.StartTime()).Seconds()
		return fmt.Sprintf("%0.1f/s", speedSize)
	}
}

// WithSpeedUnitAvg //进度平均速度带单位,例 13MB/s
func WithSpeedUnitAvg() Format {
	return func(b *Bar) string {
		speedSize := float64(b.Current()) / time.Since(b.StartTime()).Seconds()
		f, unit := volume.SizeUnit(int64(speedSize))
		return fmt.Sprintf("%0.1f%s/s", f, unit)
	}
}

// WithUsed 已经耗时,例 2m20s
func WithUsed() Format {
	return func(b *Bar) string {
		return time.Since(b.StartTime()).String()
	}
}

// WithUsedSecond 已经耗时,例 600s
func WithUsedSecond() Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%0.1fs", time.Since(b.StartTime()).Seconds())
	}
}

// WithRemain 预计剩余时间(根据所有数据来计算) 例 1m18s
func WithRemain() Format {
	return func(b *Bar) string {
		rate := float64(b.Current()) / float64(b.Total())
		spend := time.Since(b.StartTime())
		remain := "-"
		if rate > 0 {
			sub := time.Duration(float64(spend)/rate - float64(spend))
			remain = (sub - sub%time.Second).String()
		}
		return remain
	}
}

// WithRemain2 预计剩余时间(根据最近的几个数据来计算) 例 1m18s
// n: 滑动窗口大小（即"最近的几次写入"），不传则在首次采样时根据 Total 自适应到 [500,1000]
func WithRemain2(n ...int) Format {
	type node struct {
		current int64
		time    time.Time
	}

	var (
		mu       sync.Mutex
		buf      []node // 环形缓冲区
		head     int    // 队首下标（最旧）
		size     int    // 当前元素数
		capacity int    // 窗口容量，0 表示尚未初始化
		hooked   bool   // 是否已注册 OnSet 回调，避免重复注册
	)

	// 初始化容量（懒加载，依赖首次调用时 Bar 的 Total）
	initCap := func(b *Bar) {
		if capacity > 0 {
			return
		}
		_n := conv.Range(int(b.Total()/10), 500, 1000)
		_n = conv.Default(_n, n...)
		if _n <= 0 {
			_n = 500
		}
		capacity = _n
		buf = make([]node, capacity)
	}

	push := func(cur int64, t time.Time) {
		// 写入环形缓冲区尾部
		tail := (head + size) % capacity
		buf[tail] = node{current: cur, time: t}
		if size < capacity {
			size++
		} else {
			// 已满，覆盖最旧元素，head 后移
			head = (head + 1) % capacity
		}
	}

	return func(b *Bar) string {
		mu.Lock()
		if !hooked {
			initCap(b)
			b.OnSet(func(b *Bar) {
				mu.Lock()
				defer mu.Unlock()
				push(b.Current(), b.LastTime())
			})
			hooked = true
		}

		if size == 0 {
			mu.Unlock()
			return "-"
		}
		start := buf[head]
		end := buf[(head+size-1)%capacity]
		mu.Unlock()

		subCurrent := end.current - start.current
		subTime := end.time.Sub(start.time)

		// 需要同时满足：有数据进展 + 时间有流逝，避免除零或负值
		if subCurrent <= 0 || subTime <= 0 {
			return "-"
		}

		avgSpeed := float64(subTime) / float64(subCurrent) //平均耗时(每单位)
		remainNumber := b.Total() - end.current            //剩余数量
		if remainNumber <= 0 {
			return "0s"
		}
		remainTime := time.Duration(float64(remainNumber) * avgSpeed) //剩余时间
		return (remainTime - remainTime%time.Second).String()
	}
}

// WithCustomSize 大小,例 58B,需传指针,不然不会变
func WithCustomSize(size *int64) Format {
	return func(b *Bar) string {
		return volume.SizeString(*size)
	}
}

// WithCustomRateSizeUnit 大小,例 58B/100B,需传指针,不然不会变
func WithCustomRateSizeUnit(size, total *int64) Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%s/%s", volume.SizeString(*size), volume.SizeString(*total))
	}
}
