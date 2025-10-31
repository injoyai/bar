package bar

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"strings"
)

type PlanOption func(p *Plan)

func NewPlan(op ...PlanOption) *Plan {
	p := &Plan{
		prefix:  "[",
		suffix:  "]",
		style:   "■",
		padding: " ",
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

func (p *Plan) String(current, total int64) string {
	if total <= 0 {
		return fmt.Sprintf("%s%s%s", p.prefix, strings.Repeat(" ", p.width), p.suffix)
	}

	//计算当前进度
	rate := float64(current) / float64(total)
	rate = conv.Range(rate, 0, 1)

	//进度条的数量
	count := int(float64(p.width) * rate)
	if count < 0 {
		count = 0
	}

	//处理样式
	styleRunes := []rune(p.style)
	lenStyle := len(styleRunes)
	if lenStyle == 0 {
		styleRunes = []rune{'■'}
		lenStyle = 1
	}

	var b strings.Builder
	for i := 0; i < count; i++ {
		b.WriteRune(styleRunes[i%lenStyle])
	}

	paddingRunes := []rune(p.padding)
	lenPadding := len(paddingRunes)
	if lenPadding == 0 {
		paddingRunes = []rune{' '}
		lenPadding = 1
	}
	// 补全剩余部分（未完成区域）
	for i := count; i < p.width; i++ {
		b.WriteRune(paddingRunes[i%lenPadding])
	}

	barStr := fmt.Sprintf("%s%s%s", p.prefix, b.String(), p.suffix)

	if p.color != nil {
		barStr = p.color.Sprint(barStr)
	}
	return barStr
}
