package bar

import (
	"fmt"
	"github.com/injoyai/bar/internal/volume"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"time"
)

// WithPlan 进度条,例 [>>>   ]
func WithPlan(op ...PlanOption) Format {
	p := &Plan{
		prefix: "[",
		suffix: "]",
		style:  "■",
		color:  nil,
		width:  50,
	}
	for _, v := range op {
		v(p)
	}
	return func(b *Bar) string {
		p.current = b.Current()
		p.total = b.Total()
		return p.String()
	}
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
		return time.Now().Sub(b.StartTime()).String()
	}
}

// WithUsedSecond 已经耗时,例 600s
func WithUsedSecond() Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%0.1fs", time.Now().Sub(b.StartTime()).Seconds())
	}
}

// WithRemain 预计剩余时间 例 1m18s
func WithRemain() Format {
	return func(b *Bar) string {
		rate := float64(b.Current()) / float64(b.Total())
		spend := time.Now().Sub(b.StartTime())
		remain := "0s"
		if rate > 0 {
			sub := time.Duration(float64(spend)/rate - float64(spend))
			remain = (sub - sub%time.Second).String()
		}
		return remain
	}
}

// WithCurrentSize 大小,例 58B,需传指针,不然不会变
func WithCurrentSize(size *int64) Format {
	return func(b *Bar) string {
		return volume.SizeString(*size)
	}
}

// WithCurrentRateSizeUnit 大小,例 58B/100B,需传指针,不然不会变
func WithCurrentRateSizeUnit(size, total *int64) Format {
	return func(b *Bar) string {
		return fmt.Sprintf("%s/%s", volume.SizeString(*size), volume.SizeString(*total))
	}
}
