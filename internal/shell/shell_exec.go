package shell

import (
	"bytes"
	"io"
	"os/exec"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	Bash       = &Shell{_bash{}}
	SH         = &Shell{_sh{}}
	CMD        = &Shell{_cmd{}}
	Default    *Shell
	defaultUse Use
)

func init() {
	if runtime.GOOS == "windows" {
		defaultUse = _cmd{}
	} else {
		defaultUse = _bash{}
	}
	Default = &Shell{defaultUse}
}

func Run(args ...string) error {
	return Default.Run(args...)
}

/*



 */

type _cmd struct{}

func (_cmd) Prefix() [2]string { return [2]string{"cmd", "/c"} }

func (_cmd) Decode(p []byte) ([]byte, error) { return GbkToUtf8(p) }

type _bash struct{}

func (_bash) Prefix() [2]string { return [2]string{"bash", "-c"} }

func (_bash) Decode(p []byte) ([]byte, error) { return p, nil }

type _sh struct{}

func (_sh) Prefix() [2]string { return [2]string{"sh", "-c"} }

func (_sh) Decode(p []byte) ([]byte, error) { return p, nil }

type Use interface {
	Prefix() [2]string
	Decode(p []byte) ([]byte, error)
}

type Shell struct {
	Use
}

func (this *Shell) Run(args ...string) error {
	pre := this.Prefix()
	list := append(pre[1:], args...)
	cmd := exec.Command(pre[0], list...)
	return cmd.Run()
}

type Result struct {
	buf    *bytes.Buffer
	str    *string
	decode func(p []byte) ([]byte, error)
}

func (this *Result) String() string {
	if this.str == nil {
		// 优先尝试解码（如 GBK -> UTF-8）；解码失败再回退到原始内容
		if this.decode != nil {
			if bs, err := this.decode(this.buf.Bytes()); err == nil {
				s := string(bs)
				this.str = &s
				return s
			}
		}
		s := this.buf.String()
		this.str = &s
	}
	return *this.str
}

func GbkToUtf8(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}
