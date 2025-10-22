package bar

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"time"
)

//type Coroutine interface {
//	Bar
//	Wait()
//	Go(f func())
//	GoRetry(f func() error, retry int)
//}

type Bar interface {
	fmt.Stringer
	io.Closer
	Add(n int64)              //添加数量
	Set(current int64)        //设置当前数量
	SetCurrent(current int64) //设置当前数量
	SetTotal(total int64)     //设置总数量

	/*
		SetFormat 设置样式,例:
		b.SetFormat(
			WithPlan(op...),
			WithRateSize(),
			WithSpeed(),
			WithRemain(),
		)
	*/
	SetFormat(format ...Format)
	SetPrefix(prefix string) //设置前缀
	SetSuffix(suffix string) //设置后缀
	SetWriter(w io.Writer)   //设置writer
	OnSet(f func())          //设置事件
	OnFinal(f func(b Bar))   //完成事件

	Last() int64                  //最后数量
	Current() int64               //当前数量
	Total() int64                 //总数量
	StartTime() time.Time         //开始时间
	LastTime() time.Time          //最后时间
	Flush() bool                  //刷入writer
	Done() <-chan struct{}        //完成
	Logf(format string, a ...any) //在bar上方输出日志
	Log(a ...any)                 //在bar上方输出日志

	Download(source, filename string, proxy ...string) (int64, error) //通过http下载
	Copy(w io.Writer, r io.Reader) (int64, error)                     //复制
}

type Option func(b Bar)
type PlanOption func(p Plan)
type Format func(b Bar) string

type Plan interface {
	SetPrefix(s string)         //设置前缀
	SetSuffix(s string)         //设置后缀
	SetStyle(s string)          //设置样式 ■ □ # >
	SetColor(c color.Attribute) //设置颜色
	SetWidth(w int)             //设置宽度
}
