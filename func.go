package bar

import (
	"github.com/injoyai/bar/internal/m3u8"
	"github.com/injoyai/bar/internal/util"
	"github.com/injoyai/base/chans"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

func Copy(w io.Writer, r io.Reader, total int64) (int64, error) {
	return New(WithTotal(total)).Copy(w, r)
}

func Download(url, filename string, proxy ...string) (int64, error) {
	return New().Download(url, filename, proxy...)
}

func DownloadHLS(source, dir string, op ...HLSOption) error {

	cfg := &DownloadHLSConfig{
		Proxy:       "",
		Coroutine:   10,
		ShowDetails: false,
	}

	for _, v := range op {
		v(cfg)
	}

	os.MkdirAll(dir, os.ModePerm)

	ls, err := m3u8.Decode(source)
	if err != nil {
		return err
	}

	current := int64(0)
	total := int64(0)
	index := int64(0)
	b := New(
		WithTotal(int64(len(ls))),
		WithFormat(
			WithPlan(),
			WithRateSize(),
			WithCurrentRateSizeUnit(&current, &total),
			WithRemain()),
	)

	h := util.NewClient().SetTimeout(0)
	if err := h.SetProxy(cfg.Proxy); err != nil {
		return err
	}

	f := func(u string, n int64, log bool) {
		atomic.AddInt64(&index, 1)
		atomic.AddInt64(&current, n)
		atomic.StoreInt64(&total, (current/index)*int64(len(ls)))
		if log {
			b.Log(u)
		}
		b.Add(1)
		b.Flush()
	}

	wg := chans.NewWaitLimit(cfg.Coroutine)
	for _, u := range ls {
		wg.Add()
		go func(u string) {
			defer wg.Done()

			_u, err := url.Parse(u)
			if err != nil {
				b.Log("[错误]", err)
				return
			}

			filename := filepath.Join(dir, filepath.Base(_u.Path))
			if !strings.HasSuffix(filename, ".ts") {
				filename += ".ts"
			}

			stat, exist, err := Stat(filename)
			if err != nil {
				b.Log("[错误]", err)
				b.Flush()
				return
			} else if exist {
				f(u, stat.Size(), false)
				return
			}

			n, err := h.GetToFile(u, filename)
			if err != nil {
				b.Log("[错误]", err)
				return
			}

			f(u, n, cfg.ShowDetails)
		}(u)
	}

	wg.Wait()

	return nil
}

type DownloadHLSConfig struct {
	Proxy       string
	Coroutine   int
	ShowDetails bool
}

type HLSOption func(c *DownloadHLSConfig)

func WithHLSProxy(proxy string) HLSOption {
	return func(c *DownloadHLSConfig) {
		c.Proxy = proxy
	}
}
func WithHLSCoroutine(coroutine int) HLSOption {
	return func(c *DownloadHLSConfig) {
		c.Coroutine = coroutine
	}
}
func WithHLSShowDetails(b ...bool) HLSOption {
	return func(c *DownloadHLSConfig) {
		c.ShowDetails = len(b) == 0 || b[0]
	}
}

// Stat 获取文件信息
func Stat(filename string) (os.FileInfo, bool, error) {
	stat, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, false, err
	} else if err != nil && os.IsNotExist(err) {
		return nil, false, nil
	}
	return stat, true, nil
}
