package bar

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/injoyai/bar/internal/util"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
)

type Option func(b *Bar)
type Format func(b *Bar) string

// WithCurrent 设置当前数量
func WithCurrent(current int64) Option {
	return func(b *Bar) {
		b.SetCurrent(current)
	}
}

// WithTotal 设置总数量
func WithTotal(total int64) Option {
	return func(b *Bar) {
		b.SetTotal(total)
	}
}

// WithOption 批量设置
func WithOption(op ...Option) Option {
	return func(b *Bar) {
		for _, v := range op {
			v(b)
		}
	}
}

// WithFormat 设置样式
func WithFormat(fs ...Format) Option {
	return func(b *Bar) {
		b.SetFormat(fs...)
	}
}

// WithFormatDefault 设置默认样式,不带单位
func WithFormatDefault(op ...PlanOption) Option {
	return func(b *Bar) {
		b.SetFormat(
			WithPlan(op...),
			WithRateSize(),
			WithSpeed(),
			WithRemain2(),
		)
	}
}

// WithFormatDefaultUnit 设置默认样式,带单位
func WithFormatDefaultUnit(op ...PlanOption) Option {
	return func(b *Bar) {
		b.SetFormat(
			WithPlan(op...),
			WithRateSizeUnit(),
			WithSpeedUnit(),
			WithRemain2(),
		)
	}
}

// WithFormatSplit 设置分隔符
func WithFormatSplit(split string) Option {
	return func(b *Bar) {
		b.SetFormatSplit(split)
	}
}

// WithPrefix 设置前缀
func WithPrefix(prefix string) Option {
	return func(b *Bar) {
		b.SetPrefix(prefix)
	}
}

// WithSuffix 设置后缀
func WithSuffix(suffix string) Option {
	return func(b *Bar) {
		b.SetSuffix(suffix)
	}
}

// WithWriter 设置writer
func WithWriter(writer io.Writer) Option {
	return func(b *Bar) {
		b.SetWriter(writer)
	}
}

func WithFinal(f Option) Option {
	return func(b *Bar) {
		b.OnFinal(f)
	}
}

// WithAutoFlush 设置后自动刷新
func WithAutoFlush() Option {
	return func(b *Bar) {
		b.OnSet(func(b *Bar) {
			b.Flush()
		})
	}
}

// WithIntervalFlush 设置定时刷新
func WithIntervalFlush(interval time.Duration) Option {
	return func(b *Bar) {
		go func() {
			t := time.NewTimer(interval)
			defer t.Stop()
			for {
				select {
				case <-b.Done():
					return
				case <-t.C:
					b.Flush()
				}
			}
		}()
	}
}

// WithFlush 刷入writer
func WithFlush() Option {
	return func(b *Bar) {
		b.Flush()
	}
}

func New(op ...Option) *Bar {
	b := &Bar{
		current:     0,
		total:       0,
		formatSplit: "  ",
		writer:      os.Stdout,
		startTime:   time.Now(),
		Closer:      safe.NewCloser(),
	}
	b.SetCloseFunc(func(err error) error {
		if b.writer != nil {
			b.writer.Write([]byte("\n"))
		}
		if b.onFinal != nil {
			b.onFinal(b)
		}
		return nil
	})
	WithFormatDefault()(b)
	for _, v := range op {
		v(b)
	}
	return b
}

type Bar struct {
	current     int64          //当前数量
	total       int64          //总数量
	prefix      string         //前缀
	suffix      string         //后缀
	format      Format         //格式化
	formatSplit string         //分隔符
	writer      io.Writer      //输出
	onSet       []func(b *Bar) //设置事件
	onFinal     func(b *Bar)   //完成事件

	startTime time.Time  //开始时间
	last      int64      //最后一次增加的值
	lastTime  time.Time  //最后一次时间
	mu        sync.Mutex //并发锁

	*safe.Closer //closer
}

func (this *Bar) Add(n int64) {
	this.mu.Lock()
	this.current = this.current + n
	if this.current > this.total {
		this.current = this.total
	}
	this.last = n
	this.lastTime = time.Now()
	this.mu.Unlock()

	for _, f := range this.onSet {
		if f != nil {
			f(this)
		}
	}
}

