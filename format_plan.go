package bar

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"strings"
)

type PlanOption func(p *Plan)

type Plan struct {
	prefix string       //前缀 例 [
	suffix string       //后缀 例 ]
	style  string       //进度条风格 例 >
	color  *color.Color //整体颜色
	width  int          //宽度

	current int64 //当前
	total   int64 //总数
}

func (this *Plan) SetPrefix(prefix string) {
	this.prefix = prefix
}

func (this *Plan) SetSuffix(suffix string) {
	this.suffix = suffix
}

func (this *Plan) SetStyle(style string) {
	this.style = style
}

func (this *Plan) SetWidth(width int) {
	this.width = width
}

func (this *Plan) SetColor(a color.Attribute) {
	this.color = color.New(a)
}

func (this *Plan) String() string {
	lenStyle := len([]rune(this.style))
	rate := float64(this.current) / float64(this.total)
	count := int(float64(this.width) * rate / float64(lenStyle))
	count = conv.Select(count < 0, 0, count)
	nowWidth := strings.Repeat(this.style, count)
	if rate == 1 {
		for i := 0; len([]rune(nowWidth)) < this.width; i++ {
			nowWidth += string([]rune(this.style)[i%lenStyle])
		}
	}
	barStr := fmt.Sprintf(fmt.Sprintf("%s%%-%ds%s", this.prefix, this.width, this.suffix), nowWidth)
	if this.color != nil {
		barStr = this.color.Sprint(barStr)
	}
	return barStr
}
