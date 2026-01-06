package bar

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/injoyai/conv"
)

const (
	DefaultStyle   = "#"
	DefaultPadding = " "
)

type PlanOption func(p *Plan)

func NewPlan(op ...PlanOption) *Plan {
	p := &Plan{
		prefix:  "[",
		suffix:  "]",
		style:   DefaultStyle,
		padding: DefaultPadding,
		color:   nil,
		width:   50,
	}
	for _, o := range op {
		o(p)
	}
	return p
}

type Plan struct {
	prefix  string       //前缀 例 [
	suffix  string       //后缀 例 ]
	style   string       //进度条风格 例 ■
	padding string       //填充 例 .
	color   *color.Color //整体颜色
	width   int          //宽度
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

func (this *Plan) SetPadding(padding string) {
	this.padding = padding
}

func (this *Plan) SetWidth(width int) {
	this.width = width
}

func (this *Plan) SetColor(a color.Attribute) {
	this.color = color.New(a)
}

func (this *Plan) String(rate float64) string {

	//归一化
	rate = conv.Range(rate, 0, 1)

	//进度条的数量
	count := int(float64(this.width) * rate)
	if count < 0 {
		count = 0
	}

	//处理样式
	styleRunes := []rune(this.style)
	lenStyle := len(styleRunes)
	if lenStyle == 0 {
		styleRunes = []rune(DefaultStyle)
		lenStyle = 1
	}

	var b strings.Builder
	for i := 0; i < count; i++ {
		b.WriteRune(styleRunes[i%lenStyle])
	}

	paddingRunes := []rune(this.padding)
	lenPadding := len(paddingRunes)
	if lenPadding == 0 {
		paddingRunes = []rune(DefaultPadding)
		lenPadding = 1
	}
	// 补全剩余部分（未完成区域）
	for i := count; i < this.width; i++ {
		b.WriteRune(paddingRunes[i%lenPadding])
	}

	barStr := fmt.Sprintf("%s%s%s", this.prefix, b.String(), this.suffix)

	if this.color != nil {
		barStr = this.color.Sprint(barStr)
	}
	return barStr
}
