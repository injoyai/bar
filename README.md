### 实现控制台上的进度条

### 如何使用
* 正常使用
```go
import (
    "github.com/injoyai/bar"
    "time"
)

func main() {
    x := bar.New(
        bar.WithTotal(60),  
        bar.WithFormatDefault(func(p *bar.Plan) {
            p.SetStyle("■")
            p.SetPadding(".")
        }),
    )
    for {
        time.Sleep(time.Millisecond * 100)
        x.Add(1)
        if x.Flush() {
            break
        }
    }
}
```
```shell
[■■■■■■■■■■■■■■■■■.................................]  21/60  9.1/s  4s
```
* 协程模式
```go
import (
    "github.com/injoyai/bar"
	"time"
)

func main() {
    // 协程数量
    limit:=2
    // 能自动增长,和并发数量控制
    b := bar.NewCoroutine(100,limit)
    for i := 0; i < 100; i++ {
        b.Go(func () {
            <-time.After(time.Second * 1)
        })
    }
    b.Wait()
}
```
```shell
[■■■■■■■■■■■■■■■■■                                 ]  21/100  9.1/s  4s
```
* 动画效果
```go
import (
	"github.com/injoyai/bar"
	"time"
)

func main() {
	b := bar.New(
		bar.WithTotal(100),
		bar.WithFormat(
			bar.WithAnimationMoon(),
			bar.WithText("|"),
			bar.WithAnimation(bar.Animations[69]),
			bar.WithText("|"),
			bar.WithAnimationSnake(),
			bar.WithRate(),
		),
	)
	for i := 0; i < 100; i++ {
		b.Add(1)
		b.Flush()
		time.Sleep(time.Millisecond * 500)
	}
}

```
```shell
🌓  |  ∙∙∙  |  ⣻  10.0%
```