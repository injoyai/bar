package volume

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/injoyai/conv"
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
	if this == 0 {
		return 0, mapSizeUnit[0]
	}

	f := float64(this)
	i := 0

	for f >= 1024 && i < len(mapSizeUnit)-1 {
		f /= 1024
		i++
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
