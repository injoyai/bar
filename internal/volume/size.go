package volume

import (
	"fmt"
	"github.com/injoyai/conv"
	"math"
	"strings"
	"time"
)

const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
	EB = PB * 1024
)

// 64位最大只能到15.999EB
var mapSizeUnit = []string{
	"B",
	"KB",
	"MB",
	"GB",
	"TB",
	"PB",
	"EB", //64位最大单位
	"ZB",
	"YB",
	"BB",
	"NB",
	"DB",
	"CB",
	"XB",
}

type Volume uint64

func (this Volume) Uint64() uint64 {
	return uint64(this)
}

func (this Volume) String() string {
	size, unit := this.SizeUnit()
	return fmt.Sprintf("%v%s", size, unit)
}

func (this Volume) SizeUnit() (float64, string) {
	i := 0
	//先用uint64进行循环除以1024,float64的指数比uint64小
	//当值比float64的最大值小的时候,能转成float64时,使用float64进行除以1024
	for ; float64(this) < 0; i++ {
		this = this / 1024
	}
	f := float64(this)
	for ; f >= 1024*1024; i++ {
		f = f / 1024
	}
	if f >= 1024 {
		return f / 1024, mapSizeUnit[i+1]
	}
	return f, mapSizeUnit[i]
}

// ParseVolume 解析体积,
func ParseVolume(s string) Volume {
	total := Volume(0)
	size := ""
	unit := ""
	hasUnit := false

	add := func() {
		if hasUnit {
			for i, u := range mapSizeUnit {
				if strings.ToUpper(unit) == u {
					total += Volume(conv.Float64(size) * math.Pow(1024, float64(i)))
				}
			}
			hasUnit = false
			size = ""
			unit = ""
		}
	}

	for _, v := range s {
		if (v >= '0' && v <= '9') || v == '.' {
			add()
			size += string(v)
		} else {
			unit += string(v)
			hasUnit = true
		}
	}
	add()

	return total
}

// SizeUnit 字节数量和单位 例 15.8,"MB"
// 64位最大值是 18446744073709551616 = 15.999EB
func SizeUnit(b int64) (float64, string) {
	return Volume(b).SizeUnit()
}

// SizeString 字节数量字符表现方式,例 15.8MB, 会四舍五入
func SizeString(b int64, decimal ...int) string {
	size, unit := SizeUnit(b)
	d := conv.Default(1, decimal...)
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%s", d), size, unit)
}

// SizeSpeed 每秒速度 例15.8MB/s
func SizeSpeed(b int64, sub time.Duration, decimal ...int) string {
	size, unit := SizeUnit(b)
	spend := size / sub.Seconds()
	d := conv.Default(1, decimal...)
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%s/s", d), spend, unit)
}