func (this *Bar) Set(current int64) {
	this.SetCurrent(current)
}

func (this *Bar) SetCurrent(current int64) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if current > this.total {
		current = this.total
	}
	this.last = current - this.current
	this.lastTime = time.Now()
	this.current = current

	for _, f := range this.onSet {
		if f != nil {
			f(this)
		}
	}
}

func (this *Bar) SetTotal(total int64) {
	this.total = total
}

func (this *Bar) SetFormat(fs ...Format) {
	switch len(fs) {
	case 0:
		this.format = func(b *Bar) string { return "" }

	case 1:
		this.format = fs[0]

	default:
		ls := make([]string, len(fs))
		this.format = func(b *Bar) string {
			for i, v := range fs {
				ls[i] = v(b)
			}
			return strings.Join(ls, this.formatSplit)
		}

	}
}

func (this *Bar) SetFormatSplit(split string) {
	this.formatSplit = split
}

func (this *Bar) SetPrefix(prefix string) {
	this.prefix = prefix
}

func (this *Bar) SetSuffix(suffix string) {
	this.suffix = suffix
}

func (this *Bar) SetWriter(w io.Writer) {
	this.writer = w
}

func (this *Bar) OnSet(f func(b *Bar)) {
	this.onSet = append(this.onSet, f)
}

func (this *Bar) OnFinal(f Option) {
	this.onFinal = f
}

/*



 */

func (this *Bar) Last() int64 {
	return this.last
}

func (this *Bar) Current() int64 {
	return this.current
}

func (this *Bar) Total() int64 {
	return this.total
}

func (this *Bar) Rate() float64 {
	if this.total == 0 {
		return 0
	}
	return float64(this.current) / float64(this.total)
}

func (this *Bar) StartTime() time.Time {
	return this.startTime
}

func (this *Bar) LastTime() time.Time {
	return this.lastTime
}

func (this *Bar) Logf(format string, a ...any) {
	s := "\r\033[K" + fmt.Sprintf(format, a...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		s += "\n"
	}
	this.writer.Write([]byte(s))
}

func (this *Bar) Log(a ...any) {
	s := "\r\033[K" + fmt.Sprintln(a...)
	this.writer.Write([]byte(s))
}

func (this *Bar) Flush() (closed bool) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.Closed() {
		return true
	}
	if this.writer == nil {
		return false
	}
	s := this.String()
	if s == "" || s[0] != '\r' {
		s = "\r\033[K" + s
	}
	this.writer.Write([]byte(s))
	if this.current >= this.total {
		this.Close()
	}
	return this.Closed()
}

func (this *Bar) String() string {
	return this.prefix + this.format(this) + this.suffix
}

/*



 */

/*



 */

func (this *Bar) Download(source, filename string, proxy ...string) (int64, error) {
	defer this.Close()

	//下载大文件的时候需要设置长的超时时间
	h := util.NewClient().SetTimeout(0)
	if err := h.SetProxy(conv.Default("", proxy...)); err != nil {
		return 0, err
	}

	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return 0, err
	}

	return h.GetToFile(source, filename)
}

func (this *Bar) Copy(w io.Writer, r io.Reader) (int64, error) {
	return this.CopyN(w, r, 4<<10)
}

func (this *Bar) CopyN(w io.Writer, r io.Reader, bufSize int64) (int64, error) {
	buff := bufio.NewReader(r)
	reader := &Reader{
		Reader: buff,
		Bar:    this,
	}
	return io.CopyN(w, reader, bufSize)
}

func NewReader(r io.Reader, b *Bar) *Reader {
	return &Reader{
		Reader: r,
		Bar:    b,
	}
}

type Reader struct {
	io.Reader
	*Bar
}

func (this *Reader) Read(p []byte) (n int, err error) {
	n, err = this.Reader.Read(p)
	this.Bar.Add(int64(n))
	this.Bar.Flush()
	return
}
