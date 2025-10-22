package bar

import "github.com/injoyai/base/chans"

func NewCoroutine(total, limit int, op ...Option) *Coroutine {
	b := New(op...)
	b.SetTotal(int64(total))
	b.Flush()
	return &Coroutine{
		Bar: b,
		wg:  chans.NewWaitLimit(limit),
	}
}

type Coroutine struct {
	Bar
	wg chans.WaitLimit
}

func (this *Coroutine) Wait() {
	this.wg.Wait()
}

func (this *Coroutine) Go(f func()) {
	this.GoRetry(func() error {
		f()
		return nil
	}, 1)
}

func (this *Coroutine) GoRetry(f func() error, retry int) {
	if f == nil {
		return
	}
	this.wg.Add()
	go func() {
		defer func() {
			this.Bar.Add(1)
			this.Bar.Flush()
			this.wg.Done()
		}()
		for i := 0; i < retry; i++ {
			if err := f(); err == nil {
				return
			}
		}
	}()
}
